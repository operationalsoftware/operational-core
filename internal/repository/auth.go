package repository

import (
	"app/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (r *AuthRepository) GetAuthUserByUsername(
	ctx context.Context,
	tx pgx.Tx,
	username string,
) (*model.AuthUser, error) {

	var authUserDB model.AuthUserDB
	var authUser model.AuthUser

	err := tx.QueryRow(
		ctx,
		`
SELECT
	user_id,
	is_api_user,
	username,
	email,
	first_name,
	last_name,
	created,
	last_login,
	hashed_password,
	failed_login_attempts,
	login_blocked_until,
	session_duration_minutes
FROM
	app_user
WHERE
	username = $1
	`, username).Scan(
		&authUserDB.UserID,
		&authUserDB.IsAPIUser,
		&authUserDB.Username,
		&authUserDB.Email,
		&authUserDB.FirstName,
		&authUserDB.LastName,
		&authUserDB.Created,
		&authUserDB.LastLogin,
		&authUserDB.HashedPassword,
		&authUserDB.FailedLoginAttempts,
		&authUserDB.LoginBlockedUntil,
		&authUserDB.SessionDurationMinutes,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return &authUser, err
	}

	authUser = authUserDB.ToDomain()

	return &authUser, nil
}

func (r *AuthRepository) IncrementFailedLoginAttempts(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
) error {

	query := `
UPDATE 
	app_user 
SET 
	failed_login_attempts = failed_login_attempts + 1
WHERE 
	user_id = $1
	`

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) SetLoginBlockedUntil(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	loginBlockedUntil *time.Time,
) error {

	var dbTime = pgtype.Timestamptz{}

	if loginBlockedUntil != nil {
		dbTime.Valid = true
		dbTime.Time = *loginBlockedUntil
	}

	query := `
UPDATE 
	app_user
SET 
    login_blocked_until = $1,
    failed_login_attempts = 0
WHERE 
    user_id = $2
	`

	_, err := tx.Exec(ctx, query, dbTime, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) UpdateLastLogin(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
) error {

	query := `
UPDATE 
	app_user 
SET 
	last_login = NOW()
WHERE 
	user_id = $1
	`

	_, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}
