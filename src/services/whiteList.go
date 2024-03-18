package services

import (
	"database/sql"
	"fmt"
	"os/exec"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/cmd/cmd"
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

	insQ := `
	INSERT INTO active_devices (session_id, ac_keys_id, ip) 
    VALUES ($1, $2, $3)
    ON CONFLICT (session_id, ac_keys_id) DO UPDATE
    SET ip = EXCLUDED.ip;
	`
	if _, err := wl.db.Exec(insQ, req.SessionId, req.Key, req.UserIp); err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: "INSERT INTO active_devices " + err.Error(), Err: err}
	}

	pool := cmd.GetPool()
	pool.Submit(a(req))
	// go func() {
	// 	wl.whiteListAdd(req) // run in background
	// }()
	return nil
}
func a(req *dto.WhiteListAddDTO) func() {
	return func() {

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
		db := postgres.GetDB()
		if _, err := db.Exec(optQ, req.Key); err != nil {
			return
		}

		bt, _ := exec.Command("ping", "8.8.8.8").Output()
		fmt.Println(string(bt))
	}

}

// func (wl *WhiteListService) whiteListAdd(req *dto.WhiteListAddDTO) error {
// 	userId := req.UserId

// 	optQ := `
// 	WITH ranked_devices AS (
// 		SELECT id, ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at ASC) AS rn
// 		FROM active_devices
// 		WHERE user_id = $1
// 	)
// 	DELETE FROM active_devices
// 	WHERE id IN (
// 		SELECT id FROM ranked_devices WHERE rn > 5
// 	);
// 	`
// 	if _, err := wl.db.Exec(optQ, userId); err != nil {
// 		return &service_errors.ServiceErrors{EndUserMessage: "Optimized deletion error: " + err.Error(), Err: err}
// 	}
// 	return nil
// }

// func (wl *WhiteListService) WhiteListRemove(req *dto.WhiteListAddDTO) error {
// 	// find user
// 	// del device
// 	// bye

// 	tx, err := wl.db.Begin()
// 	if err != nil {
// 		return &service_errors.ServiceErrors{EndUserMessage: service_errors.InternalError}
// 	}

// 	q := `
// 	DELETE FROM active_devices where user_id = $1 AND device_id = $2;
// 	`

// 	if _, err = tx.Exec(q, req.UserId, req.UserDeviceID); err != nil {
// 		tx.Rollback()
// 		return &service_errors.ServiceErrors{EndUserMessage: "deletion failed"}
// 	}

// 	tx.Commit()
// 	return nil
// }
