package service

import (
	"PatientManager/app"
	"PatientManager/model"
	"PatientManager/util/auth"
	"PatientManager/util/cerror"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MobileLoginResult contains tokens for a mobile login
type MobileLoginResult struct {
	AccessToken  string
	RefreshToken string
	DeviceToken  string
}

type ILoginService interface {
	Login(email, password string) (string, string, error)
	RefreshTokens(user *model.User) (string, string, error)
}

type LoginService struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewLoginService() ILoginService {
	var service ILoginService

	app.Invoke(func(db *gorm.DB, logger *zap.SugaredLogger) {
		service = &LoginService{
			db:     db,
			logger: logger,
		}
	})

	return service
}

func (s *LoginService) Login(email, password string) (string, string, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Debugf("User not found Email = %s", email)
			return "", "", cerror.ErrInvalidCredentials
		}

		s.logger.Errorf("Failed to query user, error = %+v", err)
		return "", "", err
	}

	if !auth.VerifyPassword(user.PasswordHash, password) {
		s.logger.Debugf("Invalid password for user Email: %s, uuid: %s", user.Email, user.Uuid)
		return "", "", cerror.ErrInvalidCredentials
	}

	token, refresh, err := auth.GenerateTokens(&user)
	if err != nil {
		s.logger.Errorf("Failed to generate token error = %+v", err)
		return "", "", err
	}

	return token, refresh, nil
}

func (s *LoginService) RefreshTokens(user *model.User) (string, string, error) {
	return auth.GenerateTokens(user)
}
