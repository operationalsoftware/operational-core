package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	exec db.PGExecutor,
	user model.NewUser,
) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	permissionsJson, err := json.Marshal(user.Permissions)
	if err != nil {
		return err
	}

	insertUserStmt := `
INSERT INTO app_user (
	username,
	is_api_user,
	email,
	first_name,
	last_name,
	hashed_password,
	permissions
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = exec.Exec(
		ctx,

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

func (r *UserRepository) CreateAPIUser(
	ctx context.Context,
	exec db.PGExecutor,
	user model.NewAPIUser,
) (string, error) {

	password := user.Password // allow caller to set the password

	if password == "" {
		randomPassword, err := generateRandomPassword(24)
		if err != nil {
			return "", fmt.Errorf("error generating random password: %v", err)
		}
		password = randomPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	permissionsJson, err := json.Marshal(user.Permissions)
	if err != nil {
		return "", err
	}

	insertUserStmt := `
INSERT INTO app_user (
	username,
	is_api_user,
	hashed_password,
	permissions
)
VALUES ($1, $2, $3, $4)
`

	_, err = exec.Exec(
		ctx,
		insertUserStmt,
		user.Username,
		true,
		string(hashedPassword),
		permissionsJson,
	)

	if err != nil {
		return "", err
	}

	return password, nil
}

func (r *UserRepository) UpdateUser(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	update model.UserUpdate,
) error {

	// get the user to check if it exists
	user, err := r.GetUserByID(ctx, exec, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("User does not exist")
	}
	if user.IsAPIUser {
		return fmt.Errorf("An API user cannot be updated using this method")
	}

	query := `
UPDATE
	app_user

SET
	first_name = $1,
	last_name = $2,
	email = $3,
	username = $4,
	permissions = $5,
	session_duration_minutes = $6

WHERE
	user_id = $7
	`

	permissionsJSON, err := json.Marshal(update.Permissions)

	if err != nil {
		return err
	}

	_, err = exec.Exec(
		ctx,
		query,

		update.FirstName,
		update.LastName,
		update.Email,
		update.Username,
		string(permissionsJSON),
		update.SessionDurationMinutes,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateLastActive(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	lastActive time.Time,
) error {
	_, err := exec.Exec(ctx, `
UPDATE
	app_user
SET
	last_active = $1
WHERE
	user_id = $2
	`, lastActive, userID)
	return err
}

func (r *UserRepository) ResetPassword(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
	pr model.PasswordReset,
) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pr.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
UPDATE
	app_user
SET
	hashed_password = $1
WHERE
	user_id = $2
	`

	_, err = exec.Exec(ctx, query, string(hashedPassword), userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByID(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
) (*model.User, error) {
	query := `
SELECT
    user_id,
    is_api_user,
    username,
    email,
    first_name,
    last_name,
    created,
    last_login,
    last_active,
	session_duration_minutes,
    permissions,
	teams
FROM
    user_view
WHERE
    user_id = $1
	`

	user := model.User{}
	var permissions json.RawMessage
	var teamsJSON []byte
	err := exec.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.IsAPIUser,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Created,
		&user.LastLogin,
		&user.LastActive,
		&user.SessionDurationMinutes,
		&permissions,
		&teamsJSON,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(permissions, &user.Permissions)
	if err != nil {
		user.Permissions = model.UserPermissions{}
	}

	if err := json.Unmarshal(teamsJSON, &user.Teams); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUsername(
	ctx context.Context,
	exec db.PGExecutor,
	username string,
) (*model.User, error) {

	query := `
SELECT
    user_id,
    is_api_user,
    username,
    email,
    first_name,
    last_name,
    created,
    last_login,
    last_active,
	session_duration_minutes,
    permissions,
	teams
FROM
    user_view
WHERE
    username = $1
	`

	user := model.User{}
	var permissions json.RawMessage
	var teamsJSON []byte
	err := exec.QueryRow(ctx, query, username).Scan(
		&user.UserID,
		&user.IsAPIUser,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Created,
		&user.LastLogin,
		&user.LastActive,
		&user.SessionDurationMinutes,
		&permissions,
		&teamsJSON,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(permissions, &user.Permissions)
	if err != nil {
		user.Permissions = model.UserPermissions{}
	}

	err = json.Unmarshal(teamsJSON, &user.Teams)
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		user.Permissions = model.UserPermissions{}
	}

	return &user, nil
}

func (r *UserRepository) GetUsers(
	ctx context.Context,
	exec db.PGExecutor,
	q model.GetUsersQuery,
) ([]model.User, error) {

	whereClause, args := generateUsersWhereClause(q)

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(model.User{})
	limitPlaceholder := fmt.Sprintf("$%d", len(args)+1)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args)+2)

	if orderByClause == "" {
		orderByClause = "ORDER BY username ASC"
	}

	query := `
SELECT
    user_id,
    is_api_user,
    username,
    email,
    first_name,
    last_name,
    created,
    last_login,
    last_active,
	session_duration_minutes,
    permissions,
	teams
FROM
    user_view

` + whereClause + `
` + orderByClause + `

LIMIT ` + limitPlaceholder + ` OFFSET ` + offsetPlaceholder + `
`

	rows, err := exec.Query(ctx, query, append(args, limit, offset)...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		var permissions json.RawMessage
		var teamsJSON []byte

		err := rows.Scan(
			&user.UserID,
			&user.IsAPIUser,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Created,
			&user.LastLogin,
			&user.LastActive,
			&user.SessionDurationMinutes,
			&permissions,
			&teamsJSON,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(teamsJSON, &user.Teams); err != nil {
			return nil, err
		}

		err = json.Unmarshal(permissions, &user.Permissions)
		if err != nil {
			user.Permissions = model.UserPermissions{}
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetUserCount(
	ctx context.Context,
	exec db.PGExecutor,
	q model.GetUsersQuery,
) (int, error) {

	whereClause, args := generateUsersWhereClause(q)

	query := `
SELECT
	COUNT(*)
FROM
	user_view
` + whereClause

	var count int
	err := exec.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func generateUsersWhereClause(q model.GetUsersQuery) (string, []any) {
	if q.Search == "" {
		return "", nil
	}

	args := []any{q.Search}

	where := `
WHERE
	COALESCE(username, '') ILIKE '%' || $1 || '%'
OR	COALESCE(email, '') ILIKE '%' || $1 || '%'
OR	COALESCE(first_name, '') ILIKE '%' || $1 || '%'
OR	COALESCE(last_name, '') ILIKE '%' || $1 || '%'
`

	return where, args
}

func (r *UserRepository) GetActiveUserCountSince(
	ctx context.Context,
	exec db.PGExecutor,
	since time.Time,
) (int, error) {
	var count int
	err := exec.QueryRow(ctx, `
SELECT COUNT(*)
FROM app_user
WHERE last_active IS NOT NULL
  AND last_active >= $1
`, since).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func generateRandomPassword(length int) (string, error) {
	if length < 8 {
		return "", fmt.Errorf("password length must be at least 8 characters")
	}

	const (
		lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
		uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers          = "0123456789"
		symbols          = "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~"
	)

	allChars := lowercaseLetters + uppercaseLetters + numbers + symbols

	// Use time-based seed for randomness
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	password := make([]byte, length)

	// Ensure at least one lowercase, one uppercase, one number, and one symbol
	password[0] = lowercaseLetters[randomGenerator.Intn(len(lowercaseLetters))]
	password[1] = uppercaseLetters[randomGenerator.Intn(len(uppercaseLetters))]
	password[2] = numbers[randomGenerator.Intn(len(numbers))]
	password[3] = symbols[randomGenerator.Intn(len(symbols))]

	// Fill the rest of the password randomly
	for i := 4; i < length; i++ {
		password[i] = allChars[randomGenerator.Intn(len(allChars))]
	}

	// Shuffle the password to randomize the order
	randomGenerator.Shuffle(length, func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password), nil
}
