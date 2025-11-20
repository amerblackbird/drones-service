package postgres

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"

// 	"drones/internal/core/domain"
// 	"drones/internal/ports"
// )

// type AuditLogsRepository struct {
// 	db                   *sql.DB
// 	logger               ports.Logger
// 	createStmt           *sql.Stmt
// 	getByIDStmt          *sql.Stmt
// 	updateStmt           *sql.Stmt
// 	deleteStmt           *sql.Stmt
// 	listStmt             *sql.Stmt
// 	getByPerformedByStmt *sql.Stmt
// 	getByActionStmt      *sql.Stmt
// 	getByResourceStmt    *sql.Stmt
// 	countAuditLogsStmt   *sql.Stmt
// }

// func NewAuditLogsRepository(db *sql.DB, logger ports.Logger) ports.AuditLogsRepository {
// 	repo := &AuditLogsRepository{
// 		db:     db,
// 		logger: logger,
// 	}
// 	// Prepare statements
// 	if err := repo.prepareStatements(); err != nil {
// 		logger.Error("Audit: Failed to prepare statements", "error", err)
// 		return repo // Return anyway, fallback to non-prepared queries
// 	}

// 	return repo
// }

// func (r *AuditLogsRepository) Close() error {
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
// 	if r.getByPerformedByStmt != nil {
// 		if err := r.getByPerformedByStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close getByPerformedByStmt", "error", err)
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
// 	if r.countAuditLogsStmt != nil {
// 		if err := r.countAuditLogsStmt.Close(); err != nil {
// 			r.logger.Error("Failed to close countAuditLogsStmt", "error", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (r *AuditLogsRepository) GetDB() *sql.DB {
// 	return r.db
// }

// func (r *AuditLogsRepository) prepareStatements() error {
// 	var err error

// 	// Create statement with comprehensive fields
// 	r.createStmt, err = r.db.Prepare(`
// 		INSERT INTO audit_logs (
// 			performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_by_id
// 		) VALUES (
// 			$1, $2, $3, $4, $5, $6, $7, $8, $9
// 		) RETURNING 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by ID statement with all fields
// 	r.getByIDStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM audit_logs 
// 		WHERE id = $1 AND active = TRUE`)
// 	if err != nil {
// 		return err
// 	}

// 	// Update statement with comprehensive fields
// 	r.updateStmt, err = r.db.Prepare(`
// 		UPDATE audit_logs SET 
// 			performed_user_type = COALESCE($2, performed_user_type),
// 			action = COALESCE($3, action),
// 			resource_name = COALESCE($4, resource_name),
// 			resource_id = COALESCE($5, resource_id),
// 			performed_at = COALESCE($6, performed_at),
// 			old_data = COALESCE($7, old_data),
// 			new_data = COALESCE($8, new_data),
// 			updated_by_id = COALESCE($9, updated_by_id),
// 			updated_at = NOW()
// 		WHERE id = $1 AND active = TRUE
// 		RETURNING 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id`)
// 	if err != nil {
// 		return err
// 	}

// 	// Soft delete statement
// 	r.deleteStmt, err = r.db.Prepare(`
// 		UPDATE audit_logs SET 
// 			active = FALSE, 
// 			deleted_at = NOW(), 
// 			deleted_by_id = $2 
// 		WHERE id = $1 AND active = TRUE
// 		RETURNING 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id`)
// 	if err != nil {
// 		return err
// 	}

// 	// List statement with pagination and filtering
// 	r.listStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM audit_logs 
// 		WHERE active = TRUE 
// 		AND ($1::UUID IS NULL OR performed_by_id = $1)
// 		AND ($2::VARCHAR IS NULL OR performed_user_type = $2)
// 		AND ($3::VARCHAR IS NULL OR action = $3)
// 		AND ($4::VARCHAR IS NULL OR resource_name = $4)
// 		AND ($5::VARCHAR IS NULL OR resource_id = $5)
// 		ORDER BY performed_at DESC 
// 		LIMIT $6 OFFSET $7`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by performed by
// 	r.getByPerformedByStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM audit_logs 
// 		WHERE performed_by_id = $1 AND active = TRUE
// 		ORDER BY performed_at DESC`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by action
// 	r.getByActionStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM audit_logs 
// 		WHERE active = TRUE AND action = $1
// 		ORDER BY performed_at DESC`)
// 	if err != nil {
// 		return err
// 	}

// 	// Get by resource
// 	r.getByResourceStmt, err = r.db.Prepare(`
// 		SELECT 
// 			id, performed_by_id, performed_user_type, action, resource_name, resource_id, 
// 			performed_at, old_data, new_data, created_at, updated_at, deleted_at, 
// 			active, created_by_id, updated_by_id, deleted_by_id
// 		FROM audit_logs 
// 		WHERE active = TRUE 
// 		AND resource_name = $1 
// 		AND ($2::VARCHAR IS NULL OR resource_id = $2)
// 		ORDER BY performed_at DESC`)
// 	if err != nil {
// 		return err
// 	}

// 	// Count audit logs
// 	r.countAuditLogsStmt, err = r.db.Prepare(`
// 		SELECT COUNT(*) 
// 		FROM audit_logs 
// 		WHERE active = TRUE 
// 		AND ($1::UUID IS NULL OR performed_by_id = $1)
// 		AND ($2::VARCHAR IS NULL OR performed_user_type = $2)
// 		AND ($3::VARCHAR IS NULL OR action = $3)
// 		AND ($4::VARCHAR IS NULL OR resource_name = $4)
// 		AND ($5::VARCHAR IS NULL OR resource_id = $5)`)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // CreateAuditLog creates a new audit log in the repository
// func (r *AuditLogsRepository) CreateAuditLog(ctx context.Context, createAuditLog *domain.CreateAuditLogRequest) (*domain.AuditLog, error) {
// 	var auditLog domain.AuditLog

// 	// Convert old_data and new_data to JSON
// 	var oldDataJSON []byte
// 	var newDataJSON []byte
// 	var err error

// 	if createAuditLog.OldData != nil {
// 		oldDataJSON, err = json.Marshal(createAuditLog.OldData)
// 		if err != nil {
// 			r.logger.Error("Failed to marshal old_data to JSON", "error", err)
// 			return nil, err
// 		}
// 	}

// 	if createAuditLog.NewData != nil {
// 		newDataJSON, err = json.Marshal(createAuditLog.NewData)
// 		if err != nil {
// 			r.logger.Error("Failed to marshal new_data to JSON", "error", err)
// 			return nil, err
// 		}
// 	}

// 	err = r.createStmt.QueryRowContext(ctx,
// 		createAuditLog.PerformedBy,
// 		createAuditLog.UserType,
// 		createAuditLog.Action,
// 		createAuditLog.ResourceName,
// 		createAuditLog.ResourceID,
// 		createAuditLog.PerformedAt,
// 		oldDataJSON,
// 		newDataJSON,
// 		createAuditLog.PerformedBy, // created_by_id
// 	).Scan(
// 		&auditLog.ID,
// 		&auditLog.PerformedById,
// 		&auditLog.PerformedUserType,
// 		&auditLog.Action,
// 		&auditLog.ResourceName,
// 		&auditLog.ResourceID,
// 		&auditLog.PerformedAt,
// 		&auditLog.OldData,
// 		&auditLog.NewData,
// 		&auditLog.CreatedAt,
// 		&auditLog.UpdatedAt,
// 		&auditLog.DeletedAt,
// 		&auditLog.Active,
// 		&auditLog.CreatedByID,
// 		&auditLog.UpdatedByID,
// 		&auditLog.DeletedByID,
// 	)

// 	if err != nil {
// 		r.logger.Error("Failed to create audit log", "performedBy", createAuditLog.PerformedBy, "action", createAuditLog.Action, "error", err)
// 		return nil, err
// 	}

// 	return &auditLog, nil
// }

// // GetByID retrieves an audit log by its ID
// func (r *AuditLogsRepository) GetByID(ctx context.Context, auditLogID string) (*domain.AuditLog, error) {
// 	var auditLog domain.AuditLog

// 	err := r.getByIDStmt.QueryRowContext(ctx, auditLogID).Scan(
// 		&auditLog.ID,
// 		&auditLog.PerformedById,
// 		&auditLog.PerformedUserType,
// 		&auditLog.Action,
// 		&auditLog.ResourceName,
// 		&auditLog.ResourceID,
// 		&auditLog.PerformedAt,
// 		&auditLog.OldData,
// 		&auditLog.NewData,
// 		&auditLog.CreatedAt,
// 		&auditLog.UpdatedAt,
// 		&auditLog.DeletedAt,
// 		&auditLog.Active,
// 		&auditLog.CreatedByID,
// 		&auditLog.UpdatedByID,
// 		&auditLog.DeletedByID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			r.logger.Warn("Audit log not found", "auditLogID", auditLogID)
// 			return nil, domain.NewDomainError(domain.AuditLogNotFoundError, "Audit log not found", nil)
// 		}
// 		r.logger.Error("Failed to get audit log by ID", "auditLogID", auditLogID, "error", err)
// 		return nil, err
// 	}

// 	return &auditLog, nil
// }

// // UpdateAuditLog updates an existing audit log in the repository
// func (r *AuditLogsRepository) UpdateAuditLog(ctx context.Context, auditLogID string, update *domain.UpdateAuditLogRequest) (*domain.AuditLog, error) {
// 	var auditLog domain.AuditLog

// 	// Convert old_data and new_data to JSON if provided
// 	var oldDataJSON interface{}
// 	var newDataJSON interface{}
// 	var err error

// 	if update.OldData != nil {
// 		oldDataJSON, err = json.Marshal(update.OldData)
// 		if err != nil {
// 			r.logger.Error("Failed to marshal old_data to JSON", "error", err)
// 			return nil, err
// 		}
// 	}

// 	if update.NewData != nil {
// 		newDataJSON, err = json.Marshal(update.NewData)
// 		if err != nil {
// 			r.logger.Error("Failed to marshal new_data to JSON", "error", err)
// 			return nil, err
// 		}
// 	}

// 	err = r.updateStmt.QueryRowContext(ctx,
// 		auditLogID,
// 		update.PerformedUserType,
// 		update.Action,
// 		update.ResourceName,
// 		update.ResourceID,
// 		update.PerformedAt,
// 		oldDataJSON,
// 		newDataJSON,
// 		update.UpdatedByID,
// 	).Scan(
// 		&auditLog.ID,
// 		&auditLog.PerformedById,
// 		&auditLog.PerformedUserType,
// 		&auditLog.Action,
// 		&auditLog.ResourceName,
// 		&auditLog.ResourceID,
// 		&auditLog.PerformedAt,
// 		&auditLog.OldData,
// 		&auditLog.NewData,
// 		&auditLog.CreatedAt,
// 		&auditLog.UpdatedAt,
// 		&auditLog.DeletedAt,
// 		&auditLog.Active,
// 		&auditLog.CreatedByID,
// 		&auditLog.UpdatedByID,
// 		&auditLog.DeletedByID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			r.logger.Warn("Audit log not found for update", "auditLogID", auditLogID)
// 			return nil, domain.NewDomainError(domain.AuditLogNotFoundError, "Audit log not found", nil)
// 		}
// 		r.logger.Error("Failed to update audit log", "auditLogID", auditLogID, "error", err)
// 		return nil, err
// 	}

// 	return &auditLog, nil
// }

// // DeleteAuditLog soft deletes an audit log from the repository
// func (r *AuditLogsRepository) DeleteAuditLog(ctx context.Context, auditLogID string) error {
// 	var auditLog domain.AuditLog

// 	// Get the user ID from context or use system user ID for deletedByID
// 	// For now, we'll use a default value
// 	var deletedByID *string // This should come from authenticated user context

// 	err := r.deleteStmt.QueryRowContext(ctx, auditLogID, deletedByID).Scan(
// 		&auditLog.ID,
// 		&auditLog.PerformedById,
// 		&auditLog.PerformedUserType,
// 		&auditLog.Action,
// 		&auditLog.ResourceName,
// 		&auditLog.ResourceID,
// 		&auditLog.PerformedAt,
// 		&auditLog.OldData,
// 		&auditLog.NewData,
// 		&auditLog.CreatedAt,
// 		&auditLog.UpdatedAt,
// 		&auditLog.DeletedAt,
// 		&auditLog.Active,
// 		&auditLog.CreatedByID,
// 		&auditLog.UpdatedByID,
// 		&auditLog.DeletedByID,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			r.logger.Warn("Audit log not found for deletion", "auditLogID", auditLogID)
// 			return domain.NewDomainError(domain.AuditLogNotFoundError, "Audit log not found", nil)
// 		}
// 		r.logger.Error("Failed to delete audit log", "auditLogID", auditLogID, "error", err)
// 		return err
// 	}

// 	return nil
// }

// // ListAuditLogs retrieves audit logs with filtering and pagination
// func (r *AuditLogsRepository) ListAuditLogs(ctx context.Context, options domain.PaginationOption[domain.AuditLogFilter]) (*domain.Pagination[domain.AuditLog], error) {
// 	filter := options.Filter
// 	offset := options.Offset
// 	limit := options.Limit

// 	// Prepare filter parameters for prepared statements
// 	var performedByParam interface{}
// 	var userTypeParam interface{}
// 	var actionParam interface{}
// 	var resourceNameParam interface{}
// 	var resourceIDParam interface{}

// 	if filter != nil {
// 		if filter.PerformedBy != nil {
// 			performedByParam = *filter.PerformedBy
// 		}
// 		if filter.UserType != nil {
// 			userTypeParam = *filter.UserType
// 		}
// 		if filter.Action != nil {
// 			actionParam = *filter.Action
// 		}
// 		if filter.ResourceName != nil {
// 			resourceNameParam = *filter.ResourceName
// 		}
// 		if filter.ResourceID != nil {
// 			resourceIDParam = *filter.ResourceID
// 		}
// 	}

// 	// Count total records using prepared statement
// 	var total int
// 	err := r.countAuditLogsStmt.QueryRowContext(ctx, performedByParam, userTypeParam, actionParam, resourceNameParam, resourceIDParam).Scan(&total)
// 	if err != nil {
// 		r.logger.Error("Failed to count audit logs", "filter", filter, "error", err)
// 		return nil, err
// 	}

// 	// Get paginated results using prepared statement
// 	rows, err := r.listStmt.QueryContext(ctx, performedByParam, userTypeParam, actionParam, resourceNameParam, resourceIDParam, limit, offset)
// 	if err != nil {
// 		r.logger.Error("Failed to list audit logs", "filter", filter, "limit", limit, "offset", offset, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var auditLogs []*domain.AuditLog
// 	for rows.Next() {
// 		var auditLog domain.AuditLog
// 		err := rows.Scan(
// 			&auditLog.ID,
// 			&auditLog.PerformedById,
// 			&auditLog.PerformedUserType,
// 			&auditLog.Action,
// 			&auditLog.ResourceName,
// 			&auditLog.ResourceID,
// 			&auditLog.PerformedAt,
// 			&auditLog.OldData,
// 			&auditLog.NewData,
// 			&auditLog.CreatedAt,
// 			&auditLog.UpdatedAt,
// 			&auditLog.DeletedAt,
// 			&auditLog.Active,
// 			&auditLog.CreatedByID,
// 			&auditLog.UpdatedByID,
// 			&auditLog.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan audit log row", "error", err)
// 			continue
// 		}

// 		auditLogs = append(auditLogs, &auditLog)
// 	}

// 	if err = rows.Err(); err != nil {
// 		r.logger.Error("Failed to iterate audit log rows", "error", err)
// 		return nil, err
// 	}

// 	page := offset/limit + 1
// 	pageSize := limit
// 	totalPages := (total + limit - 1) / limit

// 	result := &domain.Pagination[domain.AuditLog]{
// 		Data:       auditLogs,
// 		Total:      total,
// 		Page:       page,
// 		PageSize:   pageSize,
// 		TotalPages: totalPages,
// 	}

// 	return result, nil
// }

// // Additional specialized methods

// // GetByPerformedBy retrieves all audit logs performed by a specific user
// func (r *AuditLogsRepository) GetByPerformedBy(ctx context.Context, performedBy string) ([]*domain.AuditLog, error) {
// 	rows, err := r.getByPerformedByStmt.QueryContext(ctx, performedBy)
// 	if err != nil {
// 		r.logger.Error("Failed to get audit logs by performed by", "performedBy", performedBy, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var auditLogs []*domain.AuditLog
// 	for rows.Next() {
// 		var auditLog domain.AuditLog
// 		err := rows.Scan(
// 			&auditLog.ID,
// 			&auditLog.PerformedById,
// 			&auditLog.PerformedUserType,
// 			&auditLog.Action,
// 			&auditLog.ResourceName,
// 			&auditLog.ResourceID,
// 			&auditLog.PerformedAt,
// 			&auditLog.OldData,
// 			&auditLog.NewData,
// 			&auditLog.CreatedAt,
// 			&auditLog.UpdatedAt,
// 			&auditLog.DeletedAt,
// 			&auditLog.Active,
// 			&auditLog.CreatedByID,
// 			&auditLog.UpdatedByID,
// 			&auditLog.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan audit log by performed by row", "error", err)
// 			return nil, err
// 		}
// 		auditLogs = append(auditLogs, &auditLog)
// 	}

// 	return auditLogs, rows.Err()
// }

// // GetByAction retrieves all audit logs with a specific action
// func (r *AuditLogsRepository) GetByAction(ctx context.Context, action string) ([]*domain.AuditLog, error) {
// 	rows, err := r.getByActionStmt.QueryContext(ctx, action)
// 	if err != nil {
// 		r.logger.Error("Failed to get audit logs by action", "action", action, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var auditLogs []*domain.AuditLog
// 	for rows.Next() {
// 		var auditLog domain.AuditLog
// 		err := rows.Scan(
// 			&auditLog.ID,
// 			&auditLog.PerformedById,
// 			&auditLog.PerformedUserType,
// 			&auditLog.Action,
// 			&auditLog.ResourceName,
// 			&auditLog.ResourceID,
// 			&auditLog.PerformedAt,
// 			&auditLog.OldData,
// 			&auditLog.NewData,
// 			&auditLog.CreatedAt,
// 			&auditLog.UpdatedAt,
// 			&auditLog.DeletedAt,
// 			&auditLog.Active,
// 			&auditLog.CreatedByID,
// 			&auditLog.UpdatedByID,
// 			&auditLog.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan audit log by action row", "error", err)
// 			return nil, err
// 		}
// 		auditLogs = append(auditLogs, &auditLog)
// 	}

// 	return auditLogs, rows.Err()
// }

// // GetByResource retrieves all audit logs for a specific resource
// func (r *AuditLogsRepository) GetByResource(ctx context.Context, resourceName string, resourceID *string) ([]*domain.AuditLog, error) {
// 	var resourceIDParam interface{}
// 	if resourceID != nil {
// 		resourceIDParam = *resourceID
// 	}

// 	rows, err := r.getByResourceStmt.QueryContext(ctx, resourceName, resourceIDParam)
// 	if err != nil {
// 		r.logger.Error("Failed to get audit logs by resource", "resourceName", resourceName, "resourceID", resourceID, "error", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var auditLogs []*domain.AuditLog
// 	for rows.Next() {
// 		var auditLog domain.AuditLog
// 		err := rows.Scan(
// 			&auditLog.ID,
// 			&auditLog.PerformedById,
// 			&auditLog.PerformedUserType,
// 			&auditLog.Action,
// 			&auditLog.ResourceName,
// 			&auditLog.ResourceID,
// 			&auditLog.PerformedAt,
// 			&auditLog.OldData,
// 			&auditLog.NewData,
// 			&auditLog.CreatedAt,
// 			&auditLog.UpdatedAt,
// 			&auditLog.DeletedAt,
// 			&auditLog.Active,
// 			&auditLog.CreatedByID,
// 			&auditLog.UpdatedByID,
// 			&auditLog.DeletedByID,
// 		)
// 		if err != nil {
// 			r.logger.Error("Failed to scan audit log by resource row", "error", err)
// 			return nil, err
// 		}
// 		auditLogs = append(auditLogs, &auditLog)
// 	}

// 	return auditLogs, rows.Err()
// }

// // CountAuditLogs counts the total number of audit logs with optional filters
// func (r *AuditLogsRepository) CountAuditLogs(ctx context.Context, filter *domain.AuditLogFilter) (int64, error) {
// 	var count int64

// 	// Prepare filter parameters
// 	var performedByParam interface{}
// 	var userTypeParam interface{}
// 	var actionParam interface{}
// 	var resourceNameParam interface{}
// 	var resourceIDParam interface{}

// 	if filter != nil {
// 		if filter.PerformedBy != nil {
// 			performedByParam = *filter.PerformedBy
// 		}
// 		if filter.UserType != nil {
// 			userTypeParam = *filter.UserType
// 		}
// 		if filter.Action != nil {
// 			actionParam = *filter.Action
// 		}
// 		if filter.ResourceName != nil {
// 			resourceNameParam = *filter.ResourceName
// 		}
// 		if filter.ResourceID != nil {
// 			resourceIDParam = *filter.ResourceID
// 		}
// 	}

// 	err := r.countAuditLogsStmt.QueryRowContext(ctx, performedByParam, userTypeParam, actionParam, resourceNameParam, resourceIDParam).Scan(&count)
// 	if err != nil {
// 		r.logger.Error("Failed to count audit logs", "filter", filter, "error", err)
// 		return 0, err
// 	}

// 	return count, nil
// }
