package model

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

var StockAccounts = []string{"STOCK", "PRODUCTION", "CONSUMED"}

type StockTransactionEntry struct {
	StockTransactionEntryID int
	StockTransactionType    string
	Account                 string
	StockCode               string
	Location                string
	Bin                     string
	Quantity                decimal.Decimal
	LotNumber               pgtype.Text
	RunningTotal            decimal.Decimal
	TransactionBy           string
	TransactionByUsername   string
	Timestamp               time.Time
	StockTransactionID      int
}

type GetTransactionsInput struct {
	Account      pgtype.Text
	StockCode    pgtype.Text
	Location     pgtype.Text
	Bin          pgtype.Text
	LotNumber    pgtype.Text
	LTETimestamp pgtype.Timestamptz
	Page         int
	PageSize     int
}

type NewStockTransaction struct {
	Timestamp            time.Time
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
	Account      pgtype.Text
	StockCode    pgtype.Text
	Location     pgtype.Text
	Bin          pgtype.Text
	LotNumber    pgtype.Text
	LTETimestamp pgtype.Timestamptz
	Page         int
	PageSize     int
}

type StockLevel struct {
	Account    string
	StockCode  string
	Location   string
	Bin        string
	LotNumber  pgtype.Text
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
