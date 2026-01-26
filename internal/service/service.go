package service

import (
	"app/internal/model"
	"app/internal/repository"
	"app/pkg/appsort"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

var ErrResourceServiceNotFound = errors.New("resource service not found")
var ErrResourceServiceNotLast = errors.New("resource service is not the most recent")

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

func (s *ServicesService) CreateServiceSchedule(
	ctx context.Context,
	schedule model.NewServiceSchedule,
) error {

	if !schedule.Threshold.GreaterThan(decimal.Zero) {
		return fmt.Errorf("threshold must be greater than zero")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	metric, err := s.servicesRepository.GetResourceServiceMetricByID(ctx, tx, schedule.ResourceServiceMetricID)
	if err != nil {
		return err
	}
	if metric == nil {
		return fmt.Errorf("service metric not found")
	}
	if metric.IsArchived {
		return fmt.Errorf("service metric is archived")
	}

	_, err = s.servicesRepository.CreateServiceSchedule(ctx, tx, schedule)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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

func (s *ServicesService) AssignServiceSchedule(
	ctx context.Context,
	serviceSchedule model.NewServiceScheduleAssignment,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // Ensures rollback on error

	schedule, err := s.servicesRepository.GetServiceScheduleByID(ctx, tx, serviceSchedule.ServiceScheduleID)
	if err != nil {
		return err
	}
	if schedule == nil {
		return fmt.Errorf("service schedule not found")
	}
	if schedule.IsArchived {
		return fmt.Errorf("service schedule is archived")
	}

	metric, err := s.servicesRepository.GetResourceServiceMetricByID(ctx, tx, schedule.ResourceServiceMetricID)
	if err != nil {
		return err
	}
	if metric == nil {
		return fmt.Errorf("service metric not found")
	}
	if metric.IsArchived {
		return fmt.Errorf("service metric is archived")
	}

	_, err = s.servicesRepository.AssignServiceSchedule(ctx, tx, serviceSchedule)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) BulkEditServiceSchedules(
	ctx context.Context,
	input model.BulkEditServiceSchedulesInput,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, scheduleID := range input.AssignServiceScheduleIDs {
		schedule, err := s.servicesRepository.GetServiceScheduleByID(ctx, tx, scheduleID)
		if err != nil {
			return err
		}
		if schedule == nil {
			return fmt.Errorf("service schedule not found")
		}
		if schedule.IsArchived {
			return fmt.Errorf("service schedule is archived")
		}

		metric, err := s.servicesRepository.GetResourceServiceMetricByID(ctx, tx, schedule.ResourceServiceMetricID)
		if err != nil {
			return err
		}
		if metric == nil {
			return fmt.Errorf("service metric not found")
		}
		if metric.IsArchived {
			return fmt.Errorf("service metric is archived")
		}
	}

	for _, resourceID := range input.ResourceIDs {
		for _, scheduleID := range input.AssignServiceScheduleIDs {
			_, err := s.servicesRepository.AssignServiceSchedule(ctx, tx, model.NewServiceScheduleAssignment{
				ResourceID:        resourceID,
				ServiceScheduleID: scheduleID,
			})
			if err != nil {
				return err
			}
		}

		for _, scheduleID := range input.UnassignServiceScheduleIDs {
			if err := s.servicesRepository.UnassignServiceSchedule(ctx, tx, resourceID, scheduleID); err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *ServicesService) UnassignServiceSchedule(
	ctx context.Context,
	resourceID int,
	scheduleID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = s.servicesRepository.UnassignServiceSchedule(ctx, tx, resourceID, scheduleID)
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

func (s *ServicesService) HasNewerServiceForResource(
	ctx context.Context,
	resourceID int,
	serviceID int,
	startedAt time.Time,
) (bool, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	hasNewer, err := s.servicesRepository.HasNewerServiceForResource(
		ctx,
		tx,
		resourceID,
		serviceID,
		startedAt,
	)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return hasNewer, nil
}

func (s *ServicesService) DeleteResourceService(
	ctx context.Context,
	serviceID int,
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

	hasNewer, err := s.servicesRepository.HasNewerServiceForResource(
		ctx,
		tx,
		serviceRecord.ResourceID,
		serviceRecord.ResourceServiceID,
		serviceRecord.StartedAt,
	)
	if err != nil {
		return err
	}
	if hasNewer {
		return ErrResourceServiceNotLast
	}

	if err := s.resourceRepository.ReopenClosedMetricRecords(
		ctx,
		tx,
		serviceRecord.ResourceID,
		serviceRecord.ResourceServiceID,
	); err != nil {
		return err
	}

	if err := s.servicesRepository.DeleteResourceService(ctx, tx, serviceID); err != nil {
		return err
	}

	if err := s.galleryRepository.DeleteGallery(ctx, tx, serviceRecord.GalleryID); err != nil {
		return err
	}

	if err := s.commentRepository.DeleteCommentThread(ctx, tx, serviceRecord.CommentThreadID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ServicesService) GetLastServiceForResource(
	ctx context.Context,
	resourceID int,
	excludeServiceID int,
	beforeStartedAt time.Time,
) (*model.ResourceService, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	service, err := s.servicesRepository.GetLastServiceForResource(
		ctx,
		tx,
		resourceID,
		excludeServiceID,
		beforeStartedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
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
	sort appsort.Sort,
) ([]model.ServiceMetric, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.ServiceMetric{}, 0, err
	}
	defer tx.Rollback(ctx)

	metrics, err := s.servicesRepository.ListServiceMetrics(ctx, tx, includeArchived, sort)
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

func (s *ServicesService) GetServiceSchedules(
	ctx context.Context,
	includeArchived bool,
	sort appsort.Sort,
) ([]model.ServiceSchedule, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.ServiceSchedule{}, 0, err
	}
	defer tx.Rollback(ctx)

	schedules, err := s.servicesRepository.ListServiceSchedules(ctx, tx, includeArchived, sort)
	if err != nil {
		return []model.ServiceSchedule{}, 0, err
	}

	count, err := s.servicesRepository.GetServiceSchedulesCount(ctx, tx, includeArchived)
	if err != nil {
		return []model.ServiceSchedule{}, 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return []model.ServiceSchedule{}, 0, err
	}

	return schedules, count, nil
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

func (s *ServicesService) GetServiceScheduleByID(
	ctx context.Context,
	scheduleID int,
) (*model.ServiceSchedule, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	schedule, err := s.servicesRepository.GetServiceScheduleByID(ctx, tx, scheduleID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return schedule, nil
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

func (s *ServicesService) UpdateServiceSchedule(
	ctx context.Context,
	update model.UpdateServiceSchedule,
) error {

	if !update.Threshold.GreaterThan(decimal.Zero) {
		return fmt.Errorf("threshold must be greater than zero")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	metric, err := s.servicesRepository.GetResourceServiceMetricByID(ctx, tx, update.ResourceServiceMetricID)
	if err != nil {
		return err
	}
	if metric == nil {
		return fmt.Errorf("service metric not found")
	}
	if metric.IsArchived {
		return fmt.Errorf("service metric is archived")
	}

	err = s.servicesRepository.UpdateServiceSchedule(ctx, tx, update)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ServicesService) GetResourceServiceMetricStatuses(
	ctx context.Context,
	userID int,
	q model.ResourceServiceMetricStatusesQuery,
) ([]model.ResourceServiceMetricStatus, int, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.ResourceServiceMetricStatus{}, 0, err
	}
	defer tx.Rollback(ctx)

	resources, err := s.servicesRepository.GetResourceServiceMetricStatuses(ctx, tx, userID, q)
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
