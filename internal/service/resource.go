package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/validate"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResourceService struct {
	db                 *pgxpool.Pool
	commentRepository  *repository.CommentRepository
	galleryRepository  *repository.GalleryRepository
	resourceRepository *repository.ResourceRepository
	servicesRepository *repository.ServiceRepository
}

func NewResourceService(
	db *pgxpool.Pool,
	commentRepository *repository.CommentRepository,
	galleryRepository *repository.GalleryRepository,
	resourceRepository *repository.ResourceRepository,
	servicesRepository *repository.ServiceRepository,
) *ResourceService {
	return &ResourceService{
		db:                 db,
		commentRepository:  commentRepository,
		galleryRepository:  galleryRepository,
		resourceRepository: resourceRepository,
		servicesRepository: servicesRepository,
	}
}

func (s *ResourceService) CreateResource(
	ctx context.Context,
	resource model.NewResource,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	_, err = s.resourceRepository.CreateResource(ctx, tx, resource)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ResourceService) CreateResourceMetricRecord(
	ctx context.Context,
	record model.NewResourceServiceMetricRecord,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	err = s.resourceRepository.CreateResourceMetricRecord(ctx, tx, record)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ResourceService) CreateResourceService(
	ctx context.Context,
	service model.NewResourceService,
	userID int,
) (int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	galleryId, err := s.galleryRepository.CreateGallery(
		ctx,
		tx,
		userID,
	)
	if err != nil {
		return 0, err
	}

	service.GalleryID = galleryId

	// Create a dedicated comment thread for this resource service
	threadID, err := s.commentRepository.CreateCommentThread(ctx, tx)
	if err != nil {
		return 0, err
	}
	service.CommentThreadID = threadID

	serviceID, err := s.servicesRepository.CreateService(ctx, tx, service, userID)
	if err != nil {
		return 0, err
	}

	err = s.resourceRepository.CloseOpenMetricRecords(
		ctx,
		tx,
		service.ResourceID,
		serviceID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return serviceID, nil
}

func (s *ResourceService) GetResourceByID(
	ctx context.Context,
	resourceID int,
	userID *int,
) (*model.Resource, error) {
	resource, err := s.resourceRepository.GetResourceByID(ctx, s.db, resourceID, userID)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *ResourceService) GetResourceServiceMetrics(
	ctx context.Context,
	resourceID int,
) ([]model.ServiceMetric, error) {
	metrics, err := s.resourceRepository.GetResourceServiceMetrics(ctx, s.db, resourceID)
	if err != nil {
		return []model.ServiceMetric{}, err
	}

	return metrics, nil
}

func (s *ResourceService) GetResourceServiceSchedules(
	ctx context.Context,
	resourceID int,
	userID int,
) (
	[]model.ResourceServiceMetricStatus,
	error,
) {

	currentMetrics, err := s.resourceRepository.ListResourceMetricSchedules(
		ctx,
		s.db,
		resourceID,
		userID,
	)
	if err != nil {
		return []model.ResourceServiceMetricStatus{},
			err
	}

	return currentMetrics, nil
}

func (s *ResourceService) GetServiceMetricLifetimeTotals(
	ctx context.Context,
	resourceID int,
) (
	[]model.ServiceMetricLifetimeTotal,
	error,
) {

	totals, err := s.resourceRepository.ListMetricLifetimeTotals(
		ctx,
		s.db,
		resourceID,
	)
	if err != nil {
		return []model.ServiceMetricLifetimeTotal{},
			err
	}

	return totals, nil
}

func (s *ResourceService) GetResources(
	ctx context.Context,
	q model.GetResourcesQuery,
) ([]model.Resource, int, model.ResourceAvailableFilters, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.Resource{}, 0, model.ResourceAvailableFilters{}, err
	}
	defer tx.Rollback(ctx)

	resources, err := s.resourceRepository.ListResources(ctx, tx, q)
	if err != nil {
		return []model.Resource{}, 0, model.ResourceAvailableFilters{}, err
	}

	count, err := s.resourceRepository.Count(ctx, tx, q)
	if err != nil {
		return []model.Resource{}, 0, model.ResourceAvailableFilters{}, err
	}

	availableFilters, err := s.resourceRepository.GetAvailableFilters(ctx, tx, q)
	if err != nil {
		return []model.Resource{}, 0, model.ResourceAvailableFilters{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.Resource{}, 0, model.ResourceAvailableFilters{}, err
	}

	return resources, count, availableFilters, nil
}

func (s *ResourceService) UpdateResource(
	ctx context.Context,
	resourceID int,
	update model.ResourceUpdate,
) (validate.ValidationErrors, error) {

	validationErrors := make(validate.ValidationErrors)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return validationErrors, err
	}
	defer tx.Rollback(ctx)

	resource, err := s.resourceRepository.GetResourceByID(ctx, tx, resourceID, nil)
	if err != nil {
		return validationErrors, err
	}
	if resource == nil {
		return validationErrors, fmt.Errorf("resource does not exist")
	}

	if !resource.IsArchived && update.IsArchived {
		hasSchedules, err := s.servicesRepository.HasActiveServiceSchedules(ctx, tx, resourceID)
		if err != nil {
			return validationErrors, err
		}
		if hasSchedules {
			validationErrors.Add("IsArchived", "Cannot archive resource with active service schedules")
			return validationErrors, nil
		}
	}

	err = s.resourceRepository.UpdateResource(ctx, tx, resourceID, update)
	if err != nil {
		return validationErrors, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return validationErrors, err
	}

	return validationErrors, nil
}
