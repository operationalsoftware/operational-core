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

	password, err := generateRandomPassword(24)
	if err != nil {
		log.Panic(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
	permissions = $5

WHERE
	user_id = $6
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
		userID,
	)

	if err != nil {
		return err
	}

	return nil
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
    permissions
FROM
    app_user
WHERE
    user_id = $1
	`

	userDB := model.UserDB{}
	err := exec.QueryRow(ctx, query, userID).Scan(
		&userDB.UserID,
		&userDB.IsAPIUser,
		&userDB.Username,
		&userDB.Email,
		&userDB.FirstName,
		&userDB.LastName,
		&userDB.Created,
		&userDB.LastLogin,
		&userDB.Permissions,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user := userDB.ToDomain()
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
    permissions
FROM
    app_user
WHERE
    username = $1
	`

	userDB := model.UserDB{}
	err := exec.QueryRow(ctx, query, username).Scan(
		&userDB.UserID,
		&userDB.IsAPIUser,
		&userDB.Username,
		&userDB.Email,
		&userDB.FirstName,
		&userDB.LastName,
		&userDB.Created,
		&userDB.LastLogin,
		&userDB.Permissions,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user := userDB.ToDomain()
	return &user, nil
}

func (r *UserRepository) GetUsers(
	ctx context.Context,
	exec db.PGExecutor,
	q model.GetUsersQuery,
) ([]model.User, error) {

	offset := (q.Page - 1) * q.PageSize
	limit := q.PageSize
	orderByClause := q.Sort.ToOrderByClause(map[string]string{})

	if orderByClause == "" {
		orderByClause = "ORDER BY username ASC"
	}

	query := fmt.Sprintf(`
SELECT
    user_id,
    is_api_user,
    username,
    email,
    first_name,
    last_name,
    created,
    last_login,
    permissions
FROM
    app_user

%s

LIMIT $1 OFFSET $2
	`,
		orderByClause,
	)

	rows, err := exec.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var userDB model.UserDB
		err := rows.Scan(
			&userDB.UserID,
			&userDB.IsAPIUser,
			&userDB.Username,
			&userDB.Email,
			&userDB.FirstName,
			&userDB.LastName,
			&userDB.Created,
			&userDB.LastLogin,
			&userDB.Permissions,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, userDB.ToDomain())
	}

	return users, nil
}

func (r *UserRepository) GetUserCount(
	ctx context.Context,
	exec db.PGExecutor,
) (int, error) {

	query := `
SELECT
	COUNT(*)
FROM
	app_user
	`

	var count int
	err := exec.QueryRow(ctx, query).Scan(&count)
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
