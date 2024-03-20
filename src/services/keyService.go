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
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type KeyService struct {
	db  *sql.DB
	cfg *config.Config
}

func NewKeyService(cfg *config.Config) *KeyService {
	db := postgres.GetDB()
	return &KeyService{
		db:  db,
		cfg: cfg,
	}
}

// service for /key
func (u *KeyService) GenerateKey(req *dto.GenerateKeyDTO) (*dto.KeyAcDTO, *service_errors.ServiceErrors) {
	if req.DeviceName == "" {
		req.DeviceName = "unknown-device"
	}
	tx, err := u.db.Begin()
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "database transaction start error: " + err.Error(), Status: fiber.StatusInternalServerError, Err: err}
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
			return nil, &service_errors.ServiceErrors{EndUserMessage: "the key already exists try getting another key", Status: fiber.StatusForbidden}
		}
		return nil, &service_errors.ServiceErrors{EndUserMessage: "insert query failed" + err.Error(), Status: fiber.StatusInternalServerError}
	}

	saveSession := `
		INSERT INTO active_devices(session_id, ac_keys_id, device_name) VALUES($1, $2, $3);
	`
	sessionId := uuid.New()
	_, err = tx.Exec(saveSession, sessionId.String(), key, req.DeviceName)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, &service_errors.ServiceErrors{EndUserMessage: fmt.Sprintf("session id (%s) already exists", sessionId.String()), Status: fiber.StatusBadRequest}
		}
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to commit transaction", Status: fiber.StatusInternalServerError}
	}

	result, err := common.GenerateJwt(key, sessionId.String(), u.cfg)
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

	// Save session
	saveSession := `INSERT INTO active_devices(session_id, ac_keys_id, device_name) VALUES($1, $2, $3) ON CONFLICT (session_id, ac_keys_id) DO NOTHING;`
	sessionId := uuid.New()

	_, err = tx.Exec(saveSession, sessionId.String(), req.Key, req.DeviceName)
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to save session", Status: fiber.StatusInternalServerError}
	}

	optQ := `
		WITH ranked_devices AS (
			SELECT id, ROW_NUMBER() OVER (PARTITION BY ac_keys_id ORDER BY created_at DESC) AS rn
			FROM active_devices
			WHERE ac_keys_id = $1
		)
		DELETE FROM active_devices
		WHERE id IN (
			SELECT id FROM ranked_devices WHERE rn > 5
		);
	`
	_, err = tx.Exec(optQ, req.Key)
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to remove the 6th device", Status: fiber.StatusInternalServerError}
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to commit transaction", Status: fiber.StatusInternalServerError}
	}

	// Generate JWT
	tk, err := common.GenerateJwt(req.Key, sessionId.String(), u.cfg)
	if err != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "failed to generate JWT", Status: fiber.StatusInternalServerError}
	}
	return tk, nil
}

// refresh
func (u *KeyService) Refresh(req *dto.RefreshTokenDTO) (*dto.KeyAcDTO, *service_errors.ServiceErrors) {
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

	redisErr = cache.Set[bool](req.RefreshToken, true, time.Until(time.Unix(int64((claims[constants.ExpKey]).(float64)), 0))*time.Minute)
	if redisErr != nil {
		log.Fatal(redisErr)
		return nil, &service_errors.ServiceErrors{EndUserMessage: "can't black list old refresh token", Status: fiber.StatusInternalServerError}
	}

	key, _ := claims[constants.Key].(string)
	session, _ := claims[constants.SessionIdKey].(string)

	res, e := common.GenerateJwt(key, session, u.cfg)

	if e != nil {
		return nil, &service_errors.ServiceErrors{EndUserMessage: "JWT generation gone wrong", Status: fiber.StatusInternalServerError}
	}
	return res, nil
}

func (k *KeyService) ShowAllActiveDevices(req *dto.IKeyDTO) ([]dto.DeviceDTO, *service_errors.ServiceErrors) {
	q := `
		SELECT device_name, session_id, ip FROM active_devices WHERE ac_keys_id = $1
	`
	r, err := k.db.Query(q, req.Key)
	if err != nil {
		return nil, nil
	}
	var res []dto.DeviceDTO
	for r.Next() {

		var device dto.DeviceDTO
		r.Scan(&device.DeviceName, &device.SessionId, &device.Ip)
		res = append(res, device)
	}

	return res, nil
}

func (k *KeyService) DeleteSession(req *dto.RemoveDeviceDTO) *service_errors.ServiceErrors {
	q := `
		DELETE FROM active_devices WHERE session_id = $1 AND device_name = $2;
	`
	r, err := k.db.Exec(q, req.SessionId, req.DeviceName)
	if count, _ := r.RowsAffected(); count == 0 {
		return &service_errors.ServiceErrors{EndUserMessage: "device not found", Status: fiber.StatusNotFound}
	}

	if err != nil {
		log.Fatal(err)
		return &service_errors.ServiceErrors{EndUserMessage: "can't delete device something went wrong", Status: fiber.StatusInternalServerError}
	}
	return nil
}

func (k *KeyService) DeleteAllSessions(req *dto.SessionKeyDTO) *service_errors.ServiceErrors {
	tx, err := k.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: "failed to start db transaction", Status: fiber.StatusInternalServerError}
	}

	q := `
		DELETE FROM active_devices WHERE ac_keys_id = $1 AND session_id != $2;
	`
	r, err := tx.Exec(q, req.Key, req.SessionId)
	if count, _ := r.RowsAffected(); count == 0 {
		return &service_errors.ServiceErrors{EndUserMessage: "device not found", Status: fiber.StatusNotFound}
	}

	if err != nil {
		log.Fatal(err)
		return &service_errors.ServiceErrors{EndUserMessage: "can't delete devices something went wrong", Status: fiber.StatusInternalServerError}
	}
	tx.Commit()
	return nil
}
