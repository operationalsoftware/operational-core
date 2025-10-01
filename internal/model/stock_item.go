package model

import (
	"app/pkg/appsort"
	"time"
)

type StockItem struct {
	StockItemID int
	StockCode   string `sortable:"true"`
	Description string `sortable:"true"`
	GalleryID   int
	CreatedAt   time.Time `sortable:"true"`
}

type StockItemChange struct {
	StockItemID      int
	StockCode        *string
	Description      *string
	ChangeByUsername string
	ChangedAt        time.Time
	IsCreation       bool
}

type PostStockItemChange struct {
	StockItemID int
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
	GalleryID   int
}

type GetStockItemsQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int
}
