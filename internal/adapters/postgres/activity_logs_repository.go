package postgres

// import (
// 	"context"
// 	"database/sql"

// 	"drones/internal/core/domain"
// 	"drones/internal/ports"
// )

// type ActivityLogsRepository struct {
// 	db                *sql.DB
// 	logger            ports.Logger
// 	createStmt        *sql.Stmt
// 	getByIDStmt       *sql.Stmt
// 	updateStmt        *sql.Stmt
// 	deleteStmt        *sql.Stmt
// 	listStmt          *sql.Stmt
// 	getByUserIDStmt   *sql.Stmt
// 	getByActionStmt   *sql.Stmt
// 	getByResourceStmt *sql.Stmt
// 	countLogsStmt     *sql.Stmt
// }

// func NewActivityLogsRepository(db *sql.DB, logger ports.Logger) ports.ActivityLogsRepository {
// 	repo := &ActivityLogsRepository{
// 		db:     db,
// 		logger: logger,
// 	}
// 	// Prepare statements
// 	if err := repo.prepareStatements(); err != nil {
// 		logger.Error("Activity logs: Failed to prepare statements", "error", err)
// 		return repo // Return anyway, fallback to non-prepared queries
// 	}

// 	return repo
// }

// func (r *ActivityLogsRepository) prepareStatements() error {
// 	var err error

// 	// Create statement
// 	r.createStmt, err = r.db.Prepare(`
// 		INSERT INTO activity_logs (
// 			actor_id, actor_type, action, performed_at, ip, device, location, 
// 			resource_name, resource_id, created_by_id
// 		) VALUES (
// 			$1, $2, $3, NOW(), $4, $5, $6, $7, $8, $9
// 		) RETURNING 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by ID statement
// 	r.getByIDStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM activity_logs 
// 		WHERE id = $1 AND active = TRUE`)
// 	if err != nil {
// 		return err
// 	}

// 	// Update statement
// 	r.updateStmt, err = r.db.Prepare(`
// 		UPDATE activity_logs SET 
// 			actor_type = COALESCE($2, actor_type),
// 			action = COALESCE($3, action),
// 			performed_at = COALESCE($4, performed_at),
// 			ip = COALESCE($5, ip),
// 			device = COALESCE($6, device),
// 			location = COALESCE($7, location),
// 			resource_name = COALESCE($8, resource_name),
// 			resource_id = COALESCE($9, resource_id),
// 			updated_by_id = COALESCE($10, updated_by_id),
// 			updated_at = NOW()
// 		WHERE id = $1 AND active = TRUE
// 		RETURNING 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id`)
// 	if err != nil {
// 		return err
// 	}

// 	// Soft delete statement
// 	r.deleteStmt, err = r.db.Prepare(`
// 		UPDATE activity_logs SET 
// 			active = FALSE, 
// 			deleted_at = NOW(), 
// 			deleted_by_id = $2 
// 		WHERE id = $1 AND active = TRUE
// 		RETURNING 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id`)
// 	if err != nil {
// 		return err
// 	}

// 	// List statement with pagination and filtering
// 	r.listStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM activity_logs 
// 		WHERE active = TRUE 
// 		AND ($1::VARCHAR IS NULL OR actor_id = $1)
// 		AND ($2::VARCHAR IS NULL OR action = $2)
// 		AND ($3::TIMESTAMP IS NULL OR performed_at >= $3)
// 		AND ($4::TIMESTAMP IS NULL OR performed_at <= $4)
// 		ORDER BY performed_at DESC 
// 		LIMIT $5 OFFSET $6`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by user ID
// 	r.getByUserIDStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM activity_logs 
// 		WHERE actor_id = $1 AND active = TRUE
// 		ORDER BY performed_at DESC`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by action
// 	r.getByActionStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM activity_logs 
// 		WHERE active = TRUE AND action = $1
// 		ORDER BY performed_at DESC`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by resource
// 	r.getByResourceStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, actor_id, actor_type, action, performed_at, ip, device, location,
// 			resource_name, resource_id, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM activity_logs 
// 		WHERE active = TRUE 
// 		AND resource_name = $1 
// 		AND ($2::UUID IS NULL OR resource_id = $2::UUID)
// 		ORDER BY performed_at DESC`)
// 	if err != nil {
// 		return err
// 	}

// 	// Count logs
// 	r.countLogsStmt, err = r.db.Prepare(`
// 		SELECT COUNT(*) 
// 		FROM activity_logs 
// 		WHERE active = TRUE 
// 		AND ($1::VARCHAR IS NULL OR actor_id = $1)
// 		AND ($2::VARCHAR IS NULL OR action = $2)
// 		AND ($3::TIMESTAMP IS NULL OR performed_at >= $3)
// 		AND ($4::TIMESTAMP IS NULL OR performed_at <= $4)`)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r *ActivityLogsRepository) Close() error {
// 	if r.createStmt != nil {
// 		if err := r.createStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close createStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.getByIDStmt != nil {
// 		if err := r.getByIDStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close getByIDStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.updateStmt != nil {
// 		if err := r.updateStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close updateStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.deleteStmt != nil {
// 		if err := r.deleteStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close deleteStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.listStmt != nil {
// 		if err := r.listStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close listStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.getByUserIDStmt != nil {
// 		if err := r.getByUserIDStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close getByUserIDStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.getByActionStmt != nil {
// 		if err := r.getByActionStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close getByActionStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.getByResourceStmt != nil {
// 		if err := r.getByResourceStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close getByResourceStmt", "error", err)
// 			return err
// 		}
// 	}
// 	if r.countLogsStmt != nil {
// 		if err := r.countLogsStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close countLogsStmt", "error", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (r *ActivityLogsRepository) GetDB() *sql.DB {
// 	return r.db
// }

// // CreateActivityLog creates a new activity log in the repository
// func (r *ActivityLogsRepository) CreateActivityLog(ctx context.Context, createLog *domain.CreateActivityLogRequest) (*domain.ActivityLog, error) {
// 	var log domain.ActivityLog

// 	// Set default values based on action
// 	resourceName := ""
// 	switch createLog.Action {
// 	case "login":
// 		resourceName = "authentication"
// 	case "logout":
// 		resourceName = "authentication"
// 	case "create_order":
// 		resourceName = "order"
// 	case "update_order":
// 		resourceName = "order"
// 	case "delete_order":
// 		resourceName = "order"
// 	case "create_drone":
// 		resourceName = "drone"
// 	case "update_drone":
// 		resourceName = "drone"
// 	case "delete_drone":
// 		resourceName = "drone"
// 	default:
// 		resourceName = "system"
// 	}

// 	err := r.createStmt.QueryRowContext(ctx,
// 		createLog.UserID,
// 		"user", // default actor_type
// 		createLog.Action,
// 		createLog.Metadata.IP,
// 		createLog.Metadata.Device,
// 		createLog.Metadata.Location,
// 		resourceName,
// 		nil,              // resource_id - to be set by caller if needed
// 		createLog.UserID, // created_by_id
// 	).Scan(
// 		&log.ID,
// 		&log.ActorId,
// 		&log.ActorType,
// 		&log.Action,
// 		&log.PerformedAt,
// 		&log.IP,
// 		&log.Device,
// 		&log.Location,
// 		&log.ResourceName,
// 		&log.ResourceID,
// 		&log.CreatedAt,
// 		&log.UpdatedAt,
// 		&log.DeletedAt,
// 		&log.Active,
// 		&log.CreatedByID,
// 		&log.UpdatedByID,
// 		&log.DeletedByID,
// 	)

// 	if err != nil {
// 		r.logger.Error("Failed to create activity log", "userID", createLog.UserID, "action", createLog.Action, "error", err)
// 		return nil, err
// 	}

// 	return &log, nil
// }

// // GetByID retrieves an activity log by its ID
// func (r *ActivityLogsRepository) GetByID(ctx context.Context, logID string) (*domain.ActivityLog, error) {
// 	var log domain.ActivityLog

// 	err := r.getByIDStmt.QueryRowContext(ctx, logID).Scan(
// 		&log.ID,
// 		&log.ActorId,
// 		&log.ActorType,
// 		&log.Action,
// 		&log.PerformedAt,
// 		&log.IP,
// 		&log.Device,
// 		&log.Location,
// 		&log.ResourceName,
// 		&log.ResourceID,
// 		&log.CreatedAt,
// 		&log.UpdatedAt,
// 		&log.DeletedAt,
// 		&log.Active,
// 		&log.CreatedByID,
// 		&log.UpdatedByID,
// 		&log.DeletedByID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			r.logger.Warn("Activity log not found", "logID", logID)
// 			return nil, domain.NewDomainError(domain.ActivityLogNotFoundError, "Activity log not found", nil)
// 		}
// 		r.logger.Error("Failed to get activity log by ID", "logID", logID, "error", err)
// 		return nil, err
// 	}

// 	return &log, nil
// }

// // UpdateActivityLog updates an existing activity log in the repository
// func (r *ActivityLogsRepository) UpdateActivityLog(ctx context.Context, logID string, update *domain.UpdateActivityLogRequest) (*domain.ActivityLog, error) {
// 	var log domain.ActivityLog

// 	err := r.updateStmt.QueryRowContext(ctx,
// 		logID,
// 		update.ActorType,
// 		update.Action,
// 		update.PerformedAt,
// 		update.Metadata.IP,
// 		update.Metadata.Device,
// 		update.Metadata.Location,
// 		update.ResourceName,
// 		update.ResourceID,
// 		update.UpdatedByID,
// 	).Scan(
// 		&log.ID,
// 		&log.ActorId,
// 		&log.ActorType,
// 		&log.Action,
// 		&log.PerformedAt,
// 		&log.IP,
// 		&log.Device,
// 		&log.Location,
// 		&log.ResourceName,
// 		&log.ResourceID,
// 		&log.CreatedAt,
// 		&log.UpdatedAt,
// 		&log.DeletedAt,
// 		&log.Active,
// 		&log.CreatedByID,
// 		&log.UpdatedByID,
// 		&log.DeletedByID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			r.logger.Warn("Activity log not found for update", "logID", logID)
// 			return nil, domain.NewDomainError(domain.ActivityLogNotFoundError, "Activity log not found", nil)
// 		}
// 		r.logger.Error("Failed to update activity log", "logID", logID, "error", err)
// 		return nil, err
// 	}

// 	return &log, nil
// }

// // DeleteActivityLog soft deletes an activity log from the repository
// func (r *ActivityLogsRepository) DeleteActivityLog(ctx context.Context, logID string) error {
// 	var log domain.ActivityLog

// 	// Get the user ID from context or use system user ID for deletedByID
// 	var deletedByID *string // This should come from authenticated user context

// 	err := r.deleteStmt.QueryRowContext(ctx, logID, deletedByID).Scan(
// 		&log.ID,
// 		&log.ActorId,
// 		&log.ActorType,
// 		&log.Action,
// 		&log.PerformedAt,
// 		&log.IP,
// 		&log.Device,
// 		&log.Location,
// 		&log.ResourceName,
// 		&log.ResourceID,
// 		&log.CreatedAt,
// 		&log.UpdatedAt,
// 		&log.DeletedAt,
// 		&log.Active,
// 		&log.CreatedByID,
// 		&log.UpdatedByID,
// 		&log.DeletedByID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			r.logger.Warn("Activity log not found for deletion", "logID", logID)
// 			return domain.NewDomainError(domain.ActivityLogNotFoundError, "Activity log not found", nil)
// 		}
// 		r.logger.Error("Failed to delete activity log", "logID", logID, "error", err)
// 		return err
// 	}

// 	return nil
// }

// // ListActivityLogs retrieves activity logs with filtering and pagination
// func (r *ActivityLogsRepository) ListActivityLogs(ctx context.Context, options domain.PaginationOption[domain.ActivityLogFilter]) (*domain.Pagination[domain.ActivityLog], error) {
// 	filter := options.Filter
// 	offset := options.Offset
// 	limit := options.Limit

// 	// Prepare filter parameters for prepared statements
// 	var userIDParam interface{}
// 	var actionParam interface{}
// 	var startTimeParam interface{}
// 	var endTimeParam interface{}

// 	if filter != nil {
// 		if filter.UserID != nil {
// 			userIDParam = *filter.UserID
// 		}
// 		if filter.Action != nil {
// 			actionParam = *filter.Action
// 		}
// 		if filter.StartTime != nil {
// 			startTimeParam = *filter.StartTime
// 		}
// 		if filter.EndTime != nil {
// 			endTimeParam = *filter.EndTime
// 		}
// 	}

// 	// Count total records using prepared statement
// 	var total int
// 	err := r.countLogsStmt.QueryRowContext(ctx, userIDParam, actionParam, startTimeParam, endTimeParam).Scan(&total)
// 	if err != nil {
// 		r.logger.Error("Failed to count activity logs", "filter", filter, "error", err)
// 		return nil, err
// 	}

// 	// Get paginated results using prepared statement
// 	rows, err := r.listStmt.QueryContext(ctx, userIDParam, actionParam, startTimeParam, endTimeParam, limit, offset)
// 	if err != nil {
// 		r.logger.Error("Failed to list activity logs", "filter", filter, "limit", limit, "offset", offset, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var logs []*domain.ActivityLog
// 	for rows.Next() {
// 		var log domain.ActivityLog
// 		err := rows.Scan(
// 			&log.ID,
// 			&log.ActorId,
// 			&log.ActorType,
// 			&log.Action,
// 			&log.PerformedAt,
// 			&log.IP,
// 			&log.Device,
// 			&log.Location,
// 			&log.ResourceName,
// 			&log.ResourceID,
// 			&log.CreatedAt,
// 			&log.UpdatedAt,
// 			&log.DeletedAt,
// 			&log.Active,
// 			&log.CreatedByID,
// 			&log.UpdatedByID,
// 			&log.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan activity log row", "error", err)
// 			continue
// 		}

// 		logs = append(logs, &log)
// 	}

// 	if err = rows.Err(); err != nil {
// 		r.logger.Error("Failed to iterate activity log rows", "error", err)
// 		return nil, err
// 	}

// 	page := offset/limit + 1
// 	pageSize := limit
// 	totalPages := (total + limit - 1) / limit

// 	result := &domain.Pagination[domain.ActivityLog]{
// 		Data:       logs,
// 		Total:      total,
// 		Page:       page,
// 		PageSize:   pageSize,
// 		TotalPages: totalPages,
// 	}

// 	return result, nil
// }

// // Additional specialized methods

// // GetByUserID retrieves all activity logs for a specific user
// func (r *ActivityLogsRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.ActivityLog, error) {
// 	rows, err := r.getByUserIDStmt.QueryContext(ctx, userID)
// 	if err != nil {
// 		r.logger.Error("Failed to get activity logs by user ID", "userID", userID, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var logs []*domain.ActivityLog
// 	for rows.Next() {
// 		var log domain.ActivityLog
// 		err := rows.Scan(
// 			&log.ID,
// 			&log.ActorId,
// 			&log.ActorType,
// 			&log.Action,
// 			&log.PerformedAt,
// 			&log.IP,
// 			&log.Device,
// 			&log.Location,
// 			&log.ResourceName,
// 			&log.ResourceID,
// 			&log.CreatedAt,
// 			&log.UpdatedAt,
// 			&log.DeletedAt,
// 			&log.Active,
// 			&log.CreatedByID,
// 			&log.UpdatedByID,
// 			&log.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan activity log by user ID row", "error", err)
// 			return nil, err
// 		}
// 		logs = append(logs, &log)
// 	}

// 	return logs, rows.Err()
// }

// // GetByAction retrieves all activity logs with a specific action
// func (r *ActivityLogsRepository) GetByAction(ctx context.Context, action string) ([]*domain.ActivityLog, error) {
// 	rows, err := r.getByActionStmt.QueryContext(ctx, action)
// 	if err != nil {
// 		r.logger.Error("Failed to get activity logs by action", "action", action, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var logs []*domain.ActivityLog
// 	for rows.Next() {
// 		var log domain.ActivityLog
// 		err := rows.Scan(
// 			&log.ID,
// 			&log.ActorId,
// 			&log.ActorType,
// 			&log.Action,
// 			&log.PerformedAt,
// 			&log.IP,
// 			&log.Device,
// 			&log.Location,
// 			&log.ResourceName,
// 			&log.ResourceID,
// 			&log.CreatedAt,
// 			&log.UpdatedAt,
// 			&log.DeletedAt,
// 			&log.Active,
// 			&log.CreatedByID,
// 			&log.UpdatedByID,
// 			&log.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan activity log by action row", "error", err)
// 			return nil, err
// 		}
// 		logs = append(logs, &log)
// 	}

// 	return logs, rows.Err()
// }

// // GetByResource retrieves all activity logs for a specific resource
// func (r *ActivityLogsRepository) GetByResource(ctx context.Context, resourceName string, resourceID *string) ([]*domain.ActivityLog, error) {
// 	var resourceIDParam interface{}
// 	if resourceID != nil {
// 		resourceIDParam = *resourceID
// 	}

// 	rows, err := r.getByResourceStmt.QueryContext(ctx, resourceName, resourceIDParam)
// 	if err != nil {
// 		r.logger.Error("Failed to get activity logs by resource", "resourceName", resourceName, "resourceID", resourceID, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var logs []*domain.ActivityLog
// 	for rows.Next() {
// 		var log domain.ActivityLog
// 		err := rows.Scan(
// 			&log.ID,
// 			&log.ActorId,
// 			&log.ActorType,
// 			&log.Action,
// 			&log.PerformedAt,
// 			&log.IP,
// 			&log.Device,
// 			&log.Location,
// 			&log.ResourceName,
// 			&log.ResourceID,
// 			&log.CreatedAt,
// 			&log.UpdatedAt,
// 			&log.DeletedAt,
// 			&log.Active,
// 			&log.CreatedByID,
// 			&log.UpdatedByID,
// 			&log.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan activity log by resource row", "error", err)
// 			return nil, err
// 		}
// 		logs = append(logs, &log)
// 	}

// 	return logs, rows.Err()
// }

// // CountActivityLogs counts the total number of activity logs with optional filters
// func (r *ActivityLogsRepository) CountActivityLogs(ctx context.Context, filter *domain.ActivityLogFilter) (int64, error) {
// 	var count int64

// 	// Prepare filter parameters
// 	var userIDParam interface{}
// 	var actionParam interface{}
// 	var startTimeParam interface{}
// 	var endTimeParam interface{}

// 	if filter != nil {
// 		if filter.UserID != nil {
// 			userIDParam = *filter.UserID
// 		}
// 		if filter.Action != nil {
// 			actionParam = *filter.Action
// 		}
// 		if filter.StartTime != nil {
// 			startTimeParam = *filter.StartTime
// 		}
// 		if filter.EndTime != nil {
// 			endTimeParam = *filter.EndTime
// 		}
// 	}

// 	err := r.countLogsStmt.QueryRowContext(ctx, userIDParam, actionParam, startTimeParam, endTimeParam).Scan(&count)
// 	if err != nil {
// 		r.logger.Error("Failed to count activity logs", "filter", filter, "error", err)
// 		return 0, err
// 	}

// 	return count, nil
// }
