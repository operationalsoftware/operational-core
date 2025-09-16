package model

import "time"

type Gallery struct {
	GalleryID int
	Items     []GalleryItem
	CreatedAt time.Time
	CreatedBy int
}

type GalleryItem struct {
	GalleryItemID int       `json:"gallery_item_id"`
	GalleryID     int       `json:"gallery_id"`
	Position      int       `json:"position"`
	FileID        string    `json:"file_id"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     int       `json:"created_by"`
	DownloadURL   string
}

type NewGalleryItem struct {
	GalleryID   int
	Filename    string
	ContentType string
	SizeBytes   int
}

type UpdateGalleryItem struct {
	GalleryItemID int `json:"gallery_item_id"`
	NewPosition   int `json:"new_position"`
}
