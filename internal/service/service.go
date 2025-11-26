package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrResourceServiceNotFound = errors.New("resource service not found")

type ServicesService struct {
	db                 *pgxpool.Pool
	commentRepository  *repository.CommentRepository
	galleryRepository  *repository.GalleryRepository
	resourceRepository *repository.ResourceRepository
	servicesRepository *repository.ServiceRepository
}

func NewServicesService(
	db *pgxpool.Pool,
	commentRepository *repository.CommentRepository,
	galleryRepository *repository.GalleryRepository,
	resourceRepository *repository.ResourceRepository,
	servicesRepository *repository.ServiceRepository,
) *ServicesService {
	return &ServicesService{
		db:                 db,
		commentRepository:  commentRepository,
		galleryRepository:  galleryRepository,
		resourceRepository: resourceRepository,
		servicesRepository: servicesRepository,
	}
}

func (s *ServicesService) CreateResourceServiceMetric(
	ctx context.Context,
	metric model.NewServiceMetric,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	_, err = s.servicesRepository.CreateResourceServiceMetric(ctx, tx, metric)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) DeleteResourceServiceMetric(
	ctx context.Context,
	metricID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.servicesRepository.DeleteResourceServiceMetric(ctx, tx, metricID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) GetActiveServiceID(
	ctx context.Context,
	resourceID int,
) (int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	serviceID, err := s.servicesRepository.GetActiveServiceID(ctx,
		tx,
		resourceID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return serviceID, nil
}

func (s *ServicesService) UpdateResourceService(
	ctx context.Context,
	service model.UpdateResourceService,
	userID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	serviceRecord, err := s.servicesRepository.GetResourceServiceByID(ctx, tx, service.ResourceServiceID)
	if err != nil {
		return err
	}
	if serviceRecord == nil {
		return ErrResourceServiceNotFound
	}

	err = s.servicesRepository.UpdateService(ctx, tx, service, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) CreateResourceServiceSchedule(
	ctx context.Context,
	serviceSchedule model.NewResourceServiceSchedule,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	if !serviceSchedule.Threshold.GreaterThan(decimal.Zero) {
		return fmt.Errorf("threshold must be greater than zero")
	}

	metric, err := s.servicesRepository.GetResourceServiceMetricByID(ctx, tx, serviceSchedule.ResourceServiceMetricID)
	if err != nil {
		return err
	}
	if metric == nil || metric.IsArchived {
		return fmt.Errorf("service metric not found or archived")
	}

	_, err = s.servicesRepository.CreateServiceSchedule(ctx, tx, model.NewResourceServiceSchedule{
		ResourceID:              serviceSchedule.ResourceID,
		ResourceServiceMetricID: serviceSchedule.ResourceServiceMetricID,
		Threshold:               serviceSchedule.Threshold,
	})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) ArchiveResourceServiceSchedule(
	ctx context.Context,
	resourceID int,
	scheduleID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.servicesRepository.ArchiveServiceSchedule(ctx, tx, resourceID, scheduleID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) UpdateService(
	ctx context.Context,
	resourceID int,
	serviceID int,
	action string,
	userID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	serviceRecord, err := s.servicesRepository.GetResourceServiceByID(ctx, tx, serviceID)
	if err != nil {
		return err
	}
	if serviceRecord == nil {
		return ErrResourceServiceNotFound
	}
	if serviceRecord.ResourceID != resourceID {
		return ErrResourceServiceNotFound
	}

	switch action {
	case "complete":
		err = s.servicesRepository.CompleteService(
			ctx,
			tx,
			serviceID,
			userID,
		)
	case "reopen":
		err = s.servicesRepository.ReopenService(
			ctx,
			tx,
			serviceID,
			userID,
		)
	case "cancel":
		err = s.servicesRepository.CancelService(
			ctx,
			tx,
			serviceID,
			userID,
		)
		if err == nil {
			err = s.resourceRepository.ReopenClosedMetricRecords(
				ctx,
				tx,
				resourceID,
				serviceID)
		}
	default:
		return fmt.Errorf("unsupported service action %q", action)
	}

	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) GetResourceServiceByID(
	ctx context.Context,
	serviceID int,
) (
	*model.ResourceService,
	error,
) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	service, err := s.servicesRepository.GetResourceServiceByID(ctx, tx, serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrResourceServiceNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *ServicesService) GetServiceChangelog(
	ctx context.Context,
	serviceID int,
) (
	[]model.ResourceServiceChange,
	error,
) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.ResourceServiceChange{}, err
	}
	defer tx.Rollback(ctx)

	changelog, err := s.servicesRepository.GetServiceChangelog(ctx, tx, serviceID)
	if err != nil {
		return []model.ResourceServiceChange{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.ResourceServiceChange{}, err
	}

	return changelog, nil
}

func (s *ServicesService) GetServiceMetrics(
	ctx context.Context,
	includeArchived bool,
) ([]model.ServiceMetric, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.ServiceMetric{}, 0, err
	}
	defer tx.Rollback(ctx)

	metrics, err := s.servicesRepository.ListServiceMetrics(ctx, tx, includeArchived)
	if err != nil {
		return []model.ServiceMetric{}, 0, err
	}

	count, err := s.servicesRepository.GetMetricsCount(ctx, tx, includeArchived)
	if err != nil {
		return []model.ServiceMetric{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.ServiceMetric{}, 0, err
	}

	return metrics, count, nil
}

func (s *ServicesService) GetResourceServiceMetricByID(
	ctx context.Context,
	metricID int,
) (*model.ServiceMetric, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	metric, err := s.servicesRepository.GetResourceServiceMetricByID(ctx, tx, metricID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *ServicesService) UpdateResourceServiceMetric(
	ctx context.Context,
	update model.UpdateResourceServiceMetric,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.servicesRepository.UpdateResourceServiceMetric(ctx, tx, update)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ServicesService) GetResourceServiceMetricStatuses(
	ctx context.Context,
	q model.ResourceServiceMetricStatusesQuery,
) ([]model.ResourceServiceMetricStatus, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.ResourceServiceMetricStatus{}, 0, err
	}
	defer tx.Rollback(ctx)

	resources, err := s.servicesRepository.GetResourceServiceMetricStatuses(ctx, tx, q)
	if err != nil {
		return []model.ResourceServiceMetricStatus{}, 0, err
	}

	count, err := s.servicesRepository.GetResourceServiceMetricStatusCount(ctx, tx, q)
	if err != nil {
		return []model.ResourceServiceMetricStatus{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.ResourceServiceMetricStatus{}, 0, err
	}

	return resources, count, nil
}

func (s *ServicesService) ListServices(
	ctx context.Context,
	q model.GetServicesQuery,
) ([]model.ResourceService, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.ResourceService{}, 0, err
	}
	defer tx.Rollback(ctx)

	services, err := s.servicesRepository.ListResourceServices(ctx, tx, q)
	if err != nil {
		return []model.ResourceService{}, 0, err
	}

	count, err := s.servicesRepository.ListResourceServicesCount(ctx, tx, q)
	if err != nil {
		return []model.ResourceService{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.ResourceService{}, 0, err
	}

	return services, count, nil
}
