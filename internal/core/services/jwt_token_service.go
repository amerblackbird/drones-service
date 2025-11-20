package services

import (
	"context"
	"errors"
	"time"

	config "drones/configs"
	"drones/internal/ports"
	utils "drones/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenServiceImpl struct {
	config *config.JwtConfig
}

func NewJWTService(config *config.JwtConfig) ports.JwtTokenService {
	return &JwtTokenServiceImpl{config: config}
}

// Helper to parse duration strings like "1h", "30m", etc.
func parseExpireDuration(expireStr string) (int64, error) {
	dur, err := time.ParseDuration(expireStr)
	if err != nil {
		return 0, err
	}
	return time.Now().Add(dur).Unix(), nil
}

func (j *JwtTokenServiceImpl) GenerateToken(ctx context.Context, userID string, userType string) (string, error) {

	accessExp, err := parseExpireDuration(j.config.ExpiresIn)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"utype": userType,
		"exp":   accessExp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (j *JwtTokenServiceImpl) VerifyToken(ctx context.Context, tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	userID, _ := claims["sub"].(string)
	userType, _ := claims["utype"].(string)

	return userID, userType, nil
}

func (j *JwtTokenServiceImpl) GenerateVerfiyCred() (string, string, string, string) {
	return utils.GenerateVerfiyCred()
}

func (j *JwtTokenServiceImpl) VerifyOtpCode(hash, salt, code string) bool {
	return utils.VerifyOtpHash(code, salt, hash)
}
