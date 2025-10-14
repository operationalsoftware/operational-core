package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db             *pgxpool.Pool
	authRepository *repository.AuthRepository
}

func NewAuthService(
	db *pgxpool.Pool,
	authRepository *repository.AuthRepository,
) *AuthService {
	return &AuthService{
		db:             db,
		authRepository: authRepository,
	}
}

func (s *AuthService) VerifyPasswordLogin(
	ctx context.Context,
	input model.VerifyPasswordLoginInput,
) (model.VerifyPasswordLoginOutput, error) {

	out := model.VerifyPasswordLoginOutput{}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	authUser, err := s.authRepository.GetAuthUserByUsername(
		ctx,
		tx,
		input.Username,
	)
	if err != nil {
		return out, err
	}
	if authUser == nil {
		out.FailureReason = INVALID_USERNAME_PASSWORD_MSG
		return out, nil
	}

	// Check if login blocked
	now := time.Now()
	if authUser.LoginBlockedUntil != nil && now.Before(*authUser.LoginBlockedUntil) {
		out.FailureReason = LOGIN_BLOCKED_MSG
		return out, nil
	}

	// Check password
	var passwordErr error
	if authUser.HashedPassword == "" {
		passwordErr = bcrypt.ErrMismatchedHashAndPassword
	} else {
		passwordErr = bcrypt.CompareHashAndPassword([]byte(authUser.HashedPassword), []byte(input.Password))
	}

	if passwordErr == bcrypt.ErrMismatchedHashAndPassword {

		if authUser.FailedLoginAttempts < 5 {

			err := s.authRepository.IncrementFailedLoginAttempts(
				ctx, tx, authUser.UserID,
			)
			if err != nil {
				return out, err
			}

			out.FailureReason = INVALID_USERNAME_PASSWORD_MSG
		} else {

			loginBlockedUntil := now.Add(5 * time.Minute)

			err := s.authRepository.SetLoginBlockedUntil(
				ctx, tx, authUser.UserID, &loginBlockedUntil,
			)
			if err != nil {
				return out, err
			}

			out.FailureReason = LOGIN_BLOCKED_MSG
		}

		err = tx.Commit(ctx)
		if err != nil {
			return out, err
		}

		return out, nil
	} else if passwordErr != nil {
		return out, passwordErr
	}

	// Successful login
	err = s.authRepository.UpdateLastLogin(
		ctx, tx, authUser.UserID,
	)
	if err != nil {
		return out, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return out, err
	}

	// This means the user is verified
	out.AuthUser = *authUser

	return out, nil
}

const INVALID_USERNAME_PASSWORD_MSG = "Invalid username or password. Please try again"
const LOGIN_BLOCKED_MSG = "Login temporarily blocked, please wait and try again"

func (s *AuthService) AuthenticateByEmail(
	ctx context.Context,
	email string,
) (*model.AuthUser, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	authUser, err := s.authRepository.GetAuthUserByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}
	if authUser == nil {
		return nil, nil
	}

	// Update last_login for visibility
	if err := s.authRepository.UpdateLastLogin(ctx, tx, authUser.UserID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return authUser, nil
}

// AuthenticateUserByEmail links an existing user by email; no auto-creation.
func (s *AuthService) AuthenticateUserByEmail(
	ctx context.Context,
	email, externalID string,
	firstName, lastName string,
) (*model.AuthUser, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	authUser, err := s.authRepository.GetAuthUserByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}
	if authUser == nil {
		// Do not auto-create; enforce pre-created users only.
		return nil, errors.New("microsoft login allowed only for pre-created users")
	}

	if externalID == "" {
		return nil, errors.New("missing microsoft external id")
	}

	// If already linked, ensure it's the same Microsoft account.
	if authUser.ExternalID != nil && *authUser.ExternalID != "" && *authUser.ExternalID != externalID {
		return nil, errors.New("account already linked to a different microsoft id")
	}

	// Link external ID if not already set.
	if authUser.ExternalID == nil || *authUser.ExternalID == "" {
		if err := s.authRepository.SetExternalID(ctx, tx, authUser.UserID, externalID); err != nil {
			return nil, err
		}
		// Reflect change in the returned struct.
		authUser.ExternalID = &externalID
	}

	// Update last login
	if err := s.authRepository.UpdateLastLogin(ctx, tx, authUser.UserID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return authUser, nil
}
