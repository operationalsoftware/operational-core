package model

import (
	"database/sql"
	"log"
)

func GetUsers(db *sql.DB) []User {

	users := []User{}
	getUsersQuery := `
SELECT
	UserID,
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName,
	Created,
	LastLogin

FROM
	User

ORDER BY
	Username ASC
	`
	rows, err := db.Query(getUsersQuery)

	if err != nil {
		log.Panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.UserId,
			&user.IsAPIUser,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Created,
			&user.LastLogin,
		)
		if err != nil {
			log.Panic(err)
		}
		users = append(users, user)
	}

	err = rows.Err()

	if err != nil {
		log.Panic(err)
	}

	return users
}
