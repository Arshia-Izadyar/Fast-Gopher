package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	var err error
	cnn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DbName,
	)

	db, err = sql.Open("postgres", cnn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}