package services

import (
	"database/sql"
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
)

/*
 * 1. user request to be whitelisted
 * 2. immediately add user to whitelist
 * 3. check for user ip count
 * 4. if user had more than 5 ips remove the oldest ip
 * 5. else let user use service
 */

type WhiteListService struct {
	db  *sql.DB
	cfg *config.Config
}

func NewWhiteListService(cfg *config.Config) *WhiteListService {
	db := postgres.GetDB()
	return &WhiteListService{
		db:  db,
		cfg: cfg,
	}
}

func (wl *WhiteListService) WhiteListRequest(req *dto.WhiteListAddDTO) error {

	tx, err := wl.db.Begin()
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: "first " + err.Error(), Err: err}
	}

	insQ := `
	INSERT INTO active_devices (device_id, user_id, ips) 
    VALUES ($1, $2, $3)
    ON CONFLICT (device_id, user_id) DO UPDATE
    SET ips = EXCLUDED.ips;
	`
	_, err = tx.Exec(insQ, req.UserDeviceID, req.UserId, req.UserIp)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return &service_errors.ServiceError{EndUserMessage: "INSERT INTO active_devices " + err.Error(), Err: err}
	}
	go func() {
		wl.whiteListAdd(req) // run in background
	}()
	tx.Commit()
	return nil

}

func (wl *WhiteListService) whiteListAdd(req *dto.WhiteListAddDTO) error {
	userId := req.UserId
	tx, err := wl.db.Begin()
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: "tx : " + err.Error(), Err: err}
	}

	var count int
	countQ := `
	SELECT COUNT(*) FROM active_devices WHERE user_id = $1;
	`
	err = wl.db.QueryRow(countQ, userId).Scan(&count)
	if err != nil {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: "count : " + err.Error(), Err: err}
	}

	if count > 5 {
		rmQ := `
		DELETE FROM active_devices
			WHERE id = (
				SELECT id FROM active_devices
				WHERE user_id = $1
				ORDER BY created_at ASC
				LIMIT 1
			)`
		if _, err := wl.db.Exec(rmQ, userId); err != nil {
			tx.Rollback()
			return &service_errors.ServiceError{EndUserMessage: "c : " + err.Error(), Err: err}
		}
		tx.Commit()
	}
	return nil

}

func (wl *WhiteListService) WhiteListRemove(req *dto.WhiteListAddDTO) error {
	// find user
	// del device
	// bye

	tx, err := wl.db.Begin()
	if err != nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.InternalError}
	}

	q := `
	DELETE FROM active_devices where user_id = $1 AND device_id = $2;
	`

	if _, err = tx.Exec(q, req.UserId, req.UserDeviceID); err != nil {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: "deletion failed"}
	}

	tx.Commit()
	return nil
}
