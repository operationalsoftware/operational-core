package model

import (
	"database/sql"
	"time"
)

type User struct {
	UserId              int
	Username            string
	IsAPIUser           bool
	Email               sql.NullString
	FirstName           sql.NullString
	LastName            sql.NullString
	Created             time.Time
	LastLogin           sql.NullTime
	HashedPassword      string
	FailedLoginAttempts int
	LoginBlockedUntil   sql.NullTime
}
