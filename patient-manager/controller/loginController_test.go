package controller_test

import (
	"PatientManager/app"
	"PatientManager/config"
	"PatientManager/controller"
	"PatientManager/dto"
	"PatientManager/model"
	"PatientManager/service"
	"PatientManager/util/auth"
	"PatientManager/util/cerror"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// --- Mock LoginService (remains the same) ---
type MockLoginService struct {
	mock.Mock
}

// LoginMobile implements service.ILoginService.
func (m *MockLoginService) LoginMobile(email string, password string) (*service.MobileLoginResult, error) {
	args := m.Called(email, password)
	var res *service.MobileLoginResult
	if v := args.Get(0); v != nil {
		res = v.(*service.MobileLoginResult)
	}
	return res, args.Error(1)
}

func (m *MockLoginService) Login(email, password string) (string, string, error) {
	args := m.Called(email, password)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockLoginService) RefreshTokens(user *model.User) (string, string, error) {
	args := m.Called(user)
	return args.String(0), args.String(1), args.Error(2)
}

// --- LoginController Test Suite ---
type LoginControllerTestSuite struct {
	suite.Suite
	router           *gin.Engine
	mockLoginService *MockLoginService
	sugar            *zap.SugaredLogger
}

// SetupSuite runs once before all tests in the suite
func (suite *LoginControllerTestSuite) SetupSuite() {
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	zapLogger, _ := loggerCfg.Build()
	suite.sugar = zapLogger.Sugar()
	zap.ReplaceGlobals(zapLogger) // For app.Invoke

	gin.SetMode(gin.TestMode)

	config.AppConfig = &config.AppConfiguration{
		Env:        config.Test,
		AccessKey:  "login-ctrl-test-access-key",
		RefreshKey: "login-ctrl-test-refresh-key",
		Port:       8080,
	}

	suite.mockLoginService = new(MockLoginService)

	app.Test() // Initialize DIG container
	app.Provide(func() *zap.SugaredLogger { return suite.sugar })
	app.Provide(func() service.ILoginService { return suite.mockLoginService })

	suite.router = gin.Default()
	apiGroup := suite.router.Group("/api")

	loginCtrl := controller.NewLoginController()
	loginCtrl.RegisterEndpoints(apiGroup)
}

// TearDownSuite runs once after all tests in the suite
func (suite *LoginControllerTestSuite) TearDownSuite() {
	if suite.sugar != nil {
		suite.sugar.Sync()
	}
}

// SetupTest runs before each test in the suite
func (suite *LoginControllerTestSuite) SetupTest() {
	// Reset mocks before each test for isolation
	suite.mockLoginService.ExpectedCalls = nil
	suite.mockLoginService.Calls = nil
}

// TestLoginController runs the test suite
func TestLoginController(t *testing.T) {
	suite.Run(t, new(LoginControllerTestSuite))
}

// --- Test Cases (now methods of the suite) ---

func (suite *LoginControllerTestSuite) TestLogin_Success() {
	loginDto := dto.LoginDto{Email: "test@example.com", Password: "password123"}
	expectedAccessToken := "new.access.token"
	expectedRefreshToken := "new.refresh.token"

	suite.mockLoginService.On("Login", loginDto.Email, loginDto.Password).Return(expectedAccessToken, expectedRefreshToken, nil).Once()

	jsonValue, _ := json.Marshal(loginDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var responseDto dto.TokenDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAccessToken, responseDto.AccessToken)
	assert.Equal(suite.T(), expectedRefreshToken, responseDto.RefreshToken)

	suite.mockLoginService.AssertExpectations(suite.T())
}

func (suite *LoginControllerTestSuite) TestLogin_InvalidCredentials() {
	loginDto := dto.LoginDto{Email: "wrong@example.com", Password: "wrongpassword"}
	suite.mockLoginService.On("Login", loginDto.Email, loginDto.Password).Return("", "", cerror.ErrInvalidCredentials).Once()

	jsonValue, _ := json.Marshal(loginDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), cerror.ErrInvalidCredentials.Error())
	suite.mockLoginService.AssertExpectations(suite.T())
}

func (suite *LoginControllerTestSuite) TestLogin_ServiceError() {
	loginDto := dto.LoginDto{Email: "test@example.com", Password: "password123"}
	suite.mockLoginService.On("Login", loginDto.Email, loginDto.Password).Return("", "", errors.New("internal server error")).Once()

	jsonValue, _ := json.Marshal(loginDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "internal server error")
	suite.mockLoginService.AssertExpectations(suite.T())
}

func (suite *LoginControllerTestSuite) TestLogin_BindingError() {
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email": "test@example.com", "password":`)) // Malformed JSON
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *LoginControllerTestSuite) TestRefreshToken_Success() {
	userUUID := uuid.New()
	userEmail := "refresh@example.com"
	userRole := model.RoleOsoba

	claims := &auth.Claims{
		Email: userEmail,
		Uuid:  userUUID.String(),
		Role:  userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, _ := refreshToken.SignedString([]byte(config.AppConfig.RefreshKey))

	refreshDto := dto.RefreshDto{RefreshToken: refreshTokenString}
	expectedAccessToken := "new.access.token.after.refresh"
	expectedNewRefreshToken := "new.refresh.token.after.refresh"

	suite.mockLoginService.On("RefreshTokens", mock.MatchedBy(func(user *model.User) bool {
		return user.Uuid == userUUID && user.Email == userEmail && user.Role == userRole
	})).Return(expectedAccessToken, expectedNewRefreshToken, nil).Once()

	jsonValue, _ := json.Marshal(refreshDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDto dto.TokenDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAccessToken, responseDto.AccessToken)
	assert.Equal(suite.T(), expectedNewRefreshToken, responseDto.RefreshToken)

	suite.mockLoginService.AssertExpectations(suite.T())
}

func (suite *LoginControllerTestSuite) TestRefreshToken_InvalidToken_Signature() {
	userUUID := uuid.New()
	claims := &auth.Claims{
		Email: "test@example.com", Uuid: userUUID.String(), Role: model.RoleOsoba,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	invalidTokenString, _ := token.SignedString([]byte("wrong-refresh-key"))

	refreshDto := dto.RefreshDto{RefreshToken: invalidTokenString}
	jsonValue, _ := json.Marshal(refreshDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "signature is invalid")
}

func (suite *LoginControllerTestSuite) TestRefreshToken_Expired() {
	userUUID := uuid.New()
	claims := &auth.Claims{
		Email: "test@example.com", Uuid: userUUID.String(), Role: model.RoleOsoba,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, _ := token.SignedString([]byte(config.AppConfig.RefreshKey))

	refreshDto := dto.RefreshDto{RefreshToken: expiredTokenString}
	jsonValue, _ := json.Marshal(refreshDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "token is expired")
}

func (suite *LoginControllerTestSuite) TestRefreshToken_ServiceError() {
	userUUID := uuid.New()
	userEmail := "refresh.fail@example.com"
	userRole := model.RoleFirma
	claims := &auth.Claims{
		Email: userEmail, Uuid: userUUID.String(), Role: userRole,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour))},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validRefreshTokenString, _ := token.SignedString([]byte(config.AppConfig.RefreshKey))

	refreshDto := dto.RefreshDto{RefreshToken: validRefreshTokenString}

	suite.mockLoginService.On("RefreshTokens", mock.MatchedBy(func(user *model.User) bool {
		return user.Uuid == userUUID
	})).Return("", "", errors.New("service failed to refresh")).Once()

	jsonValue, _ := json.Marshal(refreshDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "service failed to refresh")
	suite.mockLoginService.AssertExpectations(suite.T())
}

func (suite *LoginControllerTestSuite) TestRefreshToken_BindingError() {
	req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", strings.NewReader(`{"refreshToken":`)) // Malformed
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
