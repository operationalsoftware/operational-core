package authmodel

import (
	"app/internal/db"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type AuthUser struct {
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

type VerifyPasswordLoginInput struct {
	Username string
	Password string
}

type VerifyPasswordLoginOutput struct {
	AuthUser      AuthUser
	FailureReason string
}

func VerifyPasswordLogin(ex db.SQLExecutor, input VerifyPasswordLoginInput) (VerifyPasswordLoginOutput, error) {

	INVALID_EMAIL_PASSWORD_MSG := "Invalid email or password. Please try again"
	LOGIN_BLOCKED_MSG := "Login temporarily blocked, please wait and try again"

	out := VerifyPasswordLoginOutput{}

	var authUser AuthUser
	err := ex.QueryRow(`
SELECT
	UserId,
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName,
	Created,
	LastLogin,
	HashedPassword,
	FailedLoginAttempts,
	LoginBlockedUntil
FROM
	User
WHERE
	Username = ?
	`, input.Username).Scan(
		&authUser.UserId,
		&authUser.IsAPIUser,
		&authUser.Username,
		&authUser.Email,
		&authUser.FirstName,
		&authUser.LastName,
		&authUser.Created,
		&authUser.LastLogin,
		&authUser.HashedPassword,
		&authUser.FailedLoginAttempts,
		&authUser.LoginBlockedUntil,
	)

	if err == sql.ErrNoRows {
		out.FailureReason = INVALID_EMAIL_PASSWORD_MSG
		return out, nil
	} else if err != nil {
		return out, err
	}

	// Check if login blocked
	now := time.Now()
	if authUser.LoginBlockedUntil.Valid && now.Before(authUser.LoginBlockedUntil.Time) {
		out.FailureReason = LOGIN_BLOCKED_MSG
		return out, nil
	}

	// Check password
	passwordErr := bcrypt.CompareHashAndPassword([]byte(authUser.HashedPassword), []byte(input.Password))

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword {
		// Update failed login attempts
		authUser.FailedLoginAttempts++

		_, err := ex.Exec(`
UPDATE 
	User 
SET 
	FailedLoginAttempts = ? 
WHERE 
	UserID = ?
		`, authUser.FailedLoginAttempts, authUser.UserId)

		if err != nil {
			log.Panic(err)
		}

		// Check if login blocked
		if authUser.FailedLoginAttempts >= 4 {
			loginBlockedUntil := now.Add(5 * time.Minute)

			updateLoginBlockedUntilStmt := `
UPDATE 
	User
SET 
	LoginBlockedUntil = ?, 
	FailedLoginAttempts = 0
WHERE 
	UserID = ?
`
			_, err := ex.Exec(updateLoginBlockedUntilStmt, loginBlockedUntil, authUser.UserId)

			if err != nil {
				return out, err
			}

			out.FailureReason = LOGIN_BLOCKED_MSG
			return out, nil
		} else {
			out.FailureReason = INVALID_EMAIL_PASSWORD_MSG
			return out, nil
		}

	} else if passwordErr != nil {
		log.Panic(passwordErr)
	}

	// Successful login
	// Update last login
	_, err = ex.Exec(`
UPDATE 
	User 
SET 
	LastLogin = ? 
WHERE 
	UserID = ?
	`, now, authUser.UserId)

	if err != nil {
		log.Fatal(err)
	}

	// This means the user is verified
	out.AuthUser = authUser

	return out, nil
}
