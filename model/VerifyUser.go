package model

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"operationalcore/db"

	"golang.org/x/crypto/bcrypt"
)

func VerifyUser(username, password string) (User, error) {
	// Query db for user
	dbInstance := db.UseDB()

	var user User
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
		&user.UserId,
		&user.IsAPIUser,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Created,
		&user.LastLogin,
		&user.HashedPassword,
		&user.FailedLoginAttempts,
		&user.LoginBlockedUntil,
	)

	if err == sql.ErrNoRows {
		return User{}, fmt.Errorf("incorrect username or password")
	} else if err != nil {
		log.Panic(err)
	}

	// Check if login blocked
	if user.LoginBlockedUntil.Valid {
		now := time.Now()

		if now.Before(user.LoginBlockedUntil.Time) {
			return User{}, fmt.Errorf("login blocked until %s", user.LoginBlockedUntil.Time.Format("15:04:05"))
		}
	}

	// Check password
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword {
		// Update failed login attempts
		user.FailedLoginAttempts++

		_, err := dbInstance.Exec(`
UPDATE 
	User 
SET 
	FailedLoginAttempts = ? 
WHERE 
	UserID = ?
		`, user.FailedLoginAttempts, user.UserId)

		if err != nil {
			log.Panic(err)
		}

		// Check if login blocked
		if user.FailedLoginAttempts >= 4 {
			now := time.Now()
			loginBlockedUntil := now.Add(5 * time.Minute)

			updateLoginBlockedUntilStmt := `
UPDATE 
	users 
SET 
	login_blocked_until = ?, 
	failed_login_attempts = 0
WHERE 
	user_id = ?
`
			_, err := dbInstance.Exec(updateLoginBlockedUntilStmt, loginBlockedUntil, user.UserId)

			if err != nil {
				log.Fatal(err)
			}

			return User{}, fmt.Errorf("login blocked until %s", loginBlockedUntil.Format("15:04:05"))
		} else {
			return User{}, fmt.Errorf("incorrect username or password")
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
	`, now, user.UserId)

	if err != nil {
		log.Fatal(err)
	}

	// This means the user is verified
	return user, nil
}
