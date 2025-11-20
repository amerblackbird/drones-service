// application/auth_service.go
package services

import (
	"context"
	config "drones/configs"

	"drones/internal/core/domain"
	"drones/internal/ports"
)

type AuthServiceImpl struct {
	usersService ports.UserService
	tokenService ports.JwtTokenService
	config       config.JwtConfig
	logger       ports.Logger
}

func NewAuthService(
	usersService ports.UserService,
	tokenService ports.JwtTokenService,
	config config.JwtConfig,
	logger ports.Logger,
) ports.AuthService {
	return &AuthServiceImpl{tokenService: tokenService, config: config, logger: logger, usersService: usersService}
}

func (s *AuthServiceImpl) Login(ctx context.Context, name string, userType string) (*domain.Auth, error) {
	// Get user details
	user, err := s.usersService.GetUserByNameAndType(ctx, name, userType)
	if err != nil {
		s.logger.Error("Failed to get user details", "name", name, "userType", userType, "error", err)
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.tokenService.GenerateToken(ctx, user.ID, userType)
	if err != nil {
		s.logger.Error("Failed to generate tokens", "error", err)
		return nil, err
	}

	// Return auth response
	return &domain.Auth{
		AccessToken: accessToken,
	}, nil
}

func (s *AuthServiceImpl) VerifyToken(ctx context.Context, tokenString string) (string, string, error) {
	return s.tokenService.VerifyToken(ctx, tokenString)
}

func (s *AuthServiceImpl) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.usersService.GetUserByID(ctx, userID)
}
