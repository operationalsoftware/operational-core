package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"fmt"
	"log"
	"strings"

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

func (r *ResourceRepository) CreateResourceMetricRecord(
	ctx context.Context,
	exec db.PGExecutor,
	record model.NewResourceServiceMetricRecord,
) error {

	query := `
INSERT INTO resource_metric_recording (
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

func (r *ResourceRepository) CloseOpenMetricRecords(
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
	resource_metric_recording

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

func (r *ResourceRepository) ReopenClosedMetricRecords(
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
	resource_metric_recording
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

	whereClause, args := generateResourceWhereClause(q)

	limitPlaceholder := fmt.Sprintf("$%d", len(args)+1)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args)+2)

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	orderByClause, _ := q.Sort.ToOrderByClause(model.Resource{})

	if orderByClause == "" {
		orderByClause = "ORDER BY resource_id ASC"
	}

	query := `
SELECT
    r.resource_id,
    r.type,
    r.reference,
    r.is_archived,
	r.service_ownership_team_id,
	r.service_ownership_team_name,
    r.last_serviced_at,
	COALESCE((
		SELECT ARRAY_AGG(ss.name ORDER BY ss.name)
		FROM service_schedule_assignment ssa
		JOIN service_schedule ss ON ss.service_schedule_id = ssa.service_schedule_id
		JOIN resource_service_metric m ON m.resource_service_metric_id = ss.resource_service_metric_id
		WHERE ssa.resource_id = r.resource_id
			AND ss.is_archived = FALSE
			AND m.is_archived = FALSE
	), ARRAY[]::text[]) AS service_schedule_names
FROM
    resource_view r
` + whereClause + `
` + orderByClause + `
` + fmt.Sprintf("LIMIT %s OFFSET %s", limitPlaceholder, offsetPlaceholder) + `
`

	rows, err := exec.Query(ctx, query, append(args, limit, offset)...)
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
			&resource.ServiceScheduleNames,
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

	whereClause, args := generateResourceWhereClause(q)

	query := `
SELECT
	COUNT(*)
FROM
	resource_view
` + whereClause

	var count int
	err := exec.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ResourceRepository) GetAvailableFilters(
	ctx context.Context,
	exec db.PGExecutor,
	baseFilters model.GetResourcesQuery,
) (model.ResourceAvailableFilters, error) {

	mapping := map[string]string{
		"TypeIn":                 "type",
		"ServiceOwnershipTeamIn": "COALESCE(service_ownership_team_name, 'Unassigned')",
		"ReferenceIn":            "reference",
	}

	avail := model.ResourceAvailableFilters{}

	var collect = func(key string, dest *[]string) error {
		queryFilters := baseFilters

		switch key {
		case "TypeIn":
			queryFilters.TypeIn = nil
		case "ServiceOwnershipTeamIn":
			queryFilters.ServiceOwnershipTeamIn = nil
		case "ReferenceIn":
			queryFilters.ReferenceIn = nil
		}

		where, args := generateResourceWhereClause(queryFilters)
		col := mapping[key]

		if col != "" {
			if where == "" {
				where = "WHERE " + col + " IS NOT NULL"
			} else {
				where += "\nAND " + col + " IS NOT NULL"
			}
		}

		query := `
SELECT DISTINCT ` + col + ` AS val
FROM resource_view
` + where + `
ORDER BY val ASC
`

		rows, err := exec.Query(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var v string
			if err := rows.Scan(&v); err != nil {
				return err
			}
			*dest = append(*dest, v)
		}

		return rows.Err()
	}

	if err := collect("TypeIn", &avail.TypeIn); err != nil {
		return avail, err
	}
	if err := collect("ServiceOwnershipTeamIn", &avail.ServiceOwnershipTeamIn); err != nil {
		return avail, err
	}
	if err := collect("ReferenceIn", &avail.ReferenceIn); err != nil {
		return avail, err
	}

	return avail, nil
}

func generateResourceWhereClause(q model.GetResourcesQuery) (string, []any) {
	var whereClauses []string
	var args []any
	argID := 1

	if !q.IsArchived {
		whereClauses = append(whereClauses, "is_archived = FALSE")
	}

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

	addInClause("type", q.TypeIn)
	addInClause("COALESCE(service_ownership_team_name, 'Unassigned')", q.ServiceOwnershipTeamIn)
	addInClause("reference", q.ReferenceIn)

	if len(whereClauses) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(whereClauses, "\nAND "), args
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
JOIN service_schedule AS ss
	ON m.resource_service_metric_id = ss.resource_service_metric_id
JOIN service_schedule_assignment AS ssa
	ON ssa.service_schedule_id = ss.service_schedule_id
WHERE
    ssa.resource_id = $1
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

func (r *ResourceRepository) ListMetricLifetimeTotals(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
) ([]model.ServiceMetricLifetimeTotal, error) {

	query := `
SELECT
    resource_id,
    resource_type,
    reference,
    metric_name,
    lifetime_total
FROM
    service_metric_lifetime_total_view
WHERE
    resource_id = $1
ORDER BY metric_name ASC
`

	rows, err := exec.Query(ctx, query, resourceID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var totals []model.ServiceMetricLifetimeTotal

	for rows.Next() {
		var total model.ServiceMetricLifetimeTotal
		err := rows.Scan(
			&total.ResourceID,
			&total.ResourceType,
			&total.ResourceReference,
			&total.MetricName,
			&total.LifetimeTotal,
		)
		if err != nil {
			return nil, err
		}
		totals = append(totals, total)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return totals, nil
}

func (r *ResourceRepository) ListResourceMetricSchedules(
	ctx context.Context,
	exec db.PGExecutor,
	resourceID int,
) ([]model.ResourceServiceMetricStatus, error) {

	query := `
SELECT
  type,
  reference,
  service_ownership_team_id,
  service_ownership_team_name,
	service_schedule_id,
  service_schedule_name,
  resource_id,
  resource_service_metric_id,
  metric_name,
  current_value,
	threshold,
	normalised_value,
	normalised_percentage,
	is_due,
	last_recorded_at,
	schedule_is_archived,
	metric_is_archived,
	last_serviced_at,
	wip_service_id,
	has_wip_service
FROM
    resource_service_metric_status_view
WHERE
	resource_id = $1
`
	query += `
	AND schedule_is_archived = FALSE
	AND metric_is_archived = FALSE
`

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
			&metric.Type,
			&metric.Reference,
			&metric.ServiceOwnershipTeamID,
			&metric.ServiceOwnershipTeamName,
			&metric.ServiceScheduleID,
			&metric.ServiceScheduleName,
			&metric.ResourceID,
			&metric.ResourceServiceMetricID,
			&metric.MetricName,
			&metric.CurrentValue,
			&metric.Threshold,
			&metric.NormalisedValue,
			&metric.NormalisedPercentage,
			&metric.IsDue,
			&metric.LastRecordedAt,
			&metric.ScheduleIsArchived,
			&metric.MetricIsArchived,
			&metric.LastServicedAt,
			&metric.WIPServiceID,
			&metric.HasWIPService,
		)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}
