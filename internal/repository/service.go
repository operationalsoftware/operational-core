package repository

import (
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/db"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type ServiceRepository struct{}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{}
}

func (r *ServiceRepository) CreateResourceServiceMetric(
	ctx context.Context,
	exec db.PGExecutor,
	serviceMetric model.NewServiceMetric,
) (int, error) {

	query := `
INSERT INTO resource_service_metric (
	name,
	description,
	is_cumulative
) VALUES ($1, $2, $3)
RETURNING resource_service_metric_id;
	`

	var newID int
	err := exec.QueryRow(
		ctx,
		query,
		serviceMetric.Name,
		serviceMetric.Description,
		serviceMetric.IsCumulative,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *ServiceRepository) DeleteResourceServiceMetric(
	ctx context.Context,
	exec db.PGExecutor,
	metricID int,
) error {

	query := `
DELETE FROM
	resource_service_metric
WHERE
	resource_service_metric_id = $1;
	`

	_, err := exec.Exec(
		ctx,
		query,
		metricID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) CreateServiceSchedule(
	ctx context.Context,
	exec db.PGExecutor,
	newSchedule model.NewServiceSchedule,
) (int, error) {

	query := `
INSERT INTO service_schedule (
	name,
	resource_service_metric_id,
	threshold
) VALUES ($1, $2, $3)
RETURNING service_schedule_id;
	`

	var newID int
	err := exec.QueryRow(
		ctx,
		query,
		newSchedule.Name,
		newSchedule.ResourceServiceMetricID,
		newSchedule.Threshold,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *ServiceRepository) CreateResourceServiceSchedule(
	ctx context.Context,
	exec db.PGExecutor,
	schedule model.NewResourceServiceSchedule,
) (int, error) {

	query := `
INSERT INTO resource_service_schedule (
	resource_id,
	service_schedule_id
) VALUES ($1, $2)
RETURNING resource_service_schedule_id;
	`

	var newID int
	err := exec.QueryRow(ctx, query, schedule.ResourceID, schedule.ServiceScheduleID).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *ServiceRepository) HasActiveServiceSchedules(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
) (bool, error) {

	query := `
	SELECT
		EXISTS (
			SELECT
				1
			FROM
				resource_service_schedule rss
			JOIN service_schedule ss ON ss.service_schedule_id = rss.service_schedule_id
			JOIN resource_service_metric m ON m.resource_service_metric_id = ss.resource_service_metric_id
			WHERE
				rss.resource_id = $1
				AND rss.is_archived = FALSE
				AND ss.is_archived = FALSE
				AND m.is_archived = FALSE
	)
`

	var exists bool
	err := exec.QueryRow(ctx, query, resourceID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ServiceRepository) ArchiveResourceServiceSchedule(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
	scheduleID int,
) error {

	query := `
UPDATE
	resource_service_schedule
SET
	is_archived = TRUE
WHERE
	resource_id = $1
	AND service_schedule_id = $2
	AND is_archived = FALSE
`

	tag, err := exec.Exec(ctx, query, resourceID, scheduleID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("service schedule not found for resource")
	}

	return nil
}

func (r *ServiceRepository) GetActiveServiceID(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
) (int, error) {
	const query = `
SELECT
	resource_service_id
FROM
	resource_service
WHERE
	resource_id = $1
	AND
	status = $2
LIMIT 1;
	`

	var serviceID int
	err := exec.QueryRow(ctx, query,
		resourceID,
		model.ServiceStatusWorkInProgress).Scan(&serviceID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return serviceID, nil
}

func (r *ServiceRepository) CreateService(
	ctx context.Context,
	exec db.PGExecutor,
	event model.NewResourceService,
	userID int,
) (int, error) {

	var exists bool
	checkQuery := `
SELECT EXISTS (
	SELECT 1
	FROM resource_service
	WHERE resource_id = $1
		AND status = $2
);
`
	err := exec.QueryRow(ctx, checkQuery,
		event.ResourceID,
		model.ServiceStatusWorkInProgress).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("resource service already exists")
	}

	query := `
INSERT INTO resource_service (
	resource_id,
	started_by,
	started_at,
	status,
	notes,
	gallery_id,
	comment_thread_id
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING resource_service_id;
	`

	changeQuery := `
INSERT INTO
  resource_service_change (
    change_by,
    resource_service_id,
    started_by
)
VALUES ($1, $2, $1)
`

	var newID int
	err = exec.QueryRow(
		ctx,
		query,
		event.ResourceID,
		userID,
		time.Now(),
		model.ServiceStatusWorkInProgress,
		event.Notes,
		event.GalleryID,
		event.CommentThreadID,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	_, err = exec.Exec(
		ctx, changeQuery,

		userID,
		newID,
	)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *ServiceRepository) UpdateService(
	ctx context.Context,
	exec db.PGExecutor,
	event model.UpdateResourceService,
	userID int,
) error {

	query := `
UPDATE 
	resource_service
SET 
	notes = $2
WHERE
	resource_service_id = $1
	`

	changeQuery := `
INSERT INTO
  resource_service_change (
    change_by,
    resource_service_id,
	notes
)
VALUES ($1, $2, $3)
`

	_, err := exec.Exec(
		ctx,
		query,
		event.ResourceServiceID,
		event.Notes,
	)
	if err != nil {
		return err
	}

	_, err = exec.Exec(
		ctx, changeQuery,

		userID,
		event.ResourceServiceID,
		event.Notes,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) CompleteService(
	ctx context.Context,
	exec db.PGExecutor,
	serviceID int,
	userID int,
) error {

	now := time.Now()

	namedParams := map[string]any{
		"resource_service_id": serviceID,
		"user_id":             userID,
		"status":              model.ServiceStatusCompleted,
		"wip_status":          model.ServiceStatusWorkInProgress,
		"now":                 now,
	}

	serviceUpdateQuery, serviceUpdateParams, err := db.BindNamed(`
UPDATE 
	resource_service
SET 
	completed_by = :user_id,
	completed_at = :now,
	status = :status
WHERE
	resource_service_id = :resource_service_id
	AND
	status = :wip_status
`, namedParams)
	if err != nil {
		return fmt.Errorf("error binding andon update params: %v", err)
	}
	_, err = exec.Exec(
		ctx, serviceUpdateQuery,
		serviceUpdateParams...,
	)
	if err != nil {
		return err
	}

	changelogQuery, changelogParams, err := db.BindNamed(`
INSERT INTO
	resource_service_change (
		resource_service_id,
		change_by,
		change_at,
		completed_by
	)
VALUES (
	:resource_service_id,
	:user_id,
	:now,
	:user_id
)
`, namedParams)
	if err != nil {
		return fmt.Errorf("error binding changelog params: %v", err)
	}
	_, err = exec.Exec(
		ctx, changelogQuery,
		changelogParams...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) ReopenService(
	ctx context.Context,
	exec db.PGExecutor,
	serviceID int,
	userID int,
) error {

	now := time.Now()

	namedParams := map[string]any{
		"resource_service_id": serviceID,
		"user_id":             userID,
		"status":              model.ServiceStatusWorkInProgress,
		"now":                 now,
	}

	serviceUpdateQuery, serviceUpdateParams, err := db.BindNamed(`
UPDATE 
	resource_service
SET 
	completed_by = NULL,
	completed_at = NULL,
	cancelled_by = NULL,
	cancelled_at = NULL,
	status = :status
WHERE
	resource_service_id = :resource_service_id
`, namedParams)
	if err != nil {
		return fmt.Errorf("error binding andon update params: %v", err)
	}
	_, err = exec.Exec(
		ctx, serviceUpdateQuery,
		serviceUpdateParams...,
	)
	if err != nil {
		return err
	}

	changelogQuery, changelogParams, err := db.BindNamed(`
INSERT INTO
	resource_service_change (
		resource_service_id,
		change_by,
		change_at,
		reopened_by
	)
VALUES (
	:resource_service_id,
	:user_id,
	:now,
	:user_id
)
`, namedParams)
	if err != nil {
		return fmt.Errorf("error binding changelog params: %v", err)
	}
	_, err = exec.Exec(
		ctx, changelogQuery,
		changelogParams...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) CancelService(
	ctx context.Context,
	exec db.PGExecutor,
	serviceID int,
	userID int,
) error {

	now := time.Now()

	namedParams := map[string]any{
		"resource_service_id": serviceID,
		"user_id":             userID,
		"status":              model.ServiceStatusCancelled,
		"wip_status":          model.ServiceStatusWorkInProgress,
		"now":                 now,
	}

	andonUpdateQuery, andonUpdateParams, err := db.BindNamed(`
UPDATE 
	resource_service
SET 
	cancelled_by = :user_id,
	cancelled_at = :now,
	status = :status
WHERE
	resource_service_id = :resource_service_id
	AND
	status = :wip_status
`, namedParams)
	if err != nil {
		return fmt.Errorf("error binding andon update params: %v", err)
	}
	_, err = exec.Exec(
		ctx, andonUpdateQuery,
		andonUpdateParams...,
	)
	if err != nil {
		return err
	}

	changelogQuery, changelogParams, err := db.BindNamed(`
INSERT INTO
	resource_service_change (
		resource_service_id,
		change_by,
		change_at,
		cancelled_by
	)
VALUES (
	:resource_service_id,
	:user_id,
	:now,
	:user_id
)
`, namedParams)
	if err != nil {
		return fmt.Errorf("error binding changelog params: %v", err)
	}
	_, err = exec.Exec(
		ctx, changelogQuery,
		changelogParams...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ServiceRepository) GetResourceServiceByID(
	ctx context.Context,
	exec db.PGExecutor,
	serviceID int,
) (*model.ResourceService, error) {
	query := `
SELECT
    rs.resource_service_id,
    rs.resource_id,
    r.type,
    r.reference,
    rs.status,
    au.username AS started_by_username,
    rs.started_at,
    rs.notes,
    rs.gallery_id,
	rs.comment_thread_id
FROM
    resource_service rs
INNER JOIN app_user au ON rs.started_by = au.user_id
INNER JOIN resource r ON rs.resource_id = r.resource_id
WHERE
    resource_service_id = $1
	`

	sm := model.ResourceService{}
	err := exec.QueryRow(ctx, query, serviceID).Scan(
		&sm.ResourceServiceID,
		&sm.ResourceID,
		&sm.ResourceType,
		&sm.ResourceReference,
		&sm.Status,
		&sm.StartedByUsername,
		&sm.StartedAt,
		&sm.Notes,
		&sm.GalleryID,
		&sm.CommentThreadID,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &sm, nil
}

func (r *ServiceRepository) GetLastServiceForResource(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
	excludeServiceID int,
	beforeStartedAt time.Time,
) (*model.ResourceService, error) {
	query := `
SELECT
    resource_id,
    resource_service_id,
    type,
    reference,
    status,
    started_at,
    started_by,
    started_by_username,
    completed_at,
    completed_by_username,
    cancelled_at,
    cancelled_by_username,
	notes,
	gallery_id,
	comment_thread_id
FROM
    resource_service_view
WHERE
    resource_id = $1
	AND resource_service_id <> $2
	AND started_at < $3
	AND cancelled_at IS NULL
ORDER BY
	started_at DESC
LIMIT 1
`

	var service model.ResourceService
	err := exec.QueryRow(ctx, query, resourceID, excludeServiceID, beforeStartedAt).Scan(
		&service.ResourceID,
		&service.ResourceServiceID,
		&service.ResourceType,
		&service.ResourceReference,
		&service.Status,
		&service.StartedAt,
		&service.StartedBy,
		&service.StartedByUsername,
		&service.CompletedAt,
		&service.CompletedByUsername,
		&service.CancelledAt,
		&service.CancelledByUsername,
		&service.Notes,
		&service.GalleryID,
		&service.CommentThreadID,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *ServiceRepository) GetServiceChangelog(
	ctx context.Context,
	exec db.PGExecutor,
	serviceID int,
) ([]model.ResourceServiceChange, error) {

	query := `
SELECT
	resource_service_id,
	resource_service_change_id,
	change_by,
	change_by_username,
	change_at,
    is_creation,
	notes,
	started_by,
	started_by_username,
	completed_by,
	completed_by_username,
	reopened_by,
	reopened_by_username,
	cancelled_by,
	cancelled_by_username
FROM
	resource_service_change_view

WHERE resource_service_id = $1
ORDER BY change_at DESC
`

	rows, err := exec.Query(ctx, query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []model.ResourceServiceChange
	for rows.Next() {
		var change model.ResourceServiceChange
		if err := rows.Scan(
			&change.ResourceServiceID,
			&change.ResourceServiceChangeID,
			&change.ChangeBy,
			&change.ChangeByUsername,
			&change.ChangeAt,
			&change.IsCreation,
			&change.Notes,
			&change.StartedBy,
			&change.StartedByUsername,
			&change.CompletedBy,
			&change.CompletedByUsername,
			&change.ReopenedBy,
			&change.ReopenedByUsername,
			&change.CancelledBy,
			&change.CancelledByUsername,
		); err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return changes, nil
}

func (r *ServiceRepository) ListServiceMetrics(
	ctx context.Context,
	exec db.PGExecutor,
	includeArchived bool,
	sort appsort.Sort,
) ([]model.ServiceMetric, error) {

	query := `
SELECT
    resource_service_metric_id,
    name,
    description,
    is_cumulative,
	is_archived
FROM
    resource_service_metric
`
	if !includeArchived {
		query += `
WHERE
	is_archived = FALSE
`
	}

	orderByClause, err := sort.ToOrderByClause(model.ServiceMetric{})
	if err != nil {
		return nil, err
	}
	if orderByClause == "" {
		orderByClause = "ORDER BY name ASC"
	}

	query += "\n" + orderByClause

	rows, err := exec.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	serviceMetrics := []model.ServiceMetric{}
	for rows.Next() {
		var metric model.ServiceMetric

		err := rows.Scan(
			&metric.ServiceMetricID,
			&metric.Name,
			&metric.Description,
			&metric.IsCumulative,
			&metric.IsArchived,
		)
		if err != nil {
			return nil, err
		}

		serviceMetrics = append(serviceMetrics, metric)
	}

	return serviceMetrics, nil
}

func (r *ServiceRepository) ListServiceSchedules(
	ctx context.Context,
	exec db.PGExecutor,
	includeArchived bool,
	sort appsort.Sort,
) ([]model.ServiceSchedule, error) {

	query := `
SELECT
    service_schedule_id,
    name,
	resource_service_metric_id,
    metric_name,
    threshold,
	is_archived
FROM
    service_schedule_view
`
	if !includeArchived {
		query += `
WHERE
	is_archived = FALSE
		AND metric_is_archived = FALSE
`
	}

	orderByClause, err := sort.ToOrderByClause(model.ServiceSchedule{})
	if err != nil {
		return nil, err
	}
	if orderByClause == "" {
		orderByClause = "ORDER BY name ASC"
	}

	query += "\n" + orderByClause

	rows, err := exec.Query(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	serviceSchedules := []model.ServiceSchedule{}
	for rows.Next() {
		var schedule model.ServiceSchedule

		err := rows.Scan(
			&schedule.ServiceScheduleID,
			&schedule.Name,
			&schedule.ResourceServiceMetricID,
			&schedule.MetricName,
			&schedule.Threshold,
			&schedule.IsArchived,
		)
		if err != nil {
			return nil, err
		}

		serviceSchedules = append(serviceSchedules, schedule)
	}

	return serviceSchedules, nil
}

func (r *ServiceRepository) GetMetricsCount(
	ctx context.Context,
	exec db.PGExecutor,
	includeArchived bool,
) (int, error) {

	query := `
SELECT
	COUNT(*)
FROM
	resource_service_metric
	`
	if !includeArchived {
		query += `
WHERE
	is_archived = FALSE
`
	}

	var count int
	err := exec.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ServiceRepository) GetServiceSchedulesCount(
	ctx context.Context,
	exec db.PGExecutor,
	includeArchived bool,
) (int, error) {

	query := `
SELECT
	COUNT(*)
FROM
	service_schedule_view
`
	if !includeArchived {
		query += `
WHERE
	is_archived = FALSE
	AND metric_is_archived = FALSE
`
	}

	var count int
	err := exec.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ServiceRepository) GetResourceServiceMetricByID(
	ctx context.Context,
	exec db.PGExecutor,
	metricID int,
) (*model.ServiceMetric, error) {

	query := `
SELECT
	resource_service_metric_id,
	name,
	description,
	is_cumulative,
	is_archived
FROM
	resource_service_metric
WHERE
	resource_service_metric_id = $1
`

	var metric model.ServiceMetric
	err := exec.QueryRow(ctx, query, metricID).Scan(
		&metric.ServiceMetricID,
		&metric.Name,
		&metric.Description,
		&metric.IsCumulative,
		&metric.IsArchived,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &metric, nil
}

func (r *ServiceRepository) GetServiceScheduleByID(
	ctx context.Context,
	exec db.PGExecutor,
	scheduleID int,
) (*model.ServiceSchedule, error) {

	query := `
SELECT
	service_schedule_id,
	name,
	resource_service_metric_id,
	metric_name,
	threshold,
	is_archived
FROM
	service_schedule_view
WHERE
	service_schedule_id = $1
`

	var schedule model.ServiceSchedule
	err := exec.QueryRow(ctx, query, scheduleID).Scan(
		&schedule.ServiceScheduleID,
		&schedule.Name,
		&schedule.ResourceServiceMetricID,
		&schedule.MetricName,
		&schedule.Threshold,
		&schedule.IsArchived,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &schedule, nil
}

func (r *ServiceRepository) UpdateResourceServiceMetric(
	ctx context.Context,
	exec db.PGExecutor,
	metric model.UpdateResourceServiceMetric,
) error {

	query := `
UPDATE
	resource_service_metric
SET
	name = $1,
	description = $2,
	is_cumulative = $3,
	is_archived = $4
WHERE
	resource_service_metric_id = $5
`

	ct, err := exec.Exec(ctx, query,
		metric.Name,
		metric.Description,
		metric.IsCumulative,
		metric.IsArchived,
		metric.ServiceMetricID,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("service metric not found")
	}

	return nil
}

func (r *ServiceRepository) UpdateServiceSchedule(
	ctx context.Context,
	exec db.PGExecutor,
	schedule model.UpdateServiceSchedule,
) error {

	query := `
UPDATE
	service_schedule
SET
	name = $1,
	resource_service_metric_id = $2,
	threshold = $3,
	is_archived = $4
WHERE
	service_schedule_id = $5
`

	ct, err := exec.Exec(ctx, query,
		schedule.Name,
		schedule.ResourceServiceMetricID,
		schedule.Threshold,
		schedule.IsArchived,
		schedule.ServiceScheduleID,
	)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("service schedule not found")
	}

	return nil
}

func (r *ServiceRepository) ListResourceServices(
	ctx context.Context,
	exec db.PGExecutor,
	q model.GetServicesQuery,
) ([]model.ResourceService, error) {

	whereClause, args := r.generateWhereClause(q)

	limitPlaceholder := fmt.Sprintf("$%d", len(args)+1)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args)+2)

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(model.ResourceService{})

	if orderByClause == "" {
		orderByClause = "ORDER BY resource_service_id DESC"
	}

	query := `
SELECT
    resource_id,
    resource_service_id,
    type,
    reference,
    status,
    started_at,
    started_by,
    started_by_username,
    completed_at,
    completed_by_username,
    cancelled_at,
    cancelled_by_username,
	notes,
	gallery_id,
	comment_thread_id
FROM
    resource_service_view

` + whereClause + `

` + orderByClause + `
` + fmt.Sprintf("LIMIT %s OFFSET %s", limitPlaceholder, offsetPlaceholder)

	rows, err := exec.Query(ctx, query, append(args, limit, offset)...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	services := []model.ResourceService{}
	for rows.Next() {
		var service model.ResourceService

		err := rows.Scan(
			&service.ResourceID,
			&service.ResourceServiceID,
			&service.ResourceType,
			&service.ResourceReference,
			&service.Status,
			&service.StartedAt,
			&service.StartedBy,
			&service.StartedByUsername,
			&service.CompletedAt,
			&service.CompletedByUsername,
			&service.CancelledAt,
			&service.CancelledByUsername,
			&service.Notes,
			&service.GalleryID,
			&service.CommentThreadID,
		)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) ListResourceServicesCount(
	ctx context.Context,
	exec db.PGExecutor,
	q model.GetServicesQuery,
) (int, error) {

	whereClause, args := r.generateWhereClause(q)

	query := `
SELECT
	COUNT(*)
FROM
	resource_service_view
` + whereClause

	var count int
	err := exec.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ServiceRepository) generateWhereClause(filters model.GetServicesQuery) (string, []any) {
	var whereClauses []string
	var args []any
	argID := 1

	addInClause := func(column string, values []string) {
		if len(values) == 0 {
			return
		}
		placeholders := make([]string, len(values))
		for i, val := range values {
			args = append(args, val)
			placeholders[i] = fmt.Sprintf("$%d", argID)
			argID++
		}
		whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ", ")))
	}

	addInClause("reference", filters.ResourceIn)

	if len(whereClauses) == 0 {
		return "", args
	}
	return "WHERE " + strings.Join(whereClauses, " AND "), args
}

func (r *ServiceRepository) GetResourceServiceMetricStatuses(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ResourceServiceMetricStatusesQuery,
) ([]model.ResourceServiceMetricStatus, error) {

	baseQuery := `
SELECT
  service_schedule_id,
  service_schedule_name,
  resource_id,
  type,
  reference,
  service_ownership_team_id,
  service_ownership_team_name,
  resource_service_metric_id,
  metric_name,
  current_value,
  threshold,
  normalised_value,
  normalised_percentage,
  is_due,
  last_recorded_at,
  last_serviced_at,
  wip_service_id,
  has_wip_service,
	schedule_is_archived,
	metric_is_archived
FROM
    resource_service_metric_status_view
WHERE
	schedule_is_archived = FALSE
	AND metric_is_archived = FALSE
`

	limit := q.PageSize
	if limit == 0 {
		limit = 50
	}
	offset := (q.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	params := map[string]any{
		"limit":  limit,
		"offset": offset,
	}

	if len(q.ServiceOwnershipTeamIDs) > 0 {
		baseQuery += "	AND service_ownership_team_id = ANY(:team_ids)\n"
		params["team_ids"] = q.ServiceOwnershipTeamIDs
	}

	baseQuery += `
ORDER BY normalised_percentage DESC
LIMIT :limit OFFSET :offset
`

	query, args, err := db.BindNamed(baseQuery, params)
	if err != nil {
		return nil, err
	}

	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	resources := []model.ResourceServiceMetricStatus{}
	for rows.Next() {
		var resource model.ResourceServiceMetricStatus

		err := rows.Scan(
			&resource.ServiceScheduleID,
			&resource.ServiceScheduleName,
			&resource.ResourceID,
			&resource.Type,
			&resource.Reference,
			&resource.ServiceOwnershipTeamID,
			&resource.ServiceOwnershipTeamName,
			&resource.ResourceServiceMetricID,
			&resource.MetricName,
			&resource.CurrentValue,
			&resource.Threshold,
			&resource.NormalisedValue,
			&resource.NormalisedPercentage,
			&resource.IsDue,
			&resource.LastRecordedAt,
			&resource.LastServicedAt,
			&resource.WIPServiceID,
			&resource.HasWIPService,
			&resource.ScheduleIsArchived,
			&resource.MetricIsArchived,
		)
		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (r *ServiceRepository) GetResourceServiceMetricStatusCount(
	ctx context.Context,
	exec db.PGExecutor,
	q model.ResourceServiceMetricStatusesQuery,
) (int, error) {

	query := `
SELECT
	COUNT(*)
FROM
	resource_service_metric_status_view
WHERE
	schedule_is_archived = FALSE
	AND metric_is_archived = FALSE
	`

	var args []any
	if len(q.ServiceOwnershipTeamIDs) > 0 {
		query += fmt.Sprintf("	AND service_ownership_team_id = ANY($%d::int[])\n", len(args)+1)
		args = append(args, q.ServiceOwnershipTeamIDs)
	}

	var count int
	err := exec.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
