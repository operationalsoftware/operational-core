package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type StockTransactionService struct {
	db                         *pgxpool.Pool
	stockTransactionRepository *repository.StockTransactionRepository
}

func NewStockTransactionService(
	db *pgxpool.Pool,
	stockTransactionRepository *repository.StockTransactionRepository,
) *StockTransactionService {
	return &StockTransactionService{
		db:                         db,
		stockTransactionRepository: stockTransactionRepository,
	}
}

func (s *StockTransactionService) PostManualStockMovement(
	ctx context.Context,
	input *model.PostManualStockMovementInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	err = s.stockTransactionRepository.PostStockTransactions(ctx, tx, &model.PostStockTransactionsInput{{
		TransactionType: "Stock Movement",
		StockCode:       input.StockCode,
		Qty:             input.Qty,
		FromLocation:    input.FromLocation,
		FromBin:         input.FromBin,
		FromLotNumber:   input.LotNumber,
		ToLocation:      input.ToLocation,
		ToBin:           input.ToBin,
		ToLotNumber:     input.LotNumber,
		TransactionNote: input.TransactionNote,
		Timestamp:       nil,
	}}, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTransactionService) PostManualProduction(
	ctx context.Context,
	input *model.PostManualGenericStockTransactionInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	err = s.stockTransactionRepository.PostStockTransactions(ctx, tx, &model.PostStockTransactionsInput{{
		TransactionType: "Production",
		StockCode:       input.StockCode,
		Qty:             input.Qty,
		FromLocation:    input.Location,
		FromBin:         input.Bin,
		FromLotNumber:   input.LotNumber,
		ToLocation:      input.Location,
		ToBin:           input.Bin,
		ToLotNumber:     input.LotNumber,
		TransactionNote: input.TransactionNote,
		Timestamp:       nil,
	}}, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTransactionService) PostManualProductionReversal(
	ctx context.Context,
	input *model.PostManualGenericStockTransactionInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	err = s.stockTransactionRepository.PostStockTransactions(ctx, tx, &model.PostStockTransactionsInput{{
		TransactionType: "Production Reversal",
		StockCode:       input.StockCode,
		Qty:             input.Qty,
		FromLocation:    input.Location,
		FromBin:         input.Bin,
		FromLotNumber:   input.LotNumber,
		ToLocation:      input.Location,
		ToBin:           input.Bin,
		ToLotNumber:     input.LotNumber,
		TransactionNote: input.TransactionNote,
		Timestamp:       nil,
	}}, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTransactionService) PostManualConsumption(
	ctx context.Context,
	input *model.PostManualGenericStockTransactionInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	err = s.stockTransactionRepository.PostStockTransactions(ctx, tx, &model.PostStockTransactionsInput{{
		TransactionType: "Consumption",
		StockCode:       input.StockCode,
		Qty:             input.Qty,
		FromLocation:    input.Location,
		FromBin:         input.Bin,
		FromLotNumber:   input.LotNumber,
		ToLocation:      input.Location,
		ToBin:           input.Bin,
		ToLotNumber:     input.LotNumber,
		TransactionNote: input.TransactionNote,
		Timestamp:       nil,
	}}, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTransactionService) PostManualConsumptionReversal(
	ctx context.Context,
	input *model.PostManualGenericStockTransactionInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	err = s.stockTransactionRepository.PostStockTransactions(ctx, tx, &model.PostStockTransactionsInput{{
		TransactionType: "Consumption Reversal",
		StockCode:       input.StockCode,
		Qty:             input.Qty,
		FromLocation:    input.Location,
		FromBin:         input.Bin,
		FromLotNumber:   input.LotNumber,
		ToLocation:      input.Location,
		ToBin:           input.Bin,
		ToLotNumber:     input.LotNumber,
		TransactionNote: input.TransactionNote,
		Timestamp:       nil,
	}}, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTransactionService) PostManualStockAdjustment(
	ctx context.Context,
	input *model.PostManualGenericStockTransactionInput,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	var transactionType model.StockTransactionType

	if input.Qty.GreaterThan(decimal.Zero) {
		transactionType = "Stock Adjust Up"
	} else {
		transactionType = "Stock Adjust Down"
	}

	err = s.stockTransactionRepository.PostStockTransactions(ctx, tx, &model.PostStockTransactionsInput{{
		TransactionType: transactionType,
		StockCode:       input.StockCode,
		Qty:             input.Qty,
		FromLocation:    input.Location,
		FromBin:         input.Bin,
		FromLotNumber:   input.LotNumber,
		ToLocation:      input.Location,
		ToBin:           input.Bin,
		ToLotNumber:     input.LotNumber,
		TransactionNote: input.TransactionNote,
		Timestamp:       nil,
	}}, userID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func (s *StockTransactionService) GetStockTransactions(
	ctx context.Context,
	input *model.GetTransactionsInput,
	userID int,
) ([]model.StockTransactionEntry, error) {

	transactions, err := s.stockTransactionRepository.GetStockTransactions(ctx, s.db, input)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *StockTransactionService) GetStockLevels(
	ctx context.Context,
	input *model.GetStockLevelsInput,
) ([]model.StockLevel, error) {

	levels, err := s.stockTransactionRepository.GetStockLevels(ctx, s.db, input)

	if err != nil {
		return nil, err
	}

	return levels, nil
}
