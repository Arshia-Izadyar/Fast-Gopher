package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/common"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/constants"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/cache"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

type KeyService struct {
	db               *sql.DB
	cfg              *config.Config
	whiteListService *WhiteListService
}

func NewKeyService(cfg *config.Config) *KeyService {
	db := postgres.GetDB()
	wl := NewWhiteListService(cfg)
	return &KeyService{
		db:               db,
		cfg:              cfg,
		whiteListService: wl,
	}
}

// service for /key
func (u *KeyService) GenerateKey(req *dto.GenerateKeyDTO) (*dto.KeyAcDTO, *service_errors.ServiceErrors) {
	tx, err := u.db.Begin()
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "database transaction start error", Status: fiber.StatusInternalServerError, Err: err}
	}
	defer tx.Rollback()

	key := common.GenerateSecureToken(u.cfg.Key.Len)
	if key == "" {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "internal error happened trying generate token please try again", Status: fiber.StatusInternalServerError}
	}
	saveKey := `
		INSERT INTO ac_keys(id) VALUES($1);
	`
	_, err = tx.Exec(saveKey, key)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, &service_errors.ServiceErrors{EndUserMessage: "the key already exists try getting another key", Status: fiber.StatusInternalServerError}
		}
		return nil, &service_errors.ServiceErrors{EndUserMessage: "insert query failed", Status: fiber.StatusInternalServerError}
	}

	saveSession := `
		INSERT INTO active_devices(session_id, ac_keys_id) VALUES($1, $2);
	`
	_, err = tx.Exec(saveSession, req.SessionId, key)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, &service_errors.ServiceErrors{EndUserMessage: fmt.Sprintf("session id (%s) already exists", req.SessionId), Status: fiber.StatusBadRequest}
		}
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to commit transaction", Status: fiber.StatusInternalServerError}
	}

	result, err := common.GenerateJwt(&dto.KeyDTO{Key: key, Premium: false, SessionId: req.SessionId}, u.cfg)
	// TODO: err handling
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to Generate jwt", Status: fiber.StatusInternalServerError}
	}
	return result, nil
}

func (u *KeyService) GenerateTokenFromKey(req *dto.KeyDTO) (*dto.KeyAcDTO, *service_errors.ServiceErrors) {

	tx, err := u.db.Begin()
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "database transaction start error", Status: fiber.StatusInternalServerError}
	}
	defer tx.Rollback()

	// Select premium
	q := `SELECT premium FROM ac_keys WHERE id = $1`
	err = tx.QueryRow(q, req.Key).Scan(&req.Premium)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &service_errors.ServiceErrors{EndUserMessage: "no active key found", Status: fiber.StatusNotFound}
		}
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to fetch key", Status: fiber.StatusInternalServerError}
	}

	// Save session
	saveSession := `INSERT INTO active_devices(session_id, ac_keys_id) VALUES($1, $2) ON CONFLICT (session_id, ac_keys_id) DO NOTHING;`
	_, err = tx.Exec(saveSession, req.SessionId, req.Key)
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to save session", Status: fiber.StatusInternalServerError}
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to commit transaction", Status: fiber.StatusInternalServerError}
	}

	// Generate JWT
	tk, err := common.GenerateJwt(req, u.cfg)
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to generate JWT", Status: fiber.StatusInternalServerError}
	}
	return tk, nil
}

// refresh
func (u *KeyService) Refresh(req *dto.RefreshTokenDTO) (*dto.KeyAcDTO, *service_errors.ServiceErrors) {

	// 	// 1. check if refresh is used
	// 	// 2. check if its is a refresh
	// 	// 3. blacklist refresh
	// 	// 4. issue new jwt

	claims, err := common.ValidateToken(req.RefreshToken, u.cfg)
	if err != nil {
		return nil, err
	}
	if claims[constants.AccessType] == true {
		return nil, &service_errors.ServiceErrors{EndUserMessage: service_errors.NotRefreshToken, Status: fiber.StatusBadRequest}
	}

	_, redisErr := cache.Get[bool](req.RefreshToken)
	if redisErr == nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: service_errors.TokenInvalid, Status: fiber.StatusBadRequest}
	}

	remainingTime := time.Unix(int64((claims[constants.ExpKey]).(float64)), 0)

	redisErr = cache.Set[bool](req.RefreshToken, true, time.Until(remainingTime)*time.Minute)
	if redisErr != nil {
		log.Fatal(redisErr)
		return nil, &service_errors.ServiceErrors{EndUserMessage: "can't black list old refresh token", Status: fiber.StatusInternalServerError}
	}

	key, _ := claims[constants.Key].(string)
	session, _ := claims[constants.SessionIdKey].(string)

	p, err := u.fetchPremiumStatus(key)
	if err != nil {
		return nil, err
	}

	res, e := common.GenerateJwt(&dto.KeyDTO{
		Key:       key,
		SessionId: session,
		Premium:   p,
	}, u.cfg)

	if e != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "JWT generation gone wrong", Status: fiber.StatusInternalServerError}
	}
	return res, nil
}

func (k *KeyService) fetchPremiumStatus(key string) (bool, *service_errors.ServiceErrors) {
	var premium bool
	q := `SELECT premium FROM ac_keys WHERE id = $1`
	err := k.db.QueryRow(q, key).Scan(&premium)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, &service_errors.ServiceErrors{EndUserMessage: "no active key found", Status: fiber.StatusNotFound}
		}
		return false, &service_errors.ServiceErrors{EndUserMessage: "failed to fetch key", Status: fiber.StatusInternalServerError}
	}
	return premium, nil
}
