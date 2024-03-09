package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAttrs struct {
}

type User struct {
	ID           uuid.UUID `db:"id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Email        string    `db:"email"`
	UserPassword string    `db:"user_password"`
	UserType     string    `db:"user_type"`
	UserAttrs    UserAttrs `db:"user_attrs"`
}

type ActiveDevices struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UserId    uuid.UUID `db:"user_id"`
	DeviceID  string    `db:"device_id"`
}
