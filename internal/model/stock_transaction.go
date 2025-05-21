package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type StockTransactionEntry struct {
	StockTransactionEntryID int
	StockTransactionType    string
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
	Account      *string
	StockCode    *string
	Location     *string
	Bin          *string
	LotNumber    *string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
}

type NewStockTransaction struct {
	Timestamp            *time.Time
	StockCode            string
	Qty                  decimal.Decimal
	FromAccount          string
	FromLocation         string
	FromBin              string
	FromLotNumber        *string
	ToAccount            string
	ToLocation           string
	ToBin                string
	ToLotNumber          *string
	StockTransactionType string
	StockTransactionNote string
}

type PostStockTransactionsInput []NewStockTransaction

type GetStockLevelsInput struct {
	Account      *string
	StockCode    *string
	Location     *string
	Bin          *string
	LotNumber    *string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
}

type StockLevel struct {
	Account    string
	StockCode  string
	Location   string
	Bin        string
	LotNumber  sql.NullString
	StockLevel decimal.Decimal
	Timestamp  time.Time
}
