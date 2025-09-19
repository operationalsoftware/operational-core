package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"encoding/json"
	"time"

	"github.com/ncw/swift/v2"
)

type GalleryRepository struct {
	secret   string
	fileRepo *FileRepository
}

func NewGalleryRepository(secret string, fileRepo *FileRepository) *GalleryRepository {
	return &GalleryRepository{
		secret:   secret,
		fileRepo: fileRepo,
	}
}

func (r *GalleryRepository) CreateGallery(
	ctx context.Context,
	exec db.PGExecutor,
	userID int,
) (int, error) {

	query := `
INSERT INTO gallery (
	created_by
)
VALUES ($1)
RETURNING gallery_id
	`
	var newGalleryID int
	err := exec.QueryRow(
		ctx,
		query,
		userID,
	).Scan(&newGalleryID)

	if err != nil {
		return 0, err
	}

	return newGalleryID, nil
}

func (r *GalleryRepository) CreateGalleryItem(
	ctx context.Context,
	exec db.PGExecutor,
	galleryID int,
	fileID string,
	userID int,
) error {

	query := `
INSERT INTO gallery_item (
	gallery_id,
	file_id,
	position,
	created_by
)
VALUES (
	$1,
	$2,
	COALESCE(
		(SELECT MAX(position) + 1 FROM gallery_item WHERE gallery_id = $1),
		1
	),
	$3
)
`
	_, err := exec.Exec(
		ctx,
		query,
		galleryID,
		fileID,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *GalleryRepository) UpdateGalleryItemPosition(
	ctx context.Context,
	exec db.PGExecutor,
	galleryID int,
	galleryItemID int,
	newPosition int,
) error {

	var oldPosition int
	exec.QueryRow(ctx, `
SELECT
	position
FROM
	gallery_item
WHERE
	gallery_item_id = $1
`, galleryItemID).Scan(&oldPosition)

	if newPosition == oldPosition {
		return nil
	}

	var err error

	if newPosition < oldPosition {
		_, err = exec.Exec(ctx, `
UPDATE
	gallery_item
SET
	position = position + 1
WHERE
	gallery_id = $1
	AND position >= $2
	AND position < $3
        `, galleryID, newPosition, oldPosition)
	} else {
		_, err = exec.Exec(ctx, `
UPDATE
	gallery_item
SET
	position = position - 1
WHERE
	gallery_id = $1
	AND position > $2
	AND position <= $3
        `, galleryID, oldPosition, newPosition)
	}

	if err != nil {
		return err
	}

	_, err = exec.Exec(ctx, `
UPDATE
	gallery_item
SET
	position = $2
WHERE
	gallery_id = $1 AND gallery_item_id = $3
`,
		galleryID,
		newPosition,
		galleryItemID)

	return err
}

func (r *GalleryRepository) DeleteGalleryItem(
	ctx context.Context,
	exec db.PGExecutor,
	galleryID int,
	galleryItemID int,
) error {

	var position int
	exec.QueryRow(ctx, `
SELECT
	position
FROM
	gallery_item
WHERE
	gallery_item_id = $1
`, galleryItemID).Scan(&position)

	query := `
DELETE FROM
	gallery_item
WHERE
	gallery_item_id = $1
`
	_, err := exec.Exec(
		ctx,
		query,
		galleryItemID,
	)
	if err != nil {
		return err
	}

	query = `
UPDATE
	gallery_item
SET
	position = position - 1
WHERE
	gallery_id = $1 AND position > $2
`
	_, err = exec.Exec(
		ctx,
		query,
		galleryID,
		position,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *GalleryRepository) GetGalleryByID(
	ctx context.Context,
	exec db.PGExecutor,
	swiftConn *swift.Connection,
	galleryID int,
) (model.Gallery, error) {

	query := `
SELECT
	gallery_id,
	items,
	created_by,
	created_at
FROM
	gallery_view

WHERE
	gallery_id = $1
`

	var g model.Gallery
	var itemsRaw []byte
	err := exec.QueryRow(ctx, query, galleryID).Scan(
		&g.GalleryID,
		&itemsRaw,
		&g.CreatedBy,
		&g.CreatedAt,
	)
	if err != nil {
		return model.Gallery{}, err
	}

	if err := json.Unmarshal(itemsRaw, &g.Items); err != nil {
		return model.Gallery{}, err
	}

	for i := range g.Items {
		url, err := r.fileRepo.GetSignedDownloadURL(
			ctx,
			swiftConn,
			exec,
			g.Items[i].FileID, 24*time.Hour)
		if err != nil {
			return model.Gallery{}, err
		}
		g.Items[i].DownloadURL = url
	}

	return g, nil
}

func (r *GalleryRepository) GetGalleryItemByID(
	ctx context.Context,
	exec db.PGExecutor,
	galleryItemID int,
) (*model.GalleryItem, error) {

	query := `
SELECT
	gallery_item_id,
	gallery_id,
	file_id,
	position,
	created_by,
	created_at
FROM
	gallery_item

WHERE
	gallery_item_id = $1
`

	var gi model.GalleryItem
	err := exec.QueryRow(ctx, query, galleryItemID).Scan(
		&gi.GalleryItemID,
		&gi.GalleryID,
		&gi.FileID,
		&gi.Position,
		&gi.CreatedBy,
		&gi.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &gi, nil
}
