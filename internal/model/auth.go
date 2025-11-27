package model

import (
	"app/pkg/pgconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type AuthUserDB struct {
	UserID                 int                `db:"user_id"`
	Username               string             `db:"username"`
	IsAPIUser              bool               `db:"is_api_user"`
	Email                  pgtype.Text        `db:"email"`
	FirstName              pgtype.Text        `db:"first_name"`
	LastName               pgtype.Text        `db:"last_name"`
	Created                time.Time          `db:"created"`
	LastLogin              pgtype.Timestamptz `db:"last_login"`
	LastActive             pgtype.Timestamptz `db:"last_active"`
	HashedPassword         string             `db:"hashed_password"`
	FailedLoginAttempts    int                `db:"failed_login_attempts"`
	LoginBlockedUntil      pgtype.Timestamptz `db:"login_blocked_until"`
	SessionDurationMinutes pgtype.Int4        `db:"session_duration_minutes"`
}

type AuthUser struct {
	UserID                 int
	Username               string
	IsAPIUser              bool
	Email                  *string
	FirstName              *string
	LastName               *string
	Created                time.Time
	LastLogin              *time.Time
	LastActive             *time.Time
	HashedPassword         string
	FailedLoginAttempts    int
	LoginBlockedUntil      *time.Time
	SessionDurationMinutes *int
}

func (a AuthUserDB) ToDomain() AuthUser {
	return AuthUser{
		UserID:                 a.UserID,
		Username:               a.Username,
		IsAPIUser:              a.IsAPIUser,
		Email:                  pgconv.PGTextToStringPtr(a.Email),
		FirstName:              pgconv.PGTextToStringPtr(a.FirstName),
		LastName:               pgconv.PGTextToStringPtr(a.LastName),
		Created:                a.Created,
		LastLogin:              pgconv.PGTimestamptzToTimePtr(a.LastLogin),
		LastActive:             pgconv.PGTimestamptzToTimePtr(a.LastActive),
		HashedPassword:         a.HashedPassword,
		FailedLoginAttempts:    a.FailedLoginAttempts,
		LoginBlockedUntil:      pgconv.PGTimestamptzToTimePtr(a.LoginBlockedUntil),
		SessionDurationMinutes: pgconv.PGInt4ToIntPtr(a.SessionDurationMinutes),
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
