package usermodel

import (
	"app/internal/appsort"
	"app/internal/db"
	"app/internal/validate"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID      int
	Username    string
	IsAPIUser   bool
	Email       sql.NullString
	FirstName   sql.NullString
	LastName    sql.NullString
	Created     time.Time
	LastLogin   sql.NullTime
	Permissions UserPermissions
}

type NewUser struct {
	Username        string
	Email           sql.NullString
	FirstName       string
	LastName        string
	Password        string
	ConfirmPassword string
	Permissions     UserPermissions
}

func validateUsername(ve *validate.ValidationErrors, username string) {
	pattern := "^[a-z0-9_]+$"

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
	}
	if !re.MatchString(username) {
		ve.Add("Username", "must contain only letters, numbers, and underscores")
	}
}

func ValidateNewUser(user NewUser) (bool, validate.ValidationErrors) {

	var validationErrors validate.ValidationErrors = make(map[string][]string)

	// validate Username
	validate.MinLength(&validationErrors, "Username", user.Username, 3)
	validate.MaxLength(&validationErrors, "Username", user.Username, 20)
	validate.Lowercase(&validationErrors, "Username", user.Username)
	validateUsername(&validationErrors, user.Username)
	// check if username is already taken
	db := db.UseDB()
	if user.Username != "" {
		_, err := ByUsername(db, user.Username)
		if err == nil {
			validationErrors.Add("Username", fmt.Sprintf("'%s' is already taken", user.Username))
		}
	}

	// validate FirstName
	validate.MinLength(&validationErrors, "FirstName", user.FirstName, 1)
	validate.MaxLength(&validationErrors, "FirstName", user.FirstName, 20)

	// validate LastName
	validate.MinLength(&validationErrors, "LastName", user.LastName, 1)
	validate.MaxLength(&validationErrors, "LastName", user.LastName, 20)

	// validate Email
	if user.Email.String != "" {
		validate.Email(&validationErrors, "Email", user.Email.String)
	}

	// validate Password
	validate.Password(&validationErrors, "Password", user.Password)

	// validate confirm password
	if user.Password != user.ConfirmPassword {
		validationErrors.Add("ConfirmPassword", "does not match")
	}

	// user permissions don't need to be validated, the struct will be populated with
	// matching permissions from the form data meaning any missing data will be
	// zero-valued (boolean, false) and any extra data will be ignored

	return len(validationErrors) == 0, validationErrors
}

func Add(db db.SQLExecutor, user NewUser) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	permissionsJson, err := json.Marshal(user.Permissions)
	if err != nil {
		return err
	}

	insertUserStmt := `
INSERT INTO User (
	Username,
	IsAPIUser,
	Email,
	FirstName,
	LastName,
	HashedPassword,
	Permissions
)
VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = db.Exec(
		insertUserStmt,
		user.Username,
		false,
		user.Email,
		user.FirstName,
		user.LastName,
		string(hashedPassword),
		permissionsJson,
	)

	if err != nil {
		return err
	}

	return nil
}

type NewAPIUser struct {
	Username    string
	Password    string
	Permissions UserPermissions
}

func ValidateNewAPIUser(user NewAPIUser) (bool, validate.ValidationErrors) {
	var validationErrors validate.ValidationErrors = make(map[string][]string)

	// validate Username
	validate.MinLength(&validationErrors, "Username", user.Username, 3)
	validate.MaxLength(&validationErrors, "Username", user.Username, 20)
	validate.Lowercase(&validationErrors, "Username", user.Username)
	validateUsername(&validationErrors, user.Username)
	// check if username is already taken
	db := db.UseDB()
	if user.Username != "" {
		_, err := ByUsername(db, user.Username)
		if err == nil {
			validationErrors.Add("Username", fmt.Sprintf("'%s' is already taken", user.Username))
		}
	}

	return len(validationErrors) == 0, validationErrors
}

func AddAPIUser(db db.SQLExecutor, user NewAPIUser) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	permissionsJson, err := json.Marshal(user.Permissions)
	if err != nil {
		return err
	}

	insertUserStmt := `
INSERT INTO User (
	Username,
	IsAPIUser,
	HashedPassword,
	Permissions
)
VALUES (?, ?, ?, ?)
`

	_, err = db.Exec(
		insertUserStmt,
		user.Username,
		true,
		string(hashedPassword),
		permissionsJson,
	)

	if err != nil {
		return err
	}

	return nil
}

type UserUpdate struct {
	Username    string
	Email       sql.NullString
	FirstName   string
	LastName    string
	Permissions UserPermissions
}

func ValidateUserUpdate(update UserUpdate) (bool, validate.ValidationErrors) {

	var validationErrors validate.ValidationErrors = make(map[string][]string)

	validate.MinLength(&validationErrors, "Username", update.Username, 3)
	validate.MaxLength(&validationErrors, "Username", update.Username, 20)
	validate.Lowercase(&validationErrors, "Username", update.Username)
	validateUsername(&validationErrors, update.Username)

	// validate FirstName
	validate.MinLength(&validationErrors, "FirstName", update.FirstName, 1)
	validate.MaxLength(&validationErrors, "FirstName", update.FirstName, 20)

	// validate LastName
	validate.MinLength(&validationErrors, "LastName", update.LastName, 1)
	validate.MaxLength(&validationErrors, "LastName", update.LastName, 20)

	// validate Email
	if update.Email.String != "" {
		validate.Email(&validationErrors, "Email", update.Email.String)
	}

	// user permissions don't need to be validated. See description in ValidateNewUser

	return len(validationErrors) == 0, validationErrors
}

func Update(db db.SQLExecutor, id int, update UserUpdate) error {

	// get the user to check if it exists
	user, err := ByID(db, id)
	if err != nil {
		return err
	}

	// can't update an API user using this method
	if user.IsAPIUser {
		return fmt.Errorf("API users cannot be updated using this method")
	}

	query := `
UPDATE
	User

SET
	FirstName = ?,
	LastName = ?,
	Email = ?,
	Username = ?,
	Permissions = ?

WHERE
	UserID = ?
	`

	permissionsJSON, err := json.Marshal(update.Permissions)

	if err != nil {
		return err
	}

	_, err = db.Exec(
		query,

		update.FirstName,
		update.LastName,
		update.Email,
		update.Username,
		string(permissionsJSON),
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

type PasswordReset struct {
	Password        string
	ConfirmPassword string
}

func ValidatePasswordReset(pr PasswordReset) (bool, validate.ValidationErrors) {

	var validationErrors validate.ValidationErrors = make(map[string][]string)

	// validate Password
	validate.Password(&validationErrors, "Password", pr.Password)

	// validate confirm password
	if pr.Password != pr.ConfirmPassword {
		validationErrors.Add("ConfirmPassword", "does not match")
	}

	return len(validationErrors) == 0, validationErrors
}

func ResetPassword(db db.SQLExecutor, id int, pr PasswordReset) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
UPDATE
	User
SET
	HashedPassword = ?
WHERE
	UserID = ?
	`

	_, err = db.Exec(query, string(hashedPassword), id)
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
	Created,
	LastLogin,
	Permissions
FROM
	User
WHERE
	UserID = ?
	`

	var user User
	var permissionsJSON string
	err := db.QueryRow(query, id).Scan(
		&user.UserID,
		&user.IsAPIUser,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Created,
		&user.LastLogin,
		&permissionsJSON,
	)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("User not found")
	} else if err != nil {
		log.Panic(err)
	}

	// Unmarshal the permissionsJSON string into a []string
	err = json.Unmarshal([]byte(permissionsJSON), &user.Permissions)
	if err != nil {
		log.Panic(err)
	}

	return user, nil
}

func ByUsername(db db.SQLExecutor, username string) (User, error) {
	query := `
SELECT
	UserID,
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName,
	Created,
	LastLogin,
	Permissions
FROM
	User
WHERE
	Username = ?
	`

	var user User
	var permissionsJSON string
	err := db.QueryRow(query, username).Scan(
		&user.UserID,
		&user.IsAPIUser,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Created,
		&user.LastLogin,
		&permissionsJSON,
	)
	if err == sql.ErrNoRows {
		return user, fmt.Errorf("User not found")
	} else if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal([]byte(permissionsJSON), &user.Permissions)
	if err != nil {
		log.Panic(err)
	}

	return user, nil
}

var ListSortableKeys = []string{
	"Username",
	"Email",
	"FirstName",
	"LastName",
	"Created",
	"LastLogin",
}

type ListQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int
}

func List(db db.SQLExecutor, q ListQuery) ([]User, error) {

	offset := (q.Page - 1) * q.PageSize
	limit := q.PageSize
	orderByClause := q.Sort.ToOrderByClause(map[string]string{})

	query := fmt.Sprintf(`
SELECT
	UserID,
	IsAPIUser,
	Username,
	Email,
	FirstName,
	LastName,
	Created,
	LastLogin,
	Permissions
FROM
	User

%s

LIMIT ? OFFSET ?
	`,
		orderByClause,
	)

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		var permissionsJSON string
		err := rows.Scan(
			&u.UserID,
			&u.IsAPIUser,
			&u.Username,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Created,
			&u.LastLogin,
			&permissionsJSON,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(permissionsJSON), &u.Permissions)
		if err != nil {
			log.Panic(err)
		}

		users = append(users, u)
	}

	return users, nil
}

func Count(db db.SQLExecutor) (int, error) {
	query := `
SELECT
	COUNT(*)
FROM
	User
	`

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
