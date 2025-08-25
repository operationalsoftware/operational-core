package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ncw/swift/v2"
)

type FileRepository struct {
	container string
	secretKey string
}

func NewFileRepository() *FileRepository {
	return &FileRepository{
		container: os.Getenv("ORBIT_CONTAINER"),
		secretKey: os.Getenv("AES_256_ENCRYPTION_KEY"),
	}
}

func (r *FileRepository) getSignedUploadURL(conn *swift.Connection, objectName string, expiresIn time.Duration) (string, error) {
	expires := time.Now().Add(expiresIn)

	// method "PUT" for upload
	uploadURL := conn.ObjectTempUrl(r.container, objectName, r.secretKey, "PUT", expires)
	return uploadURL, nil
}

func (r *FileRepository) GetSignedDownloadURL(
	ctx context.Context,
	conn *swift.Connection,
	exec db.PGExecutor,
	fileID string,
	expiresIn time.Duration,
) (string, error) {

	file, err := r.GetFileByID(ctx, exec, fileID)
	if err != nil {
		return "", err
	}

	expires := time.Now().Add(expiresIn)

	// method "GET" for download
	downloadURL := conn.ObjectTempUrl(r.container, file.StorageKey, r.secretKey, "GET", expires)
	return downloadURL, nil
}

func (r *FileRepository) GetFileByID(
	ctx context.Context,
	exec db.PGExecutor,
	fileID string,
) (*model.File, error) {

	query := `
SELECT
	file_id,
	storage_key,
	filename,
	content_type,
	size_bytes,
	entity,
	user_id
FROM
	file
WHERE
	file_id = $1
`

	var file model.File
	err := exec.QueryRow(
		ctx, query, fileID,
	).Scan(
		&file.FileID,
		&file.StorageKey,
		&file.Filename,
		&file.ContentType,
		&file.SizeBytes,
		&file.Entity,
		&file.UserID,
	)
	if err != nil {
		return nil, err
	}

	return &file, err
}

func (r *FileRepository) CreateFile(
	ctx context.Context,
	exec db.PGExecutor,
	conn *swift.Connection,
	f *model.File,
	userID int,
) (*model.File, string, error) {
	objectName := uuid.New()

	var fileID uuid.UUID
	err := exec.QueryRow(ctx, `
INSERT INTO file (
	file_id,
	storage_key,
	filename,
	content_type,
	size_bytes,
	entity,
	entity_id,
	user_id
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8
)
RETURNING file_id`,
		objectName, objectName.String(), f.Filename, f.ContentType, f.SizeBytes, f.Entity, f.EntityID, userID,
	).Scan(&fileID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to insert file metadata: %w", err)
	}

	file := &model.File{
		FileID:      fileID.String(),
		StorageKey:  objectName.String(),
		Filename:    f.Filename,
		ContentType: f.ContentType,
		SizeBytes:   f.SizeBytes,
		Entity:      f.Entity,
		UserID:      f.UserID,
		CreatedAt:   time.Now(),
	}

	// 3. Generate signed upload URL
	uploadURL, err := r.getSignedUploadURL(conn, objectName.String(), 15*time.Minute)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate signed upload URL: %w", err)
	}

	return file, uploadURL, nil
}

func (r *FileRepository) CompleteFileUpload(
	ctx context.Context,
	exec db.PGExecutor,
	fileID string,
) error {
	// Ensure team exists
	existing, err := r.GetFileByID(ctx, exec, fileID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("team with ID %d not found", fileID)
	}

	query := `
UPDATE file
SET
	status = $2
WHERE
	file_id = $1
`
	_, err = exec.Exec(
		ctx,
		query,
		fileID,
		"success",
	)
	return err
}

func (r *FileRepository) DeleteFile(
	ctx context.Context,
	exec db.PGExecutor,
	conn *swift.Connection,
	fileID, container, secretKey string,
) error {
	// 1. Fetch file metadata from DB
	var objectName string
	query := `
SELECT
	storage_key
FROM
	file
WHERE
	file_id = $1`
	err := exec.QueryRow(ctx, query, fileID).Scan(&objectName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("file not found")
		}
		return err
	}

	// 2. Delete object from Orbit
	err = conn.ObjectDelete(ctx, container, objectName)
	if err != nil {
		return fmt.Errorf("failed to delete object from Orbit: %w", err)
	}

	// 3. Delete file record from DB
	delQuery := `
DELETE FROM
	file
WHERE
	file_id = $1
`
	_, err = exec.Exec(ctx, delQuery, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	return nil
}

func (r *CommentRepository) GetFilesByEntity(
	ctx context.Context,
	exec db.PGExecutor,
	entity string,
	entityID int,
) ([]model.File, error) {

	query := `
SELECT
	file_id,
	filename,
	storage_key,
	content_type,
	size_bytes,
	status,
	user_id
FROM file
WHERE
	entity = $1 AND entity_id = $2
ORDER BY created_at ASC
`

	rows, err := exec.Query(ctx, query, entity, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []model.File
	for rows.Next() {
		var file model.File
		if err := rows.Scan(
			&file.FileID,
			&file.Filename,
			&file.StorageKey,
			&file.ContentType,
			&file.SizeBytes,
			&file.Status,
			&file.UserID,
		); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
