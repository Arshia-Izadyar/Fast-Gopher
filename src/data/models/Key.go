package models

import (
	"time"
)

type UserAttrs struct {
}

type Key struct {
	ID      string `db:"id"`
	Premium bool   `db:"premium"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ActiveSessions struct {
	ID        int    `db:"id"`
	SessionId string `db:"session_id"`
	Ip        string `db:"ip"`
	Key       string `db:"key"`

	CreatedAt time.Time `db:"created_at"`
}
