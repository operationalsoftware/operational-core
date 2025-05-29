package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type StockAccount string

const (
	StockStockAccount      StockAccount = "STOCK"
	ProductionStockAccount StockAccount = "PRODUCTION"
	ConsumedStockAccount   StockAccount = "CONSUMED"
)

var StockAccounts = []StockAccount{
	StockStockAccount,
	ProductionStockAccount,
	ConsumedStockAccount,
}

type StockTransactionType string

const (
	StockMovementTransactionType       StockTransactionType = "Stock Movement"
	ProductionTransactionType          StockTransactionType = "Production"
	ProductionReversalTransactionType  StockTransactionType = "Production Reversal"
	ConsumptionTransactionType         StockTransactionType = "Consumption"
	ConsumptionReversalTransactionType StockTransactionType = "Consumption Reversal"
)

var StockTransacationTypeMap = map[StockTransactionType]struct {
	From StockAccount
	To   StockAccount
}{
	StockMovementTransactionType: {
		From: StockStockAccount,
		To:   StockStockAccount,
	},
	ProductionTransactionType: {
		From: ProductionStockAccount,
		To:   StockStockAccount,
	},
	ProductionReversalTransactionType: {
		From: StockStockAccount,
		To:   ProductionStockAccount,
	},
	ConsumptionTransactionType: {
		From: StockStockAccount,
		To:   ConsumedStockAccount,
	},
	ConsumptionReversalTransactionType: {
		From: ConsumedStockAccount,
		To:   StockStockAccount,
	},
}

type StockTransactionEntry struct {
	StockTransactionEntryID int
	TransactionType         StockTransactionType
	Account                 StockAccount
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
	Account      StockAccount
	StockCode    string
	Location     string
	Bin          string
	LotNumber    string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
}

type NewStockTransaction struct {
	TransactionType StockTransactionType
	Timestamp       *time.Time
	StockCode       string
	Qty             decimal.Decimal
	FromLocation    string
	FromBin         string
	FromLotNumber   string
	ToLocation      string
	ToBin           string
	ToLotNumber     string
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
	Account      StockAccount
	StockCode    string
	Location     string
	Bin          string
	LotNumber    string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
}

type StockLevel struct {
	Account    StockAccount
	StockCode  string
	Location   string
	Bin        string
	LotNumber  string
	StockLevel decimal.Decimal
	Timestamp  time.Time
}
