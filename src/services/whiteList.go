package services

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/api/dto"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/data/postgres"
	"github.com/Arshia-Izadyar/Fast-Gopher/src/pkg/service_errors"
	"github.com/gofiber/fiber/v2"
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

func (wl *WhiteListService) WhiteListRequest(req *dto.WhiteListAddDTO) *service_errors.ServiceErrors {

	// insQ := `
	// INSERT INTO active_devices (session_id, ac_keys_id, ip)
	// VALUES ($1, $2, $3)
	// ON CONFLICT (session_id, ac_keys_id) DO UPDATE
	// SET ip = EXCLUDED.ip;
	// `
	tx, err := wl.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: "starting db transaction failed" + err.Error(), Err: err}
	}

	insQ := `
		UPDATE active_devices
		SET ip = $1
		WHERE session_id = $2 AND ac_keys_id = $3;
	`
	r, err := tx.Exec(insQ, req.UserIp, req.SessionId, req.Key)
	if err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: "INSERT INTO active_devices " + err.Error(), Err: err}
	}

	count, err := r.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		return &service_errors.ServiceErrors{EndUserMessage: "device is not in users active devices please request a session on key", Status: fiber.StatusForbidden}
	}

	if err := exec.Command("ipset", "-!", "add", "whitelist", req.UserIp).Run(); err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: fmt.Sprintf("Attempt 1: Failed to execute ipset command: %v\n", err), Status: fiber.StatusForbidden}
	}
	tx.Commit()

	// pool := cmd.GetPool()
	// pool.Submit(a(req))

	return nil
}
func a(req *dto.WhiteListAddDTO) func() {
	db := postgres.GetDB()

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

	return func() {
		tx, err := db.Begin()
		if err != nil {
			log.Printf("starting db transaction failed: %v\n", err)
			return
		}
		defer tx.Rollback()

		if _, err := tx.Exec(optQ, req.Key); err != nil {
			log.Printf("running the query failed: %v\n", err)
			return
		}

		if err := exec.Command("ipset", "-!", "add", "whitelist", req.UserIp).Run(); err != nil {
			log.Printf("Attempt 1: Failed to execute ipset command: %v\n", err)
			return
		}
		tx.Commit()
	}

}

// ipset add whitelist 192.168.1.1
// ipset del whitelist 192.168.1.1
// iptables -A INPUT -p tcp --dport 443 -m set --match-set whitelist src -j ACCEPT
// iptables -A INPUT -p tcp --dport 443 -j DROP
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

func (wl *WhiteListService) WhiteListRemove(req *dto.WhiteListAddDTO) *service_errors.ServiceErrors {
	// find user
	// del device
	// bye

	tx, err := wl.db.Begin()
	if err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: service_errors.InternalError}
	}

	q := `
	DELETE FROM active_devices where ac_keys_id = $1 AND session_id = $2;
	`

	if _, err = tx.Exec(q, req.Key, req.SessionId); err != nil {
		tx.Rollback()
		fmt.Println(err)
		return &service_errors.ServiceErrors{EndUserMessage: "deletion failed", Status: fiber.StatusInternalServerError}
	}
	if err := exec.Command("ipset", "-!", "del", "whitelist", req.UserIp).Run(); err != nil {
		return &service_errors.ServiceErrors{EndUserMessage: fmt.Sprintf("Attempt 1: Failed to execute ipset command: %v\n", err), Status: fiber.StatusForbidden}
	}

	tx.Commit()
	return nil
}
