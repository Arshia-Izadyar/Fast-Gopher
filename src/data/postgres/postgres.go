package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Arshia-Izadyar/Fast-Gopher/src/config"
	_ "github.com/lib/pq"
)

var psqlDb *sql.DB

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	var err error
	fmt.Println(cfg.Postgres.DbName)
	fmt.Println(cfg.Postgres.User)
	cnn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DbName,
	)

	db, err := sql.Open("postgres", cnn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	psqlDb = db
	return db, nil
}

func GetDB() *sql.DB {
	return psqlDb
}

func CloseDB() error {
	err := psqlDb.Close()
	if err != nil {
		return err
	}
	return nil
}
