package userModel

import (
	"database/sql"
	"fmt"
	"log"
	"operationalcore/db"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID    int
	Username  string
	IsAPIUser bool
	Email     sql.NullString
	FirstName sql.NullString
	LastName  sql.NullString
	Created   time.Time
	LastLogin sql.NullTime
}

type NewUser struct {
	Username  string
	IsAPIUser bool
	Email     sql.NullString
	FirstName sql.NullString
	LastName  sql.NullString
	Password  string
}

func Add(db db.SQLExecutor, user NewUser) error {

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

type UserUpdate struct {
	Username  string
	Email     sql.NullString
	FirstName sql.NullString
	LastName  sql.NullString
}

func Update(db db.SQLExecutor, id int, update UserUpdate) error {
	query := `
UPDATE
	User

SET
	FirstName = ?,
	LastName = ?,
	Email = ?,
	Username = ?

WHERE
	UserID = ?
	`

	_, err := db.Exec(
		query,

		update.FirstName,
		update.LastName,
		update.Email,
		update.Username,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func ByID(db db.SQLExecutor, id int) (User, error) {
	query := `
SELECT
	UserID,
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName,
	Created
FROM
	User
WHERE
	UserID = ?
	`

	var user User
	err := db.QueryRow(query, id).Scan(
		&user.UserID,
		&user.IsAPIUser,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Created,
	)

	if err == sql.ErrNoRows {
		return user, fmt.Errorf("User not found")
	} else if err != nil {
		log.Panic(err)
	}

	return user, nil
}

func List(db db.SQLExecutor) ([]User, error) {
	query := `
SELECT
	UserID,
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName,
	Created
FROM
	User

ORDER BY
	Username ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.UserID,
			&u.IsAPIUser,
			&u.Username,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Created,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
