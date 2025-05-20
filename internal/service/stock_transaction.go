package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockTrxService struct {
	db                 *pgxpool.Pool
	stockTrxRepository *repository.StockTrxRepository
}

func NewStockTrxService(
	db *pgxpool.Pool,
	stockTrxRepository *repository.StockTrxRepository,
) *StockTrxService {
	return &StockTrxService{
		db:                 db,
		stockTrxRepository: stockTrxRepository,
	}
}

func (s *StockTrxService) CreateStockTransaction(
	ctx context.Context,
	input *model.StockTrxPostInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	err = s.stockTrxRepository.CreateStockTransaction(ctx, tx, input, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTrxService) GetStockTransaction(
	ctx context.Context,
	input *model.StockTrxPostInput,
	userID int,
) ([]model.StockTransactionEntry, error) {

	transactions, err := s.stockTrxRepository.GetStockTransactions(ctx, s.db, input)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *StockTrxService) GetStockLevels(
	ctx context.Context,
	input model.VerifyPasswordLoginInput,
) (model.VerifyPasswordLoginOutput, error) {

	out := model.VerifyPasswordLoginOutput{}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return out, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	// authUser, err := s.stockTrxRepository.GetAuthUserByUsername(
	// 	ctx,
	// 	tx,
	// 	input.Username,
	// )
	// if err != nil {
	// 	return out, err
	// }
	// if authUser == nil {
	// 	out.FailureReason = INVALID_USERNAME_PASSWORD_MSG
	// 	return out, nil
	// }

	// // Successful login
	// err = s.authRepository.UpdateLastLogin(
	// 	ctx, tx, authUser.UserID,
	// )
	// if err != nil {
	// 	return out, err
	// }

	err = tx.Commit(ctx)
	if err != nil {
		return out, err
	}

	return out, nil
}
