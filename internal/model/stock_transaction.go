package model

import (
	"time"

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

type PostStockTransactionsInput []NewStockTransaction

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

type StockLevel struct {
	Account    string
	StockCode  string
	Location   string
	Bin        string
	LotNumber  string
	StockLevel decimal.Decimal
	Timestamp  time.Time
}
