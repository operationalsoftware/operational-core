package authservice

import (
	"app/internal/models"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyPasswordLogin(
		ctx context.Context,
		input models.VerifyPasswordLoginInput,
	) (models.VerifyPasswordLoginOutput, error)
}

type authService struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) AuthService {
	return &authService{db: db}
}

func (s *authService) VerifyPasswordLogin(
	ctx context.Context,
	input models.VerifyPasswordLoginInput,
) (models.VerifyPasswordLoginOutput, error) {

	INVALID_EMAIL_PASSWORD_MSG := "Invalid email or password. Please try again"
	LOGIN_BLOCKED_MSG := "Login temporarily blocked, please wait and try again"

	out := models.VerifyPasswordLoginOutput{}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	var authUserDB models.AuthUserDB
	err = tx.QueryRow(
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
	login_blocked_until
FROM
	app_user
WHERE
	username = $1
	`, input.Username).Scan(
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
	)

	if err == sql.ErrNoRows {
		out.FailureReason = INVALID_EMAIL_PASSWORD_MSG
		return out, nil
	} else if err != nil {
		return out, err
	}

	// Check if login blocked
	now := time.Now()
	if authUserDB.LoginBlockedUntil.Valid && now.Before(authUserDB.LoginBlockedUntil.Time) {
		out.FailureReason = LOGIN_BLOCKED_MSG
		return out, nil
	}

	// Check password
	passwordErr := bcrypt.CompareHashAndPassword([]byte(authUserDB.HashedPassword), []byte(input.Password))

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword {
		// Update failed login attempts
		authUserDB.FailedLoginAttempts++

		_, err := tx.Exec(ctx, `
UPDATE 
	app_user 
SET 
	failed_login_attempts = $1
WHERE 
	user_id = $2
		`, authUserDB.FailedLoginAttempts, authUserDB.UserID)

		if err != nil {
			log.Panic(err)
		}

		// Check if login blocked
		if authUserDB.FailedLoginAttempts >= 4 {
			loginBlockedUntil := now.Add(5 * time.Minute)

			updateLoginBlockedUntilStmt := `
UPDATE 
	app_user
SET 
    login_blocked_until = $1,
    failed_login_attempts = 0
WHERE 
    user_id = $2
`
			_, err := tx.Exec(
				ctx,
				updateLoginBlockedUntilStmt, loginBlockedUntil, authUserDB.UserID,
			)

			if err != nil {
				return out, err
			}

			out.FailureReason = LOGIN_BLOCKED_MSG
		} else {
			out.FailureReason = INVALID_EMAIL_PASSWORD_MSG
		}

		err = tx.Commit(ctx)
		if err != nil {
			return out, err
		}

		return out, nil
	} else if passwordErr != nil {
		log.Panic(passwordErr)
	}

	// Successful login
	// Update last login
	_, err = tx.Exec(
		ctx,
		`
UPDATE 
	app_user 
SET 
	last_login = $1
WHERE 
	user_id = $2
	`, now, authUserDB.UserID)

	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return out, err
	}

	// This means the user is verified
	out.AuthUser = authUserDB.ToDomain()

	return out, nil
}
