package services

import (
	"context"

	"drones/internal/core/domain"
	"drones/internal/ports"

	"go.uber.org/zap"
)

// UsersServiceImpl implements the DronesService interface
type UsersServiceImpl struct {
	usersRepo      ports.UserRepository
	eventPublisher ports.EventPublisher
	cacheService   ports.CacheService
	logger         ports.Logger
}

func NewUserRepository(
	userRepo ports.UserRepository,
	eventPublisher ports.EventPublisher,
	cacheService ports.CacheService,
	logger ports.Logger,
) ports.UserService {
	return &UsersServiceImpl{
		usersRepo:      userRepo,
		eventPublisher: eventPublisher, // Emitting events on user actions for example user created or deleted
		cacheService:   cacheService,
		logger:         logger,
	}
}

// GetUserByID implements ports.UserService.
func (s *UsersServiceImpl) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	s.logger.Info("Fetching user by ID", zap.String("userID", userID))
	// Check cache first
	user := new(domain.User)
	cacheKey := "users:" + userID
	err := s.cacheService.Get(ctx, cacheKey, user)
	if err == nil {
		return user, nil
	} else {
		s.logger.Info("User not found in cache", zap.String("userID", userID))
	}

	// Implementation for getting a user by ID
	user, err = s.usersRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Add user to cache for 1 hour
	if err := s.cacheService.Set(ctx, cacheKey, user, 3600); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UsersServiceImpl) GetUserByNameAndType(ctx context.Context, name string, userType string) (*domain.User, error) {
	s.logger.Info("Fetching user by name and type", zap.String("name", name), zap.String("userType", userType))
	// Check cache first
	user := new(domain.User)
	cacheKey := "users:" + userType + ":" + name
	err := s.cacheService.Get(ctx, cacheKey, user)
	if err == nil {
		// User found in cache
		return user, nil
	}
	// Implementation for getting a user by name and type
	user, err = s.usersRepo.GetUserByNameAndType(ctx, name, userType)
	if err != nil {
		// Log error for debugging and return it
		s.logger.Error("Failed to get user from repository", zap.String("name", name), zap.String("userType", userType), zap.Error(err))
		return nil, err
	}

	if err := s.cacheService.Set(ctx, cacheKey, user, 0); err != nil {
		// Log error for debugging but do not fail the request
		s.logger.Error("Failed to set user in cache", zap.String("name", name), zap.String("userType", userType), zap.Error(err))
	}
	return user, nil
}
