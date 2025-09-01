package auth

import (
	"PatientManager/config"
	"PatientManager/model"
	"PatientManager/util/cerror"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	Email string         `json:"email"`
	Uuid  string         `json:"uuid"`
	Role  model.UserRole `json:"role"`
}

const (
	accessTokenDuration  = 5 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
	deviceTokenDuration  = 30 * 24 * time.Hour
)

func ParseToken(authHeader string) (*jwt.Token, *Claims, error) {
	// Parse token
	if len(authHeader) <= len("Bearer ") || authHeader[:len("Bearer ")] != "Bearer " {
		zap.S().Debugf("token: %s", authHeader)
		return nil, nil, cerror.ErrInvalidTokenFormat
	}
	tokenString := authHeader[len("Bearer "):]
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(config.AppConfig.AccessKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	return token, &claims, nil
}

func GenerateTokens(user *model.User) (string, string, error) {
	if user == nil {
		return "", "", cerror.ErrUserIsNil
	}

	accessTokenClaims := &Claims{
		Email: user.Email,
		Uuid:  user.Uuid.String(),
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(config.AppConfig.AccessKey))
	if err != nil {
		zap.S().Debugf("Failed to generate access token err = %+v", err)
		return "", "", err
	}

	refreshTokenClaims := &Claims{
		Email: user.Email,
		Uuid:  user.Uuid.String(),
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenDuration)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.AppConfig.RefreshKey))
	if err != nil {
		zap.S().Debugf("Failed to generate refresh token err = %+v", err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func GenerateDeviceToken(user *model.User) (string, error) {
	if user == nil {
		return "", cerror.ErrUserIsNil
	}

	deviceTokenClaims := &Claims{
		Email: user.Email,
		Uuid:  user.Uuid.String(),
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(deviceTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	deviceToken := jwt.NewWithClaims(jwt.SigningMethodHS256, deviceTokenClaims)

	deviceTokenString, err := deviceToken.SignedString([]byte(config.AppConfig.AccessKey))
	if err != nil {
		zap.S().Debugf("Failed to generate device token err = %+v", err)
		return "", err
	}

	return deviceTokenString, nil
}
