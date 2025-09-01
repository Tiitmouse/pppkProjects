package service_test

import (
	"PatientManager/app"
	"PatientManager/config"
	"PatientManager/model"
	"PatientManager/service"
	"PatientManager/util/auth"
	"PatientManager/util/cerror"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// --- LoginService Test Suite ---
type LoginServiceTestSuite struct {
	suite.Suite
	db           *gorm.DB
	loginService service.ILoginService
	// No mock for deviceManager, as LoginService instantiates it directly.
	// We will test its interaction with the DB.
	logger      *zap.SugaredLogger
	logObserver *observer.ObservedLogs
}

// SetupSuite runs once before all tests in the suite
func (suite *LoginServiceTestSuite) SetupSuite() {
	core, obs := observer.New(zap.InfoLevel)
	suite.logger = zap.New(core).Sugar()
	suite.logObserver = obs
	zap.ReplaceGlobals(zap.New(core)) // For app.Invoke

	db, err := gorm.Open(sqlite.Open("file:loginservice_test.db?mode=memory&cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	suite.Require().NoError(err, "Failed to connect to SQLite")
	suite.db = db

	err = suite.db.AutoMigrate(model.GetAllModels()...)
	suite.Require().NoError(err, "Failed to migrate database schema")

	config.AppConfig = &config.AppConfiguration{
		Env:        config.Test,
		AccessKey:  "login-service-test-access-key",
		RefreshKey: "login-service-test-refresh-key",
	}

	// Setup DIG container and provide dependencies
	app.Test() // Initialize DIG container
	app.Provide(func() *gorm.DB { return suite.db })
	app.Provide(func() *zap.SugaredLogger { return suite.logger })
	// LoginService is what we are testing, so we get it from DIG after providing its dependencies
	suite.loginService = service.NewLoginService()
	suite.Require().NotNil(suite.loginService, "LoginService should be initialized by DIG")
}

// TearDownSuite runs once after all tests
func (suite *LoginServiceTestSuite) TearDownSuite() {
	if suite.logger != nil {
		_ = suite.logger.Sync()
	}
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		err := sqlDB.Close()
		suite.Require().NoError(err)
	}
}

// clearTables helper
func (suite *LoginServiceTestSuite) clearTables(tables ...string) {
	suite.db.Exec("PRAGMA foreign_keys = OFF")
	defer suite.db.Exec("PRAGMA foreign_keys = ON")
	for _, table := range tables {
		var modelInstance interface{}
		switch table {
		case "users":
			modelInstance = &model.User{}
		default:
			suite.T().Fatalf("Unsupported table for clearing: %s", table)
		}
		err := suite.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(modelInstance).Error
		suite.Require().NoError(err, fmt.Sprintf("Failed to clear table %s", table))
	}
}

// SetupTest runs before each test
func (suite *LoginServiceTestSuite) SetupTest() {
	suite.logObserver.TakeAll() // Clear observed logs
}

// Helper to create a user with a hashed password
func (suite *LoginServiceTestSuite) createTestUser(email, plainPassword string, role model.UserRole) *model.User {
	hashedPassword, err := auth.HashPassword(plainPassword)
	suite.Require().NoError(err)
	user := &model.User{
		Uuid:         uuid.New(),
		FirstName:    "Test",
		LastName:     string(role),
		OIB:          uuid.New().String()[:11], // Unique OIB
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
		BirthDate:    time.Now().AddDate(-25, 0, 0),
		Residence:    "Test Residence",
	}
	err = suite.db.Create(user).Error
	suite.Require().NoError(err, "Failed to create user for test")
	return user
}

// TestLoginServiceSuite runs the test suite
func TestLoginServiceSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceTestSuite))
}

// --- Test Cases ---

func (suite *LoginServiceTestSuite) TestLogin_Success() {
	email := "testlogin@example.com"
	password := "password123"
	suite.createTestUser(email, password, model.RoleOsoba)

	accessToken, refreshToken, err := suite.loginService.Login(email, password)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)
}

func (suite *LoginServiceTestSuite) TestLogin_UserNotFound() {
	accessToken, refreshToken, err := suite.loginService.Login("nonexistent@example.com", "password123")

	assert.Error(suite.T(), err)
	assert.True(suite.T(), errors.Is(err, cerror.ErrInvalidCredentials))
	assert.Empty(suite.T(), accessToken)
	assert.Empty(suite.T(), refreshToken)
}

func (suite *LoginServiceTestSuite) TestLogin_IncorrectPassword() {
	email := "wrongpass@example.com"
	password := "password123"
	suite.createTestUser(email, password, model.RoleOsoba)

	accessToken, refreshToken, err := suite.loginService.Login(email, "wrongPassword")

	assert.Error(suite.T(), err)
	assert.True(suite.T(), errors.Is(err, cerror.ErrInvalidCredentials))
	assert.Empty(suite.T(), accessToken)
	assert.Empty(suite.T(), refreshToken)
}

func (suite *LoginServiceTestSuite) TestRefreshTokens_Success() {
	user := suite.createTestUser("refreshuser@example.com", "password123", model.RoleFirma)

	accessToken, refreshToken, err := suite.loginService.RefreshTokens(user)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)

	// Further validation: parse tokens and check claims (optional, as auth.GenerateTokens is tested separately)
	_, accessClaims, errAccess := auth.ParseToken("Bearer " + accessToken)
	suite.Require().NoError(errAccess)
	assert.Equal(suite.T(), user.Uuid.String(), accessClaims.Uuid)

	// For refresh token, need to parse with refresh key
	var refreshClaims auth.Claims
	_, errRefresh := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.RefreshKey), nil
	})
	suite.Require().NoError(errRefresh)
	assert.Equal(suite.T(), user.Uuid.String(), refreshClaims.Uuid)
}
