package model

import (
	"app/pkg/appsort"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type StockItemDB struct {
	StockCode   string `db:"stock_code"`
	Description string `db:"description"`
	CreatedAt   time.Time
}

type StockItem struct {
	StockCode   string
	Description string
	CreatedAt   time.Time
}

type StockItemChange struct {
	StockCode        string
	StockCodeHistory pgtype.Text
	Description      pgtype.Text
	ChangeByUsername string
	ChangedAt        time.Time
	IsCreation       bool
}

type PostStockItemChange struct {
	StockCode        string
	StockCodeHistory *string
	Description      *string
	ChangeBy         int
}

type LabelGenerator struct {
	StockCode  string
	LabelCount int
}

type PostStockItem struct {
	StockCode   string
	Description string
}

type GetStockItemsQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int
}
