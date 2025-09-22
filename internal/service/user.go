package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/validate"
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db             *pgxpool.Pool
	userRepository *repository.UserRepository
}

func NewUserService(db *pgxpool.Pool, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		db:             db,
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	user model.NewUser,
) (
	validate.ValidationErrors,
	error,
) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	validationErrors, err := s.validateNewUser(ctx, tx, user)
	if err != nil || len(validationErrors) > 0 {
		return validationErrors, err
	}

	err = s.userRepository.CreateUser(ctx, tx, user)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	return validate.ValidationErrors{}, nil
}

func (s *UserService) CreateAPIUser(
	ctx context.Context,
	user model.NewAPIUser,
) (
	validate.ValidationErrors,
	string,
	error,
) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return validate.ValidationErrors{}, "", err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	validationErrors, err := s.validateNewAPIUser(ctx, tx, user)
	if err != nil || len(validationErrors) > 0 {
		return validationErrors, "", err
	}

	password, err := s.userRepository.CreateAPIUser(ctx, tx, user)
	if err != nil {
		return validationErrors, "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return validate.ValidationErrors{}, "", err
	}

	return validate.ValidationErrors{}, password, nil
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	userID int,
	update model.UserUpdate,
) (
	validate.ValidationErrors,
	error,
) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	validationErrors, err := s.validateUserUpdate(ctx, tx, userID, update)
	if err != nil || len(validationErrors) > 0 {
		return validationErrors, err
	}

	err = s.userRepository.UpdateUser(ctx, tx, userID, update)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	return validate.ValidationErrors{}, nil
}

func (s *UserService) ResetPassword(
	ctx context.Context,
	userID int,
	pr model.PasswordReset,
) (*model.User, validate.ValidationErrors, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return &model.User{}, validate.ValidationErrors{}, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	user, err := s.userRepository.GetUserByID(ctx, tx, userID)
	if err != nil {
		return &model.User{}, validate.ValidationErrors{}, err
	}
	if user == nil {
		return user, validate.ValidationErrors{}, err
	}

	validationErrors := s.validatePasswordReset(pr)
	if len(validationErrors) > 0 {
		return user, validationErrors, err
	}

	err = s.userRepository.ResetPassword(ctx, tx, userID, pr)
	if err != nil {
		return user, validate.ValidationErrors{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return user, validate.ValidationErrors{}, err
	}

	return user, validate.ValidationErrors{}, nil
}

func (s *UserService) GetUserByID(
	ctx context.Context,
	userID int,
) (
	*model.User,
	error,
) {

	user, err := s.userRepository.GetUserByID(ctx, s.db, userID)
	if err != nil {
		return &model.User{}, err
	}

	return user, nil
}

func (s *UserService) GetUserByUsername(
	ctx context.Context,
	username string,
) (
	*model.User,
	error,
) {
	user, err := s.userRepository.GetUserByUsername(ctx, s.db, username)
	if err != nil {
		return &model.User{}, err
	}

	return user, nil
}

func (s *UserService) GetUsers(
	ctx context.Context,
	q model.GetUsersQuery,
) ([]model.User, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.User{}, 0, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	users, err := s.userRepository.GetUsers(ctx, tx, q)
	if err != nil {
		return []model.User{}, 0, err
	}

	count, err := s.userRepository.GetUserCount(ctx, tx)
	if err != nil {
		return []model.User{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.User{}, 0, err
	}

	return users, count, nil
}

func (s *UserService) validateNewUser(
	ctx context.Context,
	tx pgx.Tx,
	user model.NewUser,
) (validate.ValidationErrors, error) {

	var ve validate.ValidationErrors = make(map[string][]string)

	err := s.validateUsername(ctx, tx, &ve, "Username", user.Username)
	if err != nil {
		return ve, err
	}
	s.validateFirstName(&ve, "FirstName", user.FirstName)
	s.validateLastName(&ve, "LastName", user.LastName)

	if user.Email != nil {
		validate.Email(&ve, "Email", *user.Email)
	}

	validate.Password(&ve, "Password", user.Password)

	if user.Password != user.ConfirmPassword {
		ve.Add("ConfirmPassword", "does not match")
	}

	// user permissions don't need to be validated, the struct will be populated with
	// matching permissions from the form data meaning any missing data will be
	// zero-valued (boolean, false) and any extra data will be ignored

	return ve, nil
}

func (s *UserService) validateNewAPIUser(
	ctx context.Context,
	tx pgx.Tx,
	user model.NewAPIUser,
) (validate.ValidationErrors, error) {
	var ve validate.ValidationErrors = make(map[string][]string)

	err := s.validateUsername(ctx, tx, &ve, "Username", user.Username)
	if err != nil {
		return ve, err
	}

	return ve, nil
}

func (s *UserService) validateUserUpdate(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	update model.UserUpdate,
) (validate.ValidationErrors, error) {

	var ve validate.ValidationErrors = make(map[string][]string)

	// get current username
	user, err := s.userRepository.GetUserByID(ctx, tx, userID)
	if err != nil {
		return validate.ValidationErrors{}, err
	}
	if user == nil {
		return validate.ValidationErrors{}, fmt.Errorf("User does not exist")
	}

	if user.Username != update.Username {
		err = s.validateUsername(ctx, tx, &ve, "Username", update.Username)
		if err != nil {
			return validate.ValidationErrors{}, err
		}
	}

	s.validateFirstName(&ve, "FirstName", update.FirstName)
	s.validateLastName(&ve, "LastName", update.LastName)

	if update.SessionDurationMinutes != nil {
		s.validateSessionDuration(&ve, "SessionDurationMinutes", *update.SessionDurationMinutes)
	}

	if update.Email != nil {
		validate.Email(&ve, "Email", *update.Email)
	}

	if update.SessionDurationMinutes != nil {
		validate.IntGT(&ve, "SessionDurationMinutes", *update.SessionDurationMinutes, 1)
	}

	// user permissions don't need to be validated. See description in validateNewUser

	return ve, nil
}

func (s *UserService) validatePasswordReset(
	pr model.PasswordReset,
) validate.ValidationErrors {

	var validationErrors validate.ValidationErrors = make(map[string][]string)

	// validate Password
	validate.Password(&validationErrors, "Password", pr.Password)

	// validate confirm password
	if pr.Password != pr.ConfirmPassword {
		validationErrors.Add("ConfirmPassword", "does not match")
	}

	return validationErrors
}

func (s *UserService) validateUsername(
	ctx context.Context,
	tx pgx.Tx,
	ve *validate.ValidationErrors,
	usernameKey string,
	username string,
) error {

	validate.MinLength(ve, usernameKey, username, 3)
	validate.MaxLength(ve, usernameKey, username, 20)
	validate.Lowercase(ve, usernameKey, username)

	pattern := "^[a-z0-9_]+$"

	// Compile the regular expression
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Error compiling regex:", err) // should never happen
	}
	if !re.MatchString(username) {
		ve.Add(usernameKey, "must contain only letters, numbers, and underscores")
	}

	// check if username is already taken
	existingUser, err := s.userRepository.GetUserByUsername(ctx, tx, username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		ve.Add(usernameKey, fmt.Sprintf("'%s' is already taken", username))
	}

	return nil
}

func (s *UserService) validateFirstName(
	ve *validate.ValidationErrors,
	firstNameKey string,
	firstName string,
) {
	validate.MinLength(ve, firstNameKey, firstName, 1)
	validate.MaxLength(ve, firstNameKey, firstName, 20)
}

func (s *UserService) validateLastName(
	ve *validate.ValidationErrors,
	lastNameKey string,
	lastName string,
) {
	validate.MinLength(ve, lastNameKey, lastName, 1)
	validate.MaxLength(ve, lastNameKey, lastName, 20)
}

func (s *UserService) validateSessionDuration(
	ve *validate.ValidationErrors,
	sessDurationKey string,
	duration int,
) {
	validate.IntGTE(ve, sessDurationKey, duration, 1)
	validate.IntLTE(ve, sessDurationKey, duration, 525600)
}
