package model

import (
	"app/pkg/appsort"
	"time"
)

type StockItemDB struct {
	StockCode   string `db:"stock_code"`
	Description string `db:"description"`
	CreatedAt   time.Time
}

type StockItem struct {
	ProductType string
	YarnType    string
	StyleNumber string
	Colour      string
	Size        string
	ToeClosing  string
	StockCode   string
	Description string
	CreatedAt   time.Time
}

type StockItemChange struct {
	StockCode        string
	StockCodeHistory *string
	Description      *string
	ChangeByUsername string
	ChangeAt         time.Time
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
	ProductType string
	YarnType    string
	StyleNumber string
	Colour      string
	Size        string
	ToeClosing  string
	StockCode   string
	Description string
}

type GetStockItemsQuery struct {
	Sort     appsort.Sort
	Page     int
	PageSize int
}

type SKUConfigItem struct {
	SKUField string
	Label    string
	Code     string
}

type SKUConfig struct {
	Label string
	Code  string
}

type SKUConfigData struct {
	ProductType []SKUConfig
	YarnType    []SKUConfig
	StyleNumber []SKUConfig
	Colour      []SKUConfig
	ToeClosing  []SKUConfig
	Size        []SKUConfig
}
