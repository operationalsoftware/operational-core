package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

var StockAccounts = []string{"STOCK", "PRODUCTION", "CONSUMED"}

type StockTransactionEntry struct {
	StockTransactionEntryID int
	TransactionType         string
	Account                 string
	StockCode               string
	Location                string
	Bin                     string
	Quantity                decimal.Decimal
	LotNumber               string
	RunningTotal            decimal.Decimal
	TransactionBy           string
	TransactionByUsername   string
	Timestamp               time.Time
	StockTransactionID      int
}

type GetTransactionsInput struct {
	Account      string
	StockCode    string
	Location     string
	Bin          string
	LotNumber    string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
}

// type GetTransactionsInput struct {
// 	Account      *string
// 	StockCode    *string
// 	Location     *string
// 	Bin          *string
// 	LotNumber    *string
// 	LTETimestamp *time.Time
// 	Page         int
// 	PageSize     int
// }

type GetTransactionsDBInput struct {
	Account      pgtype.Text
	StockCode    pgtype.Text
	Location     pgtype.Text
	Bin          pgtype.Text
	LotNumber    pgtype.Text
	LTETimestamp pgtype.Timestamptz
	Page         int
	PageSize     int
}

// func (in GetTransactionsInput) ToDBModel() GetTransactionsDBInput {
// 	return GetTransactionsDBInput{
// 		Account:      pgconv.StringPtrToPGText(in.Account),
// 		StockCode:    pgconv.StringPtrToPGText(in.StockCode),
// 		Location:     pgconv.StringPtrToPGText(in.Location),
// 		Bin:          pgconv.StringPtrToPGText(in.Bin),
// 		LotNumber:    pgconv.StringPtrToPGText(in.LotNumber),
// 		LTETimestamp: pgconv.TimePtrToPGTimestamptz(in.LTETimestamp),
// 		Page:         in.Page,
// 		PageSize:     in.PageSize,
// 	}
// }

// func (in GetTransactionsDBInput) ToDomain() GetTransactionsInput {
// 	return GetTransactionsInput{
// 		Account:      pgconv.PGTextToStringPtr(in.Account),
// 		StockCode:    pgconv.PGTextToStringPtr(in.StockCode),
// 		Location:     pgconv.PGTextToStringPtr(in.Location),
// 		Bin:          pgconv.PGTextToStringPtr(in.Bin),
// 		LotNumber:    pgconv.PGTextToStringPtr(in.LotNumber),
// 		LTETimestamp: pgconv.PGTimestamptzToTimePtr(in.LTETimestamp),
// 		Page:         in.Page,
// 		PageSize:     in.PageSize,
// 	}
// }

type NewStockTransaction struct {
	Timestamp       *time.Time
	StockCode       string
	Qty             decimal.Decimal
	FromAccount     string
	FromLocation    string
	FromBin         string
	FromLotNumber   string
	ToAccount       string
	ToLocation      string
	ToBin           string
	ToLotNumber     string
	TransactionType string
	TransactionNote string
}

type NewStockTransactionDB struct {
	Timestamp       time.Time
	StockCode       string
	Qty             decimal.Decimal
	FromAccount     string
	FromLocation    string
	FromBin         string
	FromLotNumber   pgtype.Text
	ToAccount       string
	ToLocation      string
	ToBin           string
	ToLotNumber     pgtype.Text
	TransactionType string
	TransactionNote string
}

type PostStockTransactionsInput []NewStockTransaction
type PostStockTransactionsInputDB []NewStockTransactionDB

type PostManualGenericStockTransactionInput struct {
	StockCode       string
	Qty             decimal.Decimal
	Location        string
	Bin             string
	LotNumber       string
	TransactionNote string
}

type PostManualStockMovementInput struct {
	StockCode       string
	Qty             decimal.Decimal
	FromLocation    string
	FromBin         string
	ToLocation      string
	ToBin           string
	LotNumber       string
	TransactionNote string
}

type GetStockLevelsInput struct {
	Account      string
	StockCode    string
	Location     string
	Bin          string
	LotNumber    string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
}

type GetStockLevelsDBInput struct {
	Account      pgtype.Text
	StockCode    pgtype.Text
	Location     pgtype.Text
	Bin          pgtype.Text
	LotNumber    pgtype.Text
	LTETimestamp pgtype.Timestamptz
	Page         int
	PageSize     int
}

// func (in *GetStockLevelsInput) ToDB() GetStockLevelsDBInput {

// 	return GetStockLevelsDBInput{
// 		Account:      pgconv.StringPtrToPGText(in.Account),
// 		StockCode:    pgconv.StringPtrToPGText(in.StockCode),
// 		Location:     pgconv.StringPtrToPGText(in.Location),
// 		Bin:          pgconv.StringPtrToPGText(in.Bin),
// 		LotNumber:    pgconv.StringPtrToPGText(in.LotNumber),
// 		LTETimestamp: pgconv.TimePtrToPGTimestamptz(in.LTETimestamp),
// 		Page:         in.Page,
// 		PageSize:     in.PageSize,
// 	}
// }

// func (in GetStockLevelsDBInput) ToDomain() GetStockLevelsInput {
// 	return GetStockLevelsInput{
// 		Account:      pgconv.PGTextToStringPtr(in.Account),
// 		StockCode:    pgconv.PGTextToStringPtr(in.StockCode),
// 		Location:     pgconv.PGTextToStringPtr(in.Location),
// 		Bin:          pgconv.PGTextToStringPtr(in.Bin),
// 		LotNumber:    pgconv.PGTextToStringPtr(in.LotNumber),
// 		LTETimestamp: pgconv.PGTimestamptzToTimePtr(in.LTETimestamp),
// 		Page:         in.Page,
// 		PageSize:     in.PageSize,
// 	}
// }

type StockLevel struct {
	Account    string
	StockCode  string
	Location   string
	Bin        string
	LotNumber  string
	StockLevel decimal.Decimal
	Timestamp  time.Time
}

type Movement struct {
	StockCode     string
	Qty           decimal.Decimal
	FromLocation  string
	FromBin       string
	FromLotNumber pgtype.Text
	ToLocation    string
	ToBin         string
	ToLotNumber   pgtype.Text
}

type Movements []Movement
