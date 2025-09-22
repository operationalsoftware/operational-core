package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/apphmac"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ncw/swift/v2"
)

type GalleryService struct {
	db                *pgxpool.Pool
	swiftConn         *swift.Connection
	fileRepository    *repository.FileRepository
	galleryRepository *repository.GalleryRepository
}

func NewGalleryService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	fileRepository *repository.FileRepository,
	galleryRepository *repository.GalleryRepository,
) *GalleryService {
	return &GalleryService{
		db:                db,
		swiftConn:         swiftConn,
		fileRepository:    fileRepository,
		galleryRepository: galleryRepository,
	}
}

func (s *GalleryService) GetGallery(
	ctx context.Context,
	galleryID int,
) (model.Gallery, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return model.Gallery{}, err
	}
	defer tx.Rollback(ctx)

	gallery, err := s.galleryRepository.GetGalleryByID(
		ctx,
		tx,
		s.swiftConn,
		galleryID,
	)
	if err != nil {
		return model.Gallery{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Gallery{}, err
	}

	return gallery, nil
}

func (s *GalleryService) GetGalleryImgURLs(
	ctx context.Context,
	galleryID int,
) ([]string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	gallery, err := s.galleryRepository.GetGalleryByID(
		ctx,
		tx,
		s.swiftConn,
		galleryID,
	)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(gallery.Items))
	for _, item := range gallery.Items {
		urls = append(urls, item.DownloadURL)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *GalleryService) AddGalleryItem(
	ctx context.Context,
	galleryItem *model.NewGalleryItem,
	userID int,
) (*model.File, string, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback(ctx)

	newFile, signedURL, err := s.fileRepository.CreateFile(
		ctx,
		tx,
		s.swiftConn,
		&model.File{
			Filename:    galleryItem.Filename,
			ContentType: galleryItem.ContentType,
			SizeBytes:   galleryItem.SizeBytes,
			Entity:      "Gallery",
			EntityID:    galleryItem.GalleryID,
		},
		userID,
	)
	if err != nil {
		return nil, "", err
	}

	err = s.galleryRepository.CreateGalleryItem(
		ctx,
		tx,
		galleryItem.GalleryID,
		newFile.FileID,
		userID,
	)
	if err != nil {
		return nil, "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", err
	}

	return newFile, signedURL, nil
}

func (s *GalleryService) SetGalleryItemPosition(
	ctx context.Context,
	galleryID int,
	galleryItemID int,
	newPosition int,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.galleryRepository.UpdateGalleryItemPosition(ctx, tx, galleryID, galleryItemID, newPosition)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *GalleryService) DeleteGalleryItem(
	ctx context.Context,
	galleryID int,
	galleryItemID int,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	galleryItem, err := s.galleryRepository.GetGalleryItemByID(
		ctx,
		tx,
		galleryItemID,
	)
	if err != nil {
		return err
	}

	err = s.galleryRepository.DeleteGalleryItem(
		ctx,
		tx,
		galleryID,
		galleryItemID,
	)
	if err != nil {
		return err
	}

	err = s.fileRepository.DeleteFile(
		ctx,
		tx,
		s.swiftConn,
		galleryItem.FileID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *GalleryService) GenerateTempURL(
	galleryID int,
	canEdit bool,
) string {
	expires := time.Now().Add(12 * time.Hour).Unix()
	return s.generateTempURL(galleryID, expires, canEdit, "")
}

func (s *GalleryService) GenerateEditTempURL(
	galleryID int,
	canEdit bool,
) string {
	expires := time.Now().Add(12 * time.Hour).Unix()
	return s.generateTempURL(galleryID, expires, canEdit, "/edit")
}

func (s *GalleryService) generateTempURL(
	galleryID int,
	expires int64,
	canEdit bool,
	pathSuffix string,
) string {
	secretKey := os.Getenv("AES_256_ENCRYPTION_KEY")

	allowedOperations := []string{"view"}
	if canEdit {
		allowedOperations = append(allowedOperations, "edit")
	}
	ops := strings.Join(allowedOperations, ",")

	claims := apphmac.Claims{
		Entity:            "gallery",
		EntityID:          fmt.Sprintf("%d", galleryID),
		AllowedOperations: allowedOperations,
		Expires:           expires,
	}

	hmac := apphmac.GenerateHMAC(claims, secretKey)

	galleryURL := fmt.Sprintf(
		"/gallery/%d%s?HMAC=%s&AllowedOperations=%s&Expires=%d",
		galleryID, pathSuffix, hmac, ops, expires)

	return galleryURL
}
