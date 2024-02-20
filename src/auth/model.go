package auth

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"app/db"

	"golang.org/x/crypto/bcrypt"
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

func VerifyUser(username, password string) (AuthUser, error) {
	// Query db for user
	dbInstance := db.UseDB()

	var authUser AuthUser
	err := dbInstance.QueryRow(`
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
	`, username).Scan(
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
		return AuthUser{}, fmt.Errorf("incorrect username or password")
	} else if err != nil {
		log.Panic(err)
	}

	// Check if login blocked
	if authUser.LoginBlockedUntil.Valid {
		now := time.Now()

		if now.Before(authUser.LoginBlockedUntil.Time) {
			return AuthUser{}, fmt.Errorf("login blocked until %s", authUser.LoginBlockedUntil.Time.Format("15:04:05"))
		}
	}

	// Check password
	passwordErr := bcrypt.CompareHashAndPassword([]byte(authUser.HashedPassword), []byte(password))

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword {
		// Update failed login attempts
		authUser.FailedLoginAttempts++

		_, err := dbInstance.Exec(`
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
			now := time.Now()
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
			_, err := dbInstance.Exec(updateLoginBlockedUntilStmt, loginBlockedUntil, authUser.UserId)

			if err != nil {
				log.Fatal(err)
			}

			return AuthUser{}, fmt.Errorf("login blocked until %s", loginBlockedUntil.Format("15:04:05"))
		} else {
			return AuthUser{}, fmt.Errorf("incorrect username or password")
		}

	} else if passwordErr != nil {
		log.Panic(passwordErr)
	}

	// Successfull login
	// Update last login
	now := time.Now()
	_, err = dbInstance.Exec(`
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
	return authUser, nil
}
