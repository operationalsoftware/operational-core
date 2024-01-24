package model

import (
	"database/sql"
	"operationalcore/db"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserToAdd struct {
	Username  string
	IsAPIUser bool
	Email     sql.NullString
	FirstName sql.NullString
	LastName  sql.NullString
	Password  string
}

func AddUser(db db.SQLExecutor, user UserToAdd) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	insertUserStmt := `
INSERT INTO User (
	Username,
	IsAPIUser,
	Email,
	FirstName,
	LastName,
	HashedPassword
)
VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(
		insertUserStmt,
		user.Username,
		user.IsAPIUser,
		user.Email,
		user.FirstName,
		user.LastName,
		string(hashedPassword),
	)

	if err != nil {
		return err
	}

	return nil
}
