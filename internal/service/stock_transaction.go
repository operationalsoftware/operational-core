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

func (s *StockTrxService) PostStockTransaction(
	ctx context.Context,
	input *model.PostStockTransactionsInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	err = s.stockTrxRepository.PostStockTransactions(ctx, tx, input, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTrxService) GetStockTransactions(
	ctx context.Context,
	input *model.GetTransactionsInput,
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
	input *model.GetStockLevelsInput,
) ([]model.StockLevel, error) {

	levels, err := s.stockTrxRepository.GetStockLevels(ctx, s.db, input)

	if err != nil {
		return nil, err
	}

	return levels, nil
}
