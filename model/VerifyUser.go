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

	checkUserWithUsernameStmt := `SELECT * FROM users WHERE username = ?`

	var user User

	row := dbInstance.QueryRow(checkUserWithUsernameStmt, username)

	err := row.Scan(&user.UserId, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Created, &user.LastLogin, &user.HashedPassword, &user.FailedLoginAttempts, &user.LoginBlockedUntil)

	if err == sql.ErrNoRows {
		return User{}, fmt.Errorf("incorrect username or password")
	}

	// Check if login blocked
	if user.LoginBlockedUntil != (time.Time{}) {
		now := time.Now()

		if now.Before(user.LoginBlockedUntil) {
			return User{}, fmt.Errorf("login blocked until %s", user.LoginBlockedUntil.Format("15:04:05"))
		}
	}

	// Check password
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword {
		// Update failed login attempts
		user.FailedLoginAttempts++

		updateFailedLoginAttemptsStmt := `
UPDATE 
	users 
SET 
	failed_login_attempts = ? 
WHERE 
	user_id = ?
`
		_, err := dbInstance.Exec(updateFailedLoginAttemptsStmt, user.FailedLoginAttempts, user.UserId)

		if err != nil {
			log.Fatal(err)
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
		log.Fatal(passwordErr)
	}

	// Successfull login
	// Update last login
	now := time.Now()
	updateLastLoginStmt := `
UPDATE 
	users 
SET 
	last_login = ? 
WHERE 
	user_id = ?
`
	_, err = dbInstance.Exec(updateLastLoginStmt, now, user.UserId)

	if err != nil {
		log.Fatal(err)
	}

	// This means the user is verified
	return user, nil

}
