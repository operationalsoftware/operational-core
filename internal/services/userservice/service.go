package userservice

import (
	"app/internal/models"
	"app/pkg/validate"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ValidateNewUser(ctx context.Context, user models.NewUser) (bool, validate.ValidationErrors)
	ValidateNewAPIUser(ctx context.Context, user models.NewAPIUser) (bool, validate.ValidationErrors)
	ValidateUserUpdate(ctx context.Context, update models.UserUpdate) (bool, validate.ValidationErrors)
	ValidatePasswordReset(pr models.PasswordReset) (bool, validate.ValidationErrors)

	CreateUser(ctx context.Context, user models.NewUser) error
	CreateAPIUser(ctx context.Context, user models.NewAPIUser) (string, error)
	UpdateUser(ctx context.Context, id int, update models.UserUpdate) error
	ResetPassword(ctx context.Context, id int, pr models.PasswordReset) error

	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUsers(ctx context.Context, q models.GetUsersQuery) ([]models.User, error)
	GetUserCount(ctx context.Context) (int, error)
}

type userService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) UserService {
	return &userService{db: db}
}

func (s *userService) ValidateNewUser(
	ctx context.Context, user models.NewUser,
) (bool, validate.ValidationErrors) {

	var validationErrors validate.ValidationErrors = make(map[string][]string)

	// validate Username
	validate.MinLength(&validationErrors, "Username", user.Username, 3)
	validate.MaxLength(&validationErrors, "Username", user.Username, 20)
	validate.Lowercase(&validationErrors, "Username", user.Username)
	validateUsername(&validationErrors, user.Username)
	// check if username is already taken
	if user.Username != "" {
		_, err := s.GetUserByUsername(ctx, user.Username)
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
	if user.Email != nil {
		validate.Email(&validationErrors, "Email", *user.Email)
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

func (s *userService) CreateUser(ctx context.Context, user models.NewUser) error {

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
	_, err = s.db.Exec(
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

func (s *userService) ValidateNewAPIUser(
	ctx context.Context, user models.NewAPIUser,
) (bool, validate.ValidationErrors) {
	var validationErrors validate.ValidationErrors = make(map[string][]string)

	// validate Username
	validate.MinLength(&validationErrors, "Username", user.Username, 3)
	validate.MaxLength(&validationErrors, "Username", user.Username, 20)
	validate.Lowercase(&validationErrors, "Username", user.Username)
	validateUsername(&validationErrors, user.Username)
	// check if username is already taken
	if user.Username != "" {
		_, err := s.GetUserByUsername(ctx, user.Username)
		if err == nil {
			validationErrors.Add("Username", fmt.Sprintf("'%s' is already taken", user.Username))
		}
	}

	return len(validationErrors) == 0, validationErrors
}

func (s *userService) CreateAPIUser(ctx context.Context, user models.NewAPIUser) (string, error) {

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

	_, err = s.db.Exec(
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

func (s *userService) ValidateUserUpdate(
	ctx context.Context,
	update models.UserUpdate,
) (bool, validate.ValidationErrors) {

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
	if update.Email != nil {
		validate.Email(&validationErrors, "Email", *update.Email)
	}

	// user permissions don't need to be validated. See description in ValidateNewUser

	return len(validationErrors) == 0, validationErrors
}

func (s *userService) UpdateUser(
	ctx context.Context, id int, update models.UserUpdate,
) error {

	// get the user to check if it exists
	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// can't update an API user using this method
	if user.IsAPIUser {
		return fmt.Errorf("API users cannot be updated using this method")
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

	_, err = s.db.Exec(
		ctx,
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

func (s *userService) ValidatePasswordReset(pr models.PasswordReset) (bool, validate.ValidationErrors) {

	var validationErrors validate.ValidationErrors = make(map[string][]string)

	// validate Password
	validate.Password(&validationErrors, "Password", pr.Password)

	// validate confirm password
	if pr.Password != pr.ConfirmPassword {
		validationErrors.Add("ConfirmPassword", "does not match")
	}

	return len(validationErrors) == 0, validationErrors
}

func (s *userService) ResetPassword(
	ctx context.Context, id int, pr models.PasswordReset,
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

	_, err = s.db.Exec(ctx, query, string(hashedPassword), id)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
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

	userDB := models.UserDB{}
	err := s.db.QueryRow(ctx, query, id).Scan(
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
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user := userDB.ToDomain()
	return &user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {

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

	userDB := models.UserDB{}
	err := s.db.QueryRow(ctx, query, username).Scan(
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
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user := userDB.ToDomain()
	return &user, nil
}

func (s *userService) GetUsers(ctx context.Context, q models.GetUsersQuery) ([]models.User, error) {

	offset := (q.Page - 1) * q.PageSize
	limit := q.PageSize
	orderByClause := q.Sort.ToOrderByClause(map[string]string{})

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

	rows, err := s.db.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var userDB models.UserDB
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

func (s *userService) GetUserCount(ctx context.Context) (int, error) {
	query := `
SELECT
	COUNT(*)
FROM
	app_user
	`

	var count int
	err := s.db.QueryRow(ctx, query).Scan(&count)
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
