package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/validate"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockItemService struct {
	db                  *pgxpool.Pool
	stockItemRepository *repository.StockItemRepository
}

func NewStockItemService(
	db *pgxpool.Pool,
	stockItemRepository *repository.StockItemRepository,
) *StockItemService {
	return &StockItemService{
		db:                  db,
		stockItemRepository: stockItemRepository,
	}
}

func (s *StockItemService) GetStockItems(
	ctx context.Context,
	input *model.GetStockItemsQuery,
	userID int,
) ([]model.StockItem, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.StockItem{}, 0, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	stockItems, err := s.stockItemRepository.GetStockItems(ctx, tx, input)
	if err != nil {
		return []model.StockItem{}, 0, err
	}

	count, err := s.stockItemRepository.GetStockItemsCount(ctx, tx)
	if err != nil {
		return []model.StockItem{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.StockItem{}, 0, err
	}

	return stockItems, count, nil
}

func (s *StockItemService) GetStockItem(
	ctx context.Context,
	stockItemID int,
) (*model.StockItem, error) {

	stockItem, err := s.stockItemRepository.GetStockItem(ctx, s.db, stockItemID)
	if err != nil {
		return &model.StockItem{}, err
	}

	return stockItem, nil
}

func (s *StockItemService) GetStockItemChanges(
	ctx context.Context,
	stockItemID int,
) ([]model.StockItemChange, error) {

	stockItem, err := s.stockItemRepository.GetStockItemChanges(ctx, s.db, stockItemID)
	if err != nil {
		return nil, err
	}

	return stockItem, nil
}

func (s *StockItemService) CreateStockItem(
	ctx context.Context,
	input *model.PostStockItem,
	userID int,
) (
	validate.ValidationErrors,
	error,
) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}
	defer tx.Rollback(ctx)

	validationErrors, err := s.validateNewStockItem(input)
	if err != nil || len(validationErrors) > 0 {
		return validationErrors, err
	}

	stockItem, _ := s.stockItemRepository.GetStockItemByStockCode(ctx, tx, input.StockCode)
	if stockItem != nil {
		validationErrors.Add("StockCode", "already exists")
		return validationErrors, err
	}

	err = s.stockItemRepository.CreateStockItem(ctx, tx, input)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	err = s.stockItemRepository.AddStockItemChange(ctx, tx, model.PostStockItemChange{
		StockCode:   &input.StockCode,
		Description: &input.Description,
		ChangeBy:    userID,
	})
	if err != nil {
		fmt.Println(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	return validate.ValidationErrors{}, nil
}

func (s *StockItemService) UpdateStockItem(
	ctx context.Context,
	stockCode string,
	input *model.PostStockItem,
	userID int,
) (
	validate.ValidationErrors,
	error,
) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}
	defer tx.Rollback(ctx)

	validationErrors, err := s.validateUpdateStockItem(input)
	if err != nil || len(validationErrors) > 0 {
		return validationErrors, err
	}

	stockItem, _ := s.stockItemRepository.GetStockItem(ctx, tx, stockCode)

	if input.StockCode != stockCode {
		existing, _ := s.stockItemRepository.GetStockItem(ctx, tx, input.StockCode)
		if existing != nil {
			validationErrors.Add("StockCode", "stock code already exists")
			return validationErrors, nil
		}
	}

	err = s.stockItemRepository.UpdateStockItem(ctx, tx, stockCode, input)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	change := model.PostStockItemChange{
		StockCode:   nil,
		Description: nil,
		ChangeBy:    userID,
	}

	if stockItem.Description != input.Description {
		change.Description = &input.Description
	}

	if stockItem.StockCode != input.StockCode {
		change.StockCode = &stockCode
	}

	// Only insert if at least one field changed
	if change.Description != nil || change.StockCode != nil {
		err = s.stockItemRepository.AddStockItemChange(ctx, tx, change)
		if err != nil {
			fmt.Println(err)
			return validate.ValidationErrors{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	return validate.ValidationErrors{}, nil
}

func (s *StockItemService) validateNewStockItem(
	stockItem *model.PostStockItem,
) (validate.ValidationErrors, error) {

	var ve validate.ValidationErrors = make(map[string][]string)

	if stockItem.StockCode == "" {
		ve.Add("StockCode", "is required")
	}

	if stockItem.Description == "" {
		ve.Add("Description", "is required")
	}

	return ve, nil
}

func (s *StockItemService) validateUpdateStockItem(
	stockItem *model.PostStockItem,
) (validate.ValidationErrors, error) {

	var ve validate.ValidationErrors = make(map[string][]string)

	if stockItem.StockCode == "" {
		ve.Add("StockCode", "should not be empty")
	}

	if stockItem.Description == "" {
		ve.Add("Description", "should not be empty")
	}

	return ve, nil
}
