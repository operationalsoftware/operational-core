package model

import (
	"app/pkg/appsort"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type StockItem struct {
	StockItemID int
	StockCode   string
	Description string
	CreatedAt   time.Time
}

type StockItemChange struct {
	StockItemID      int
	StockCode        pgtype.Text
	Description      pgtype.Text
	ChangeByUsername string
	ChangedAt        time.Time
	IsCreation       bool
}

type PostStockItemChange struct {
	StockCode   *string
	Description *string
	ChangeBy    int
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
