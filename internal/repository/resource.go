package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type ResourceRepository struct{}

func NewResourceRepository() *ResourceRepository {
	return &ResourceRepository{}
}

func (r *ResourceRepository) CreateResource(
	ctx context.Context,
	exec db.PGExecutor,
	resource model.NewResource,
) (int, error) {

	query := `
INSERT INTO resource (
	type,
	reference,
	service_ownership_team_id
) VALUES ($1, $2, $3)
RETURNING resource_id;
	`

	var newID int
	err := exec.QueryRow(
		ctx,
		query,
		resource.Type,
		resource.Reference,
		resource.ServiceOwnershipTeamID,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *ResourceRepository) CreateResourceUsageRecord(
	ctx context.Context,
	exec db.PGExecutor,
	record model.NewResourceUsageRecord,
) error {

	query := `
INSERT INTO resource_usage_record (
	resource_id,
	resource_service_metric_id,
	value
) VALUES ($1, $2, $3)
`

	_, err := exec.Exec(
		ctx,
		query,
		record.ResourceID,
		record.ResourceServiceMetricID,
		record.Value,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceRepository) UpdateResource(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
	update model.ResourceUpdate,
) error {

	resource, err := r.GetResourceByID(ctx, exec, resourceID)
	if err != nil {
		return err
	}
	if resource == nil {
		return fmt.Errorf("resource does not exist")
	}

	query := `
UPDATE
	resource
SET
	type = $1,
	reference = $2,
	is_archived = $3,
	service_ownership_team_id = $4
WHERE
	resource_id = $5
	`

	_, err = exec.Exec(
		ctx,
		query,

		update.Type,
		update.Reference,
		update.IsArchived,
		update.ServiceOwnershipTeamID,
		resourceID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceRepository) CloseOpenUsageRecords(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
	serviceID int,
) error {

	resource, err := r.GetResourceByID(ctx, exec, resourceID)
	if err != nil {
		return err
	}
	if resource == nil {
		return fmt.Errorf("resource does not exist")
	}

	query := `
UPDATE
	resource_usage_record

SET
	closed_by_resource_service_id = $2

WHERE
	resource_id = $1
	AND
	closed_by_resource_service_id IS NULL
`

	_, err = exec.Exec(
		ctx,
		query,

		resourceID,
		serviceID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceRepository) ReopenClosedUsageRecords(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
	serviceID int,
) error {

	resource, err := r.GetResourceByID(ctx, exec, resourceID)
	if err != nil {
		return err
	}
	if resource == nil {
		return fmt.Errorf("resource does not exist")
	}

	query := `
UPDATE
	resource_usage_record
SET
	closed_by_resource_service_id = NULL
WHERE
	resource_id = $1
	AND
	closed_by_resource_service_id = $2
`

	_, err = exec.Exec(
		ctx,
		query,
		resourceID,
		serviceID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceRepository) ArchiveResource(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
) error {

	query := `
UPDATE
	resource
SET
	is_archived = TRUE
WHERE
	resource_id = $1
`

	_, err := exec.Exec(ctx, query, resourceID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceRepository) GetResourceByID(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
) (*model.Resource, error) {
	query := `
SELECT
    resource_id,
    type,
    reference,
    is_archived,
	service_ownership_team_id,
	service_ownership_team_name,
	last_serviced_at
FROM
    resource_view
WHERE
    resource_id = $1
	`

	resource := model.Resource{}
	err := exec.QueryRow(ctx, query, resourceID).Scan(
		&resource.ResourceID,
		&resource.Type,
		&resource.Reference,
		&resource.IsArchived,
		&resource.ServiceOwnershipTeamID,
		&resource.ServiceOwnershipTeamName,
		&resource.LastServicedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *ResourceRepository) ListResources(
	ctx context.Context,
	exec db.PGExecutor,
	q model.GetResourcesQuery,
) ([]model.Resource, error) {

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(model.Resource{})

	if orderByClause == "" {
		orderByClause = "ORDER BY resource_id ASC"
	}

	whereClause := `
WHERE
	is_archived = FALSE
`
	if q.IsArchived {
		whereClause = `
WHERE
	is_archived = TRUE
`
	}

	query := `
SELECT
    resource_id,
    type,
    reference,
    is_archived,
	service_ownership_team_id,
	service_ownership_team_name,
    last_serviced_at
FROM
    resource_view
` + whereClause + `
` + orderByClause + `
LIMIT $1 OFFSET $2
`

	rows, err := exec.Query(ctx, query, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	resources := []model.Resource{}
	for rows.Next() {
		var resource model.Resource

		err := rows.Scan(
			&resource.ResourceID,
			&resource.Type,
			&resource.Reference,
			&resource.IsArchived,
			&resource.ServiceOwnershipTeamID,
			&resource.ServiceOwnershipTeamName,
			&resource.LastServicedAt,
		)
		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (r *ResourceRepository) Count(
	ctx context.Context,
	exec db.PGExecutor,
	q model.GetResourcesQuery,
) (int, error) {

	whereClause := `
WHERE
	is_archived = FALSE
`
	if q.IsArchived {
		whereClause = `
WHERE
	is_archived = TRUE
`
	}

	query := `
SELECT
	COUNT(*)
FROM
	resource_view
` + whereClause

	var count int
	err := exec.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ResourceRepository) GetResourceServiceMetrics(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
) ([]model.ServiceMetric, error) {

	query := `
SELECT
    m.resource_service_metric_id,
    m.name,
    m.description,
    m.is_cumulative,
	m.is_archived
FROM
    resource_service_metric AS m
JOIN
    resource_service_schedule AS ss
    ON m.resource_service_metric_id = ss.resource_service_metric_id
WHERE
    ss.resource_id = $1
	AND ss.is_archived = FALSE
	AND m.is_archived = FALSE;
`

	rows, err := exec.Query(ctx, query, resourceID)
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return serviceMetrics, nil
}

func (r *ResourceRepository) ListResourceMetricSchedules(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
	includeArchived bool,
) ([]model.ResourceServiceMetricStatus, error) {

	query := `
SELECT
    resource_service_schedule_id,
    resource_id,
    resource_service_metric_id,
    metric_name,
    current_value,
	threshold,
	normalised_value,
	is_due,
	last_recorded_at,
	schedule_is_archived,
	metric_is_archived
FROM
    resource_service_metric_status_view
WHERE
	resource_id = $1
`
	if !includeArchived {
		query += `
	AND schedule_is_archived = FALSE
	AND metric_is_archived = FALSE
`
	}

	query += `
ORDER BY metric_name ASC
`

	rows, err := exec.Query(ctx, query, resourceID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	metrics := []model.ResourceServiceMetricStatus{}
	for rows.Next() {
		var metric model.ResourceServiceMetricStatus

		err := rows.Scan(
			&metric.ResourceServiceScheduleID,
			&metric.ResourceID,
			&metric.ResourceServiceMetricID,
			&metric.MetricName,
			&metric.CurrentValue,
			&metric.Threshold,
			&metric.NormalisedValue,
			&metric.IsDue,
			&metric.LastRecordedAt,
			&metric.ScheduleIsArchived,
			&metric.MetricIsArchived,
		)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}
