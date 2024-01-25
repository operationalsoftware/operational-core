package model

import "database/sql"

func EditUser(db *sql.DB, user User, id int) error {
	editUserQuery := `
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
		editUserQuery,

		user.FirstName,
		user.LastName,
		user.Email,
		user.Username,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
