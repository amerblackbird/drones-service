package services

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"drones/internal/core/domain"
// 	"drones/internal/ports"
// )

// type ActivityLogsService struct {
// 	repo           ports.ActivityLogsRepository
// 	cacheService   ports.CacheService
// 	eventPublisher ports.EventPublisher
// 	logger         ports.Logger
// }

// func NewActivityLogsService(
// 	repo ports.ActivityLogsRepository,
// 	cacheService ports.CacheService,
// 	eventPublisher ports.EventPublisher,
// 	logger ports.Logger,
// ) ports.ActivityLogsService {
// 	return &ActivityLogsService{
// 		repo:           repo,
// 		cacheService:   cacheService,
// 		eventPublisher: eventPublisher,
// 		logger:         logger,
// 	}
// }

// // CreateActivityLog creates a new activity log with business validation
// func (s *ActivityLogsService) CreateActivityLog(ctx context.Context, log *domain.CreateActivityLogRequest) (*domain.ActivityLog, error) {
// 	// Validate required fields
// 	if log.UserID == "" {
// 		s.logger.Error("UserID is required for activity log")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "UserID is required", nil)
// 	}

// 	if log.Action == "" {
// 		s.logger.Error("Action is required for activity log")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "Action is required", nil)
// 	}

// 	// Create the activity log
// 	activityLog, err := s.repo.CreateActivityLog(ctx, log)
// 	if err != nil {
// 		s.logger.Error("Failed to create activity log", "userID", log.UserID, "action", log.Action, "error", err)
// 		return nil, err
// 	}

// 	// Cache the newly created activity log with 1 hour TTL
// 	cacheKey := fmt.Sprintf("activity_log:%s", activityLog.ID)
// 	err = s.cacheService.Set(ctx, cacheKey, *activityLog, 3600) // 1 hour TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache newly created activity log", "logID", activityLog.ID, "error", err)
// 		// Don't fail the operation if caching fails
// 	}

// 	// Publish activity log event for audit purposes
// 	if s.eventPublisher != nil {
// 		err = s.eventPublisher.LogActivity(ctx, log.UserID, log.Action, log.Metadata)
// 		if err != nil {
// 			s.logger.Error("Failed to publish activity log event", "userID", log.UserID, "action", log.Action, "error", err)
// 			// Don't fail the operation if event publishing fails
// 		}
// 	}

// 	s.logger.Info("Activity log created successfully", "logID", activityLog.ID, "userID", log.UserID, "action", log.Action)
// 	return activityLog, nil
// }

// // ListActivityLogs retrieves activity logs with filtering and pagination
// func (s *ActivityLogsService) ListActivityLogs(ctx context.Context, options domain.PaginationOption[domain.ActivityLogFilter]) (*domain.Pagination[domain.ActivityLog], error) {
// 	// Validate pagination parameters
// 	if options.Limit <= 0 {
// 		options.Limit = 10 // Default limit
// 	}
// 	if options.Limit > 100 {
// 		options.Limit = 100 // Maximum limit
// 	}
// 	if options.Offset < 0 {
// 		options.Offset = 0
// 	}

// 	// Check cache for frequently accessed lists
// 	cacheKey := fmt.Sprintf("activity_logs:list:%d:%d", options.Limit, options.Offset)
// 	if options.Filter != nil {
// 		// Add filter parameters to cache key
// 		if options.Filter.UserID != nil {
// 			cacheKey += fmt.Sprintf(":user:%s", *options.Filter.UserID)
// 		}
// 		if options.Filter.Action != nil {
// 			cacheKey += fmt.Sprintf(":action:%s", *options.Filter.Action)
// 		}
// 	}

// 	var cachedResult domain.Pagination[domain.ActivityLog]
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedResult)
// 	if err == nil && len(cachedResult.Data) > 0 {
// 		s.logger.Debug("Activity logs retrieved from cache", "cacheKey", cacheKey)
// 		return &cachedResult, nil
// 	}

// 	// Fetch from repository
// 	result, err := s.repo.ListActivityLogs(ctx, options)
// 	if err != nil {
// 		s.logger.Error("Failed to list activity logs", "filter", options.Filter, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 5 minute TTL
// 	err = s.cacheService.Set(ctx, cacheKey, *result, 300) // 5 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache activity logs list", "cacheKey", cacheKey, "error", err)
// 		// Don't fail the operation if caching fails
// 	}

// 	s.logger.Info("Activity logs retrieved successfully", "count", len(result.Data), "total", result.Total)
// 	return result, nil
// }

// // GetUserActivityLogs retrieves activity logs for a specific user with pagination
// func (s *ActivityLogsService) GetUserActivityLogs(ctx context.Context, userID string, options domain.PaginationOption[domain.ActivityLogFilter]) (*domain.Pagination[domain.ActivityLog], error) {
// 	if userID == "" {
// 		s.logger.Error("UserID is required for user activity logs")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "UserID is required", nil)
// 	}

// 	// Set the user filter
// 	if options.Filter == nil {
// 		options.Filter = &domain.ActivityLogFilter{}
// 	}
// 	options.Filter.UserID = &userID

// 	// Validate pagination parameters
// 	if options.Limit <= 0 {
// 		options.Limit = 10 // Default limit
// 	}
// 	if options.Limit > 50 { // Lower limit for user-specific queries
// 		options.Limit = 50
// 	}
// 	if options.Offset < 0 {
// 		options.Offset = 0
// 	}

// 	// Check cache for user activity logs
// 	cacheKey := fmt.Sprintf("user_activity_logs:%s:%d:%d", userID, options.Limit, options.Offset)
// 	if options.Filter.Action != nil {
// 		cacheKey += fmt.Sprintf(":action:%s", *options.Filter.Action)
// 	}

// 	var cachedResult domain.Pagination[domain.ActivityLog]
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedResult)
// 	if err == nil && len(cachedResult.Data) > 0 {
// 		s.logger.Debug("User activity logs retrieved from cache", "userID", userID, "cacheKey", cacheKey)
// 		return &cachedResult, nil
// 	}

// 	// Fetch from repository using the general ListActivityLogs method
// 	result, err := s.repo.ListActivityLogs(ctx, options)
// 	if err != nil {
// 		s.logger.Error("Failed to get user activity logs", "userID", userID, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 10 minute TTL (longer for user-specific data)
// 	err = s.cacheService.Set(ctx, cacheKey, *result, 600) // 10 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache user activity logs", "userID", userID, "cacheKey", cacheKey, "error", err)
// 		// Don't fail the operation if caching fails
// 	}

// 	s.logger.Info("User activity logs retrieved successfully", "userID", userID, "count", len(result.Data), "total", result.Total)
// 	return result, nil
// }

// // Helper methods for specific activity log queries

// // GetActivityLogsByAction retrieves all activity logs with a specific action
// func (s *ActivityLogsService) GetActivityLogsByAction(ctx context.Context, action string) ([]*domain.ActivityLog, error) {
// 	if action == "" {
// 		s.logger.Error("Action is required")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "Action is required", nil)
// 	}

// 	// Check cache first
// 	cacheKey := fmt.Sprintf("activity_logs:action:%s", action)
// 	var cachedLogs []*domain.ActivityLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLogs)
// 	if err == nil && len(cachedLogs) > 0 {
// 		s.logger.Debug("Activity logs by action retrieved from cache", "action", action)
// 		return cachedLogs, nil
// 	}

// 	// Fetch from repository
// 	logs, err := s.repo.GetByAction(ctx, action)
// 	if err != nil {
// 		s.logger.Error("Failed to get activity logs by action", "action", action, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 15 minute TTL
// 	err = s.cacheService.Set(ctx, cacheKey, logs, 900) // 15 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache activity logs by action", "action", action, "error", err)
// 	}

// 	s.logger.Info("Activity logs by action retrieved successfully", "action", action, "count", len(logs))
// 	return logs, nil
// }

// // GetActivityLogsByResource retrieves all activity logs for a specific resource
// func (s *ActivityLogsService) GetActivityLogsByResource(ctx context.Context, resourceName string, resourceID *string) ([]*domain.ActivityLog, error) {
// 	if resourceName == "" {
// 		s.logger.Error("ResourceName is required")
// 		return nil, domain.NewDomainError(domain.InvalidInputError, "ResourceName is required", nil)
// 	}

// 	// Check cache first
// 	cacheKey := fmt.Sprintf("activity_logs:resource:%s", resourceName)
// 	if resourceID != nil {
// 		cacheKey += fmt.Sprintf(":%s", *resourceID)
// 	}

// 	var cachedLogs []*domain.ActivityLog
// 	err := s.cacheService.Get(ctx, cacheKey, &cachedLogs)
// 	if err == nil && len(cachedLogs) > 0 {
// 		s.logger.Debug("Activity logs by resource retrieved from cache", "resourceName", resourceName, "resourceID", resourceID)
// 		return cachedLogs, nil
// 	}

// 	// Fetch from repository
// 	logs, err := s.repo.GetByResource(ctx, resourceName, resourceID)
// 	if err != nil {
// 		s.logger.Error("Failed to get activity logs by resource", "resourceName", resourceName, "resourceID", resourceID, "error", err)
// 		return nil, err
// 	}

// 	// Cache the result with 15 minute TTL
// 	err = s.cacheService.Set(ctx, cacheKey, logs, 900) // 15 minutes TTL
// 	if err != nil {
// 		s.logger.Error("Failed to cache activity logs by resource", "resourceName", resourceName, "resourceID", resourceID, "error", err)
// 	}

// 	s.logger.Info("Activity logs by resource retrieved successfully", "resourceName", resourceName, "resourceID", resourceID, "count", len(logs))
// 	return logs, nil
// }

// // InvalidateUserActivityLogsCache invalidates cached activity logs for a user
// func (s *ActivityLogsService) InvalidateUserActivityLogsCache(ctx context.Context, userID string) {
// 	// Invalidate user-specific cache entries
// 	cacheKey := fmt.Sprintf("user_activity_logs:%s:*", userID)
// 	err := s.cacheService.Delete(ctx, cacheKey)
// 	if err != nil {
// 		s.logger.Error("Failed to invalidate user activity logs cache", "userID", userID, "error", err)
// 	}
// }

// // GetRecentActivityLogs retrieves recent activity logs (last 24 hours)
// func (s *ActivityLogsService) GetRecentActivityLogs(ctx context.Context, limit int) ([]*domain.ActivityLog, error) {
// 	if limit <= 0 {
// 		limit = 20 // Default limit
// 	}
// 	if limit > 100 {
// 		limit = 100 // Maximum limit
// 	}

// 	// Calculate timestamp for 24 hours ago
// 	now := time.Now()
// 	yesterday := now.Add(-24 * time.Hour)
// 	startTime := yesterday.Unix()
// 	endTime := now.Unix()

// 	filter := &domain.ActivityLogFilter{
// 		StartTime: &startTime,
// 		EndTime:   &endTime,
// 	}

// 	options := domain.PaginationOption[domain.ActivityLogFilter]{
// 		Filter: filter,
// 		Limit:  limit,
// 		Offset: 0,
// 	}

// 	result, err := s.ListActivityLogs(ctx, options)
// 	if err != nil {
// 		s.logger.Error("Failed to get recent activity logs", "error", err)
// 		return nil, err
// 	}

// 	return result.Data, nil
// }
