package model

import (
	"app/pkg/pgconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type AuthUserDB struct {
	UserID              int                `db:"user_id"`
	Username            string             `db:"username"`
	IsAPIUser           bool               `db:"is_api_user"`
	Email               pgtype.Text        `db:"email"`
	FirstName           pgtype.Text        `db:"first_name"`
	LastName            pgtype.Text        `db:"last_name"`
	Created             time.Time          `db:"created"`
	LastLogin           pgtype.Timestamptz `db:"last_login"`
	HashedPassword      string             `db:"hashed_password"`
	FailedLoginAttempts int                `db:"failed_login_attempts"`
	LoginBlockedUntil   pgtype.Timestamptz `db:"login_blocked_until"`
}

type AuthUser struct {
	UserID              int
	Username            string
	IsAPIUser           bool
	Email               *string
	FirstName           *string
	LastName            *string
	Created             time.Time
	LastLogin           *time.Time
	HashedPassword      string
	FailedLoginAttempts int
	LoginBlockedUntil   *time.Time
}

func (a AuthUserDB) ToDomain() AuthUser {
	return AuthUser{
		UserID:              a.UserID,
		Username:            a.Username,
		IsAPIUser:           a.IsAPIUser,
		Email:               pgconv.PGTextToStringPtr(a.Email),
		FirstName:           pgconv.PGTextToStringPtr(a.FirstName),
		LastName:            pgconv.PGTextToStringPtr(a.LastName),
		Created:             a.Created,
		LastLogin:           pgconv.PGTimestamptzToTimePtr(a.LastLogin),
		HashedPassword:      a.HashedPassword,
		FailedLoginAttempts: a.FailedLoginAttempts,
		LoginBlockedUntil:   pgconv.PGTimestamptzToTimePtr(a.LoginBlockedUntil),
	}
}

type VerifyPasswordLoginInput struct {
	Username string
	Password string
}

type VerifyPasswordLoginOutput struct {
	AuthUser      AuthUser
	FailureReason string
}
