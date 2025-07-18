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
	stockCode string,
) (*model.StockItem, error) {

	stockItem, err := s.stockItemRepository.GetStockItem(ctx, s.db, stockCode)
	if err != nil {
		return &model.StockItem{}, err
	}

	return stockItem, nil
}

func (s *StockItemService) GetStockItemChanges(
	ctx context.Context,
	stockCode string,
) ([]model.StockItemChange, error) {

	stockItem, err := s.stockItemRepository.GetStockItemChanges(ctx, s.db, stockCode)
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

	stockItem, _ := s.stockItemRepository.GetStockItem(ctx, tx, input.StockCode)
	if stockItem != nil {
		validationErrors.Add("StockCode", "already exists")
		return validationErrors, err
	}

	err = s.stockItemRepository.CreateStockItem(ctx, tx, input)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	err = s.stockItemRepository.AddStockItemChange(ctx, tx, model.PostStockItemChange{
		StockCode:        input.StockCode,
		StockCodeHistory: nil,
		Description:      &input.Description,
		ChangeBy:         userID,
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
		StockCode:   input.StockCode,
		Description: nil,
		ChangeBy:    userID,
	}

	// Set fields only if they changed
	if stockItem.Description != input.Description {
		change.Description = &input.Description
	}

	if stockItem.StockCode != input.StockCode {
		change.StockCodeHistory = &stockCode
	}

	// Only insert if at least one field changed
	if change.Description != nil || change.StockCodeHistory != nil {
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

func (s *StockItemService) GetSKUConfiguration(
	ctx context.Context,
) (model.SKUConfigData, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.SKUConfigData{}, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	skuItems, err := s.stockItemRepository.GetSKUConfiguration(ctx, tx)
	if err != nil {
		return model.SKUConfigData{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.SKUConfigData{}, err
	}

	return skuItems, nil
}

func (s *StockItemService) CreateSKUConfigItem(
	ctx context.Context,
	input *model.SKUConfigItem,
) (
	validate.ValidationErrors,
	error,
) {
	validationErrors, err := s.validateSKUConfigItem(input)
	if err != nil || len(validationErrors) > 0 {
		return validationErrors, err
	}

	err = s.stockItemRepository.CreateSKUConfigItem(ctx, s.db, input)
	if err != nil {
		return validate.ValidationErrors{}, err
	}

	return validate.ValidationErrors{}, nil
}

func (s *StockItemService) DeleteSKUConfigItem(
	ctx context.Context,
	skuField string,
	skuCode string,
) error {

	err := s.stockItemRepository.DeleteSKUConfigItem(ctx, s.db, skuField, skuCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *StockItemService) validateNewStockItem(
	stockItem *model.PostStockItem,
) (validate.ValidationErrors, error) {

	var ve validate.ValidationErrors = make(map[string][]string)

	if stockItem.ProductType == "" {
		ve.Add("ProductType", "is required")
	}

	if stockItem.YarnType == "" {
		ve.Add("YarnType", "is required")
	}

	if stockItem.StyleNumber == "" {
		ve.Add("StyleNumber", "is required")
	}

	if stockItem.Colour == "" {
		ve.Add("Colour", "is required")
	}

	if stockItem.ToeClosing == "" {
		ve.Add("ToeClosing", "is required")
	}

	if stockItem.Size == "" {
		ve.Add("Size", "is required")
	}

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

	if stockItem.ProductType == "" {
		ve.Add("ProductType", "should not be empty")
	}

	if stockItem.YarnType == "" {
		ve.Add("YarnType", "should not be empty")
	}

	if stockItem.StyleNumber == "" {
		ve.Add("StyleNumber", "should not be empty")
	}

	if stockItem.Colour == "" {
		ve.Add("Colour", "should not be empty")
	}

	if stockItem.ToeClosing == "" {
		ve.Add("ToeClosing", "should not be empty")
	}

	if stockItem.Size == "" {
		ve.Add("Size", "should not be empty")
	}

	if stockItem.Description == "" {
		ve.Add("Description", "should not be empty")
	}

	return ve, nil
}

func (s *StockItemService) validateSKUConfigItem(
	skuItem *model.SKUConfigItem,
) (validate.ValidationErrors, error) {

	var ve validate.ValidationErrors = make(map[string][]string)

	if skuItem.SKUField == "" {
		ve.Add("SKUField", "should not be empty")
	}

	if skuItem.Label == "" {
		ve.Add("Label", "should not be empty")
	}

	if skuItem.Code == "" {
		ve.Add("Code", "should not be empty")
	}

	return ve, nil
}
