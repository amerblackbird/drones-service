package services

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"drones/internal/core/domain"
// 	"drones/internal/ports"
// )

// type AuditLogsService struct {
// 	repo           ports.AuditLogsRepository
// 	cacheService   ports.CacheService
// 	eventPublisher ports.EventPublisher
// 	logger         ports.Logger
// }

// func NewAuditLogsService(
// 	repo ports.AuditLogsRepository,
// 	cacheService ports.CacheService,
// 	eventPublisher ports.EventPublisher,
// 	logger ports.Logger,
// ) ports.AuditLogsService {
// 	return &AuditLogsService{
// 		repo:           repo,
// 		cacheService:   cacheService,
// 		eventPublisher: eventPublisher,
// 		logger:         logger,
// 	}
// }

// // CreateAuditLog creates a new audit log with comprehensive validation
// func (s *AuditLogsService) CreateAuditLog(ctx context.Context, log *domain.CreateAuditLogRequest) (*domain.AuditLog, error) {
// 	// Validate required fields
// 	if log.PerformedBy == "" {
// 		s.logger.Error("PerformedBy is required for audit log")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "PerformedBy is required", nil)
// 	}

// 	if log.Action == "" {
// 		s.logger.Error("Action is required for audit log")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "Action is required", nil)
// 	}

// 	if log.ResourceName == "" {
// 		s.logger.Error("ResourceName is required for audit log")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "ResourceName is required", nil)
// 	}

// 	if log.ResourceID == "" {
// 		s.logger.Error("ResourceID is required for audit log")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "ResourceID is required", nil)
// 	}

// 	// Set default values if not provided
// 	if log.UserType == "" {
// 		log.UserType = "user" // Default user type
// 	}

// 	if log.PerformedAt == "" {
// 		log.PerformedAt = time.Now().Format(time.RFC3339) // Default to current time
// 	}

// 	// Validate action is in allowed list
// 	allowedActions := []string{
// 		"create", "update", "delete", "view", "login", "logout",
// 		"approve", "reject", "assign", "unassign", "activate", "deactivate",
// 		"export", "import", "backup", "restore",
// 	}

// 	isValidAction := false
// 	for _, allowedAction := range allowedActions {
// 		if log.Action == allowedAction {
// 			isValidAction = true
// 			break
// 		}
// 	}

// 	if !isValidAction {
// 		s.logger.Warn("Potentially invalid action for audit log", "action", log.Action)
// 		// Don't fail, but log the warning for monitoring
// 	}

// 	// Create the audit log
// 	auditLog, err := s.repo.CreateAuditLog(ctx, log)
// 	if err != nil {
// 		s.logger.Error("Failed to create audit log", "performedBy", log.PerformedBy, "action", log.Action, "resource", log.ResourceName, "error", err)
// 		return nil, err
// 	}

// 	// Cache the newly created audit log with 2 hour TTL (audit logs are more critical)
// 	cacheKey := fmt.Sprintf("audit_log:%s", auditLog.ID)
// 	err = s.cacheService.Set(ctx, cacheKey, *auditLog, 7200) // 2 hours TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache newly created audit log", "logID", auditLog.ID, "error", err)
// 		// Don't fail the operation if caching fails
// 	}

// 	// Invalidate related cache entries
// 	s.invalidateRelatedCaches(ctx, log.PerformedBy, log.ResourceName, log.Action)

// 	s.logger.Info("Audit log created successfully",
// 		"logID", auditLog.ID,
// 		"performedBy", log.PerformedBy,
// 		"action", log.Action,
// 		"resource", log.ResourceName,
// 		"resourceID", log.ResourceID)
// 	return auditLog, nil
// }

// // ListAuditLogs retrieves audit logs with filtering
// func (s *AuditLogsService) ListAuditLogs(ctx context.Context, filter *domain.AuditLogFilter) ([]*domain.AuditLog, error) {
// 	// Create cache key based on filter parameters
// 	cacheKey := "audit_logs:list"
// 	if filter != nil {
// 		if filter.PerformedBy != nil {
// 			cacheKey += fmt.Sprintf(":user:%s", *filter.PerformedBy)
// 		}
// 		if filter.Action != nil {
// 			cacheKey += fmt.Sprintf(":action:%s", *filter.Action)
// 		}
// 		if filter.ResourceName != nil {
// 			cacheKey += fmt.Sprintf(":resource:%s", *filter.ResourceName)
// 		}
// 		if filter.ResourceID != nil {
// 			cacheKey += fmt.Sprintf(":resourceid:%s", *filter.ResourceID)
// 		}
// 	}

// 	// Check cache first
// 	var cachedLogs []*domain.AuditLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLogs)
// 	if err == nil && len(cachedLogs) > 0 {
// 		s.logger.Debug("Audit logs retrieved from cache", "cacheKey", cacheKey)
// 		return cachedLogs, nil
// 	}

// 	// Use pagination with default values for the list method
// 	paginationOptions := domain.PaginationOption[domain.AuditLogFilter]{
// 		Filter: filter,
// 		Limit:  1000, // Large limit for list method
// 		Offset: 0,
// 	}

// 	// Fetch from repository using pagination
// 	result, err := s.repo.ListAuditLogs(ctx, paginationOptions)
// 	if err != nil {
// 		s.logger.Error("Failed to list audit logs", "filter", filter, "error", err)
// 		return nil, err
// 	}

// 	// Extract data from pagination result
// 	auditLogs := result.Data

// 	// Cache the result with 30 minute TTL (audit logs change less frequently)
// 	err = s.cacheService.Set(ctx, cacheKey, auditLogs, 1800) // 30 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache audit logs list", "cacheKey", cacheKey, "error", err)
// 		// Don't fail the operation if caching fails
// 	}

// 	s.logger.Info("Audit logs retrieved successfully", "count", len(auditLogs), "filter", filter)
// 	return auditLogs, nil
// }

// // Additional business methods for audit logs

// // GetAuditLogByID retrieves a single audit log by ID
// func (s *AuditLogsService) GetAuditLogByID(ctx context.Context, logID string) (*domain.AuditLog, error) {
// 	if logID == "" {
// 		s.logger.Error("LogID is required")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "LogID is required", nil)
// 	}

// 	// Check cache first
// 	cacheKey := fmt.Sprintf("audit_log:%s", logID)
// 	var cachedLog domain.AuditLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLog)
// 	if err == nil && cachedLog.ID != "" {
// 		s.logger.Debug("Audit log retrieved from cache", "logID", logID)
// 		return &cachedLog, nil
// 	}

// 	// Fetch from repository
// 	auditLog, err := s.repo.GetByID(ctx, logID)
// 	if err != nil {
// 		s.logger.Error("Failed to get audit log by ID", "logID", logID, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result
// 	err = s.cacheService.Set(ctx, cacheKey, *auditLog, 7200) // 2 hours TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache audit log", "logID", logID, "error", err)
// 	}

// 	s.logger.Info("Audit log retrieved successfully", "logID", logID)
// 	return auditLog, nil
// }

// // GetAuditLogsByPerformedBy retrieves audit logs by user
// func (s *AuditLogsService) GetAuditLogsByPerformedBy(ctx context.Context, performedBy string) ([]*domain.AuditLog, error) {
// 	if performedBy == "" {
// 		s.logger.Error("PerformedBy is required")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "PerformedBy is required", nil)
// 	}

// 	// Check cache first
// 	cacheKey := fmt.Sprintf("audit_logs:user:%s", performedBy)
// 	var cachedLogs []*domain.AuditLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLogs)
// 	if err == nil && len(cachedLogs) > 0 {
// 		s.logger.Debug("User audit logs retrieved from cache", "performedBy", performedBy)
// 		return cachedLogs, nil
// 	}

// 	// Fetch from repository
// 	logs, err := s.repo.GetByPerformedBy(ctx, performedBy)
// 	if err != nil {
// 		s.logger.Error("Failed to get audit logs by performed by", "performedBy", performedBy, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 20 minute TTL
// 	err = s.cacheService.Set(ctx, cacheKey, logs, 1200) // 20 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache user audit logs", "performedBy", performedBy, "error", err)
// 	}

// 	s.logger.Info("User audit logs retrieved successfully", "performedBy", performedBy, "count", len(logs))
// 	return logs, nil
// }

// // GetAuditLogsByAction retrieves audit logs by action
// func (s *AuditLogsService) GetAuditLogsByAction(ctx context.Context, action string) ([]*domain.AuditLog, error) {
// 	if action == "" {
// 		s.logger.Error("Action is required")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "Action is required", nil)
// 	}

// 	// Check cache first
// 	cacheKey := fmt.Sprintf("audit_logs:action:%s", action)
// 	var cachedLogs []*domain.AuditLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLogs)
// 	if err == nil && len(cachedLogs) > 0 {
// 		s.logger.Debug("Action audit logs retrieved from cache", "action", action)
// 		return cachedLogs, nil
// 	}

// 	// Fetch from repository
// 	logs, err := s.repo.GetByAction(ctx, action)
// 	if err != nil {
// 		s.logger.Error("Failed to get audit logs by action", "action", action, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 25 minute TTL
// 	err = s.cacheService.Set(ctx, cacheKey, logs, 1500) // 25 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache action audit logs", "action", action, "error", err)
// 	}

// 	s.logger.Info("Action audit logs retrieved successfully", "action", action, "count", len(logs))
// 	return logs, nil
// }

// // GetAuditLogsByResource retrieves audit logs for a specific resource
// func (s *AuditLogsService) GetAuditLogsByResource(ctx context.Context, resourceName string, resourceID *string) ([]*domain.AuditLog, error) {
// 	if resourceName == "" {
// 		s.logger.Error("ResourceName is required")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "ResourceName is required", nil)
// 	}

// 	// Check cache first
// 	cacheKey := fmt.Sprintf("audit_logs:resource:%s", resourceName)
// 	if resourceID != nil {
// 		cacheKey += fmt.Sprintf(":%s", *resourceID)
// 	}

// 	var cachedLogs []*domain.AuditLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLogs)
// 	if err == nil && len(cachedLogs) > 0 {
// 		s.logger.Debug("Resource audit logs retrieved from cache", "resourceName", resourceName, "resourceID", resourceID)
// 		return cachedLogs, nil
// 	}

// 	// Fetch from repository
// 	logs, err := s.repo.GetByResource(ctx, resourceName, resourceID)
// 	if err != nil {
// 		s.logger.Error("Failed to get audit logs by resource", "resourceName", resourceName, "resourceID", resourceID, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 25 minute TTL
// 	err = s.cacheService.Set(ctx, cacheKey, logs, 1500) // 25 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache resource audit logs", "resourceName", resourceName, "resourceID", resourceID, "error", err)
// 	}

// 	s.logger.Info("Resource audit logs retrieved successfully", "resourceName", resourceName, "resourceID", resourceID, "count", len(logs))
// 	return logs, nil
// }

// // GetRecentAuditLogs retrieves recent audit logs (last 7 days)
// func (s *AuditLogsService) GetRecentAuditLogs(ctx context.Context, days int) ([]*domain.AuditLog, error) {
// 	if days <= 0 {
// 		days = 7 // Default to 7 days
// 	}
// 	if days > 30 {
// 		days = 30 // Maximum 30 days
// 	}

// 	// Check cache first
// 	cacheKey := fmt.Sprintf("audit_logs:recent:%d_days", days)
// 	var cachedLogs []*domain.AuditLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLogs)
// 	if err == nil && len(cachedLogs) > 0 {
// 		s.logger.Debug("Recent audit logs retrieved from cache", "days", days)
// 		return cachedLogs, nil
// 	}

// 	// For simplicity, we'll use the list method without date filtering
// 	// In a real implementation, you'd add date range filtering to the repository
// 	logs, err := s.ListAuditLogs(ctx, nil)
// 	if err != nil {
// 		s.logger.Error("Failed to get recent audit logs", "days", days, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 1 hour TTL (recent data changes more frequently)
// 	err = s.cacheService.Set(ctx, cacheKey, logs, 3600) // 1 hour TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache recent audit logs", "days", days, "error", err)
// 	}

// 	s.logger.Info("Recent audit logs retrieved successfully", "days", days, "count", len(logs))
// 	return logs, nil
// }

// // CountAuditLogsByFilter counts audit logs matching the filter
// func (s *AuditLogsService) CountAuditLogsByFilter(ctx context.Context, filter *domain.AuditLogFilter) (int64, error) {
// 	// Check cache first
// 	cacheKey := "audit_logs:count"
// 	if filter != nil {
// 		if filter.PerformedBy != nil {
// 			cacheKey += fmt.Sprintf(":user:%s", *filter.PerformedBy)
// 		}
// 		if filter.Action != nil {
// 			cacheKey += fmt.Sprintf(":action:%s", *filter.Action)
// 		}
// 		if filter.ResourceName != nil {
// 			cacheKey += fmt.Sprintf(":resource:%s", *filter.ResourceName)
// 		}
// 	}

// 	var cachedCount int64
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedCount)
// 	if err == nil && cachedCount > 0 {
// 		s.logger.Debug("Audit logs count retrieved from cache", "cacheKey", cacheKey)
// 		return cachedCount, nil
// 	}

// 	// Fetch from repository
// 	count, err := s.repo.CountAuditLogs(ctx, filter)
// 	if err != nil {
// 		s.logger.Error("Failed to count audit logs", "filter", filter, "error", err)
// 		return 0, err
// 	}

// 	// Cache the result with 15 minute TTL
// 	err = s.cacheService.Set(ctx, cacheKey, count, 900) // 15 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache audit logs count", "cacheKey", cacheKey, "error", err)
// 	}

// 	s.logger.Info("Audit logs count retrieved successfully", "count", count, "filter", filter)
// 	return count, nil
// }

// // Helper method to invalidate related cache entries
// func (s *AuditLogsService) invalidateRelatedCaches(ctx context.Context, performedBy, resourceName, action string) {
// 	// List of cache patterns to invalidate
// 	cachePatterns := []string{
// 		"audit_logs:list*",
// 		fmt.Sprintf("audit_logs:user:%s", performedBy),
// 		fmt.Sprintf("audit_logs:resource:%s*", resourceName),
// 		fmt.Sprintf("audit_logs:action:%s", action),
// 		"audit_logs:recent:*",
// 		"audit_logs:count*",
// 	}

// 	for _, pattern := range cachePatterns {
// 		err := s.cacheService.Delete(ctx, pattern)
// 		if err != nil {
// 			s.logger.Error("Failed to invalidate cache", "pattern", pattern, "error", err)
// 		}
// 	}
// }

// // CreateChangeAuditLog is a helper method for tracking entity changes
// func (s *AuditLogsService) CreateChangeAuditLog(ctx context.Context, performedBy, userType, action, resourceName, resourceID string, oldData, newData interface{}) (*domain.AuditLog, error) {
// 	log := &domain.CreateAuditLogRequest{
// 		PerformedBy:  performedBy,
// 		UserType:     userType,
// 		Action:       action,
// 		ResourceName: resourceName,
// 		ResourceID:   resourceID,
// 		PerformedAt:  time.Now().Format(time.RFC3339),
// 		OldData:      oldData,
// 		NewData:      newData,
// 	}

// 	return s.CreateAuditLog(ctx, log)
// }
