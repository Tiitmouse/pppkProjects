package controller_test

import (
	"PatientManager/app"
	"PatientManager/config"
	"PatientManager/controller"
	"PatientManager/dto"
	"PatientManager/model"
	"PatientManager/service"
	"PatientManager/util/auth"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

// --- Mock UserCrudService ---
type MockUserCrudService struct {
	mock.Mock
}

func (m *MockUserCrudService) Create(user *model.User, password string) (*model.User, error) {
	args := m.Called(user, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserCrudService) Read(id uuid.UUID) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserCrudService) ReadAll() ([]model.User, error) { // This specific ReadAll is not used by UserController
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserCrudService) Update(id uuid.UUID, user *model.User) (*model.User, error) {
	args := m.Called(id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserCrudService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserCrudService) GetAllUsers() ([]model.User, error) { // For SuperAdmin
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserCrudService) GetAllPoliceOfficers() ([]model.User, error) { // For MupAdmin
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserCrudService) SearchUsersByName(query string) ([]model.User, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserCrudService) GetUserByOIB(oib string) (*model.User, error) {
	args := m.Called(oib)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// --- UserController Test Suite ---
type UserControllerTestSuite struct {
	suite.Suite
	router              *gin.Engine
	mockUserCrudService *MockUserCrudService
	sugar               *zap.SugaredLogger
}

// SetupSuite runs once before all tests in the suite
func (suite *UserControllerTestSuite) SetupSuite() {
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	zapLogger, _ := loggerCfg.Build()
	suite.sugar = zapLogger.Sugar()
	zap.ReplaceGlobals(zapLogger)

	gin.SetMode(gin.TestMode)

	config.AppConfig = &config.AppConfiguration{
		Env:        config.Test,
		AccessKey:  "user-ctrl-test-access-key",
		RefreshKey: "user-ctrl-test-refresh-key",
	}

	suite.mockUserCrudService = new(MockUserCrudService)

	app.Test()
	app.Provide(func() *zap.SugaredLogger { return suite.sugar })
	app.Provide(func() service.IUserCrudService { return suite.mockUserCrudService })

	suite.router = gin.Default()
	apiGroup := suite.router.Group("/api") // Assuming routes are under /api

	userCtrl := controller.NewUserController()
	userCtrl.RegisterEndpoints(apiGroup)
}

// TearDownSuite runs once after all tests
func (suite *UserControllerTestSuite) TearDownSuite() {
	if suite.sugar != nil {
		_ = suite.sugar.Sync()
	}
}

// SetupTest runs before each test
func (suite *UserControllerTestSuite) SetupTest() {
	suite.mockUserCrudService.ExpectedCalls = nil
	suite.mockUserCrudService.Calls = nil
}

// Helper to generate a token for a test user
func generateUserTestToken(userID uuid.UUID, userEmail string, userRole model.UserRole) string {
	token, _, _ := auth.GenerateTokens(&model.User{
		Uuid:  userID,
		Email: userEmail,
		Role:  userRole,
	})
	return token
}

// TestUserController runs the test suite
func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

// --- Test Cases ---

func (suite *UserControllerTestSuite) TestCreateUser_Success() {
	adminToken := generateUserTestToken(uuid.New(), "admin@example.com", model.RoleSuperAdmin)
	newUserDto := dto.NewUserDto{
		FirstName: "Test", LastName: "User", OIB: "12345678901",
		Residence: "Testville", BirthDate: "1990-01-01", Email: "test@example.com",
		Password: "password123", Role: "osoba",
	}
	expectedUserUUID := uuid.New()
	expectedUserModel, _ := newUserDto.ToModel() // DTO to model handles most fields
	expectedUserModel.Uuid = expectedUserUUID    // Set the expected UUID

	// Mock service Create
	suite.mockUserCrudService.On("Create", mock.AnythingOfType("*model.User"), newUserDto.Password).
		Run(func(args mock.Arguments) {
			argUser := args.Get(0).(*model.User)
			// Assertions on the user passed to the service
			assert.Equal(suite.T(), newUserDto.FirstName, argUser.FirstName)
			assert.Equal(suite.T(), newUserDto.Email, argUser.Email)
		}).Return(expectedUserModel, nil).Once()

	jsonValue, _ := json.Marshal(newUserDto)
	req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	var responseDto dto.UserDto // Note: UserController.create returns UserDto, not NewUserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUserUUID.String(), responseDto.Uuid)
	assert.Equal(suite.T(), newUserDto.FirstName, responseDto.FirstName)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestCreateUser_Forbidden() {
	nonAdminToken := generateUserTestToken(uuid.New(), "user@example.com", model.RoleOsoba)
	newUserDto := dto.NewUserDto{FirstName: "Test", Role: "osoba"} // Simplified DTO
	jsonValue, _ := json.Marshal(newUserDto)

	req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+nonAdminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)

	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestCreateUser_BindingError() {
	adminToken := generateUserTestToken(uuid.New(), "admin@example.com", model.RoleSuperAdmin)
	req, _ := http.NewRequest(http.MethodPost, "/api/user/", strings.NewReader(`{}`)) // Malformed
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)

	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetUser_Success() {
	adminToken := generateUserTestToken(uuid.New(), "admin@example.com", model.RoleSuperAdmin)
	targetUserUUID := uuid.New()
	expectedUser := &model.User{
		Uuid: targetUserUUID, FirstName: "Target", LastName: "User", Role: model.RoleOsoba,
		BirthDate: time.Now().AddDate(-20, 0, 0), OIB: "98765432109", Email: "target@example.com",
	}
	suite.mockUserCrudService.On("Read", targetUserUUID).Return(expectedUser, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/"+targetUserUUID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDto dto.UserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), targetUserUUID.String(), responseDto.Uuid)
	assert.Equal(suite.T(), expectedUser.FirstName, responseDto.FirstName)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetUser_NotFound() {
	adminToken := generateUserTestToken(uuid.New(), "admin@example.com", model.RoleSuperAdmin)
	targetUserUUID := uuid.New()
	suite.mockUserCrudService.On("Read", targetUserUUID).Return(nil, gorm.ErrRecordNotFound).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/"+targetUserUUID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestUpdateUser_Success() {
	adminToken := generateUserTestToken(uuid.New(), "admin@example.com", model.RoleSuperAdmin)
	targetUserUUID := uuid.New()
	updateDto := dto.UserDto{ // UserDto is used for update
		Uuid: targetUserUUID.String(), FirstName: "UpdatedFirst", LastName: "UpdatedLast",
		OIB: "11122233344", Residence: "Updated Residence", BirthDate: "1985-05-15",
		Email: "updated.email@example.com", Role: "firma",
	}
	updatedUserModel, _ := updateDto.ToModel() // DTO to model

	suite.mockUserCrudService.On("Update", targetUserUUID, mock.MatchedBy(func(u *model.User) bool {
		return u.FirstName == updateDto.FirstName && u.Email == updateDto.Email
	})).Return(updatedUserModel, nil).Once()

	jsonValue, _ := json.Marshal(updateDto)
	req, _ := http.NewRequest(http.MethodPut, "/api/user/"+targetUserUUID.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDto dto.UserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updateDto.FirstName, responseDto.FirstName)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestDeleteUser_Success() {
	adminToken := generateUserTestToken(uuid.New(), "admin@example.com", model.RoleSuperAdmin)
	targetUserUUID := uuid.New()
	suite.mockUserCrudService.On("Delete", targetUserUUID).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/api/user/"+targetUserUUID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetLoggedInUser_Success() {
	loggedInUserUUID := uuid.New()
	loggedInUserEmail := "loggedin@example.com"
	loggedInUserRole := model.RoleOsoba
	token := generateUserTestToken(loggedInUserUUID, loggedInUserEmail, loggedInUserRole)

	expectedUser := &model.User{
		Uuid: loggedInUserUUID, FirstName: "Logged", LastName: "In", Role: loggedInUserRole,
		BirthDate: time.Now().AddDate(-20, 0, 0), OIB: "55566677788", Email: loggedInUserEmail,
	}
	suite.mockUserCrudService.On("Read", loggedInUserUUID).Return(expectedUser, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/my-data", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDto dto.UserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), loggedInUserUUID.String(), responseDto.Uuid)
	assert.Equal(suite.T(), expectedUser.FirstName, responseDto.FirstName)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetAllUsersForSuperAdmin_Success() {
	superAdminToken := generateUserTestToken(uuid.New(), "super@admin.com", model.RoleSuperAdmin)
	users := []model.User{
		{Uuid: uuid.New(), FirstName: "UserA", Role: model.RoleOsoba, BirthDate: time.Now()},
		{Uuid: uuid.New(), FirstName: "UserB", Role: model.RoleFirma, BirthDate: time.Now()},
	}
	suite.mockUserCrudService.On("GetAllUsers").Return(users, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/all-users", nil)
	req.Header.Set("Authorization", "Bearer "+superAdminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDtos []dto.UserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDtos)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), responseDtos, 2)
	assert.Equal(suite.T(), users[0].FirstName, responseDtos[0].FirstName)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetAllPoliceOfficers_Success() {
	mupAdminToken := generateUserTestToken(uuid.New(), "mup@admin.com", model.RoleMupADMIN)
	policeOfficers := []model.User{
		{Uuid: uuid.New(), FirstName: "OfficerA", Role: model.RolePolicija, BirthDate: time.Now()},
		{Uuid: uuid.New(), FirstName: "OfficerB", Role: model.RolePolicija, BirthDate: time.Now()},
	}
	suite.mockUserCrudService.On("GetAllPoliceOfficers").Return(policeOfficers, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/police-officers", nil)
	req.Header.Set("Authorization", "Bearer "+mupAdminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDtos []dto.UserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDtos)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), responseDtos, 2)
	assert.Equal(suite.T(), policeOfficers[0].FirstName, responseDtos[0].FirstName)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestSearchUsersByName_Success() {
	adminToken := generateUserTestToken(uuid.New(), "adminsearch@example.com", model.RoleSuperAdmin) // Assuming SuperAdmin can search
	searchQuery := "John"
	foundUsers := []model.User{
		{Uuid: uuid.New(), FirstName: "John", LastName: "Doe", Role: model.RoleOsoba, BirthDate: time.Now()},
		{Uuid: uuid.New(), FirstName: "Johnny", LastName: "Smith", Role: model.RoleFirma, BirthDate: time.Now()},
	}
	suite.mockUserCrudService.On("SearchUsersByName", searchQuery).Return(foundUsers, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/search?query="+searchQuery, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDtos []dto.UserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDtos)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), responseDtos, 2)
	assert.Equal(suite.T(), foundUsers[0].FirstName, responseDtos[0].FirstName)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestSearchUsersByName_EmptyQuery() {
	adminToken := generateUserTestToken(uuid.New(), "adminsearch@example.com", model.RoleSuperAdmin)
	req, _ := http.NewRequest(http.MethodGet, "/api/user/search?query=", nil) // Empty query
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Search query is required")
}

func (suite *UserControllerTestSuite) TestGetUserByOib_Success() {
	hakToken := generateUserTestToken(uuid.New(), "hak@example.com", model.RoleHAK)
	targetOIB := "11223344556"
	expectedUser := &model.User{
		Uuid: uuid.New(), FirstName: "OIB", LastName: "User", Role: model.RoleOsoba,
		BirthDate: time.Now().AddDate(-30, 0, 0), OIB: targetOIB, Email: "oib.user@example.com",
	}
	suite.mockUserCrudService.On("GetUserByOIB", targetOIB).Return(expectedUser, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/oib/"+targetOIB, nil)
	req.Header.Set("Authorization", "Bearer "+hakToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var responseDto dto.UserDto
	err := json.Unmarshal(w.Body.Bytes(), &responseDto)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedUser.Uuid.String(), responseDto.Uuid)
	assert.Equal(suite.T(), expectedUser.OIB, responseDto.OIB)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestGetUserByOib_NotFound() {
	hakToken := generateUserTestToken(uuid.New(), "hak@example.com", model.RoleHAK)
	targetOIB := "00000000000" // Non-existent OIB
	suite.mockUserCrudService.On("GetUserByOIB", targetOIB).Return(nil, gorm.ErrRecordNotFound).Once()

	req, _ := http.NewRequest(http.MethodGet, "/api/user/oib/"+targetOIB, nil)
	req.Header.Set("Authorization", "Bearer "+hakToken)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.mockUserCrudService.AssertExpectations(suite.T())
}
