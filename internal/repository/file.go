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

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) getSignedUploadURL(conn *swift.Connection, container, objectName, secretKey string, expiresIn time.Duration) (string, error) {
	if secretKey == "" {
		return "", fmt.Errorf("TempURL key is not set for container %s", container)
	}

	expires := time.Now().Add(expiresIn)

	// method "PUT" for upload
	uploadURL := conn.ObjectTempUrl(container, objectName, secretKey, "PUT", expires)
	return uploadURL, nil
}

func (r *FileRepository) GetSignedDownloadURL(
	ctx context.Context,
	conn *swift.Connection,
	exec db.PGExecutor,
	fileID string,
	expiresIn time.Duration,
) (string, error) {
	container := os.Getenv("ORBIT_CONTAINER")
	secretKey := os.Getenv("SWIFT_TEMP_URL_KEY")

	file, err := r.GetFileByID(ctx, exec, fileID)
	if err != nil {
		return "", err
	}

	expires := time.Now().Add(expiresIn)

	// method "GET" for download
	downloadURL := conn.ObjectTempUrl(container, file.ObjectName, secretKey, "GET", expires)
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
	object_name,
	original_filename,
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
		&file.ObjectName,
		&file.OriginalFilename,
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
) (*model.File, string, error) {
	container := os.Getenv("ORBIT_CONTAINER")
	secretKey := os.Getenv("SWIFT_TEMP_URL_KEY")
	// 1. Generate unique object name
	objectName := fmt.Sprintf("%s-%s", uuid.New().String(), f.OriginalFilename)

	// 2. Insert metadata into DB
	var fileID string
	err := exec.QueryRow(ctx, `
INSERT INTO file (
	object_name,
	original_filename,
	content_type,
	size_bytes,
	entity,
	user_id
)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
RETURNING file_id`,
		objectName, f.OriginalFilename, f.ContentType, f.SizeBytes, f.Entity, f.UserID,
	).Scan(&fileID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to insert file metadata: %w", err)
	}

	file := &model.File{
		FileID:           fileID,
		ObjectName:       objectName,
		OriginalFilename: f.OriginalFilename,
		ContentType:      f.ContentType,
		SizeBytes:        f.SizeBytes,
		Entity:           f.Entity,
		UserID:           f.UserID,
		CreatedAt:        time.Now(),
	}

	// 3. Generate signed upload URL
	uploadURL, err := r.getSignedUploadURL(conn, container, objectName, secretKey, 15*time.Minute)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate signed upload URL: %w", err)
	}

	return file, uploadURL, nil
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
	object_name
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
