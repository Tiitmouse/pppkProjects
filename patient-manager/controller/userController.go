package controller

import (
	"PatientManager/app"
	"PatientManager/dto"
	"PatientManager/service"
	"PatientManager/util/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserController struct {
	UserCrud service.IUserCrudService
	logger   *zap.SugaredLogger
}

func NewUserController() *UserController {
	var controller *UserController

	// Call dependency injection
	app.Invoke(func(UserService service.IUserCrudService, logger *zap.SugaredLogger) {
		// create controller
		controller = &UserController{
			UserCrud: UserService,
			logger:   logger,
		}
	})

	return controller
}

func (u *UserController) RegisterEndpoints(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.POST("", u.create)
		user.GET("/:uuid", u.get)
		user.GET("/my-data", u.getLoggedInUser)
		user.PUT("/:uuid", u.update)
		user.DELETE("/:uuid", u.delete)
	}
}

// UserExample godoc
//
//	@Summary		get user with uuid
//	@Description	get a user with uuid
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	dto.UserDto
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			uuid	path	string	true	"user uuid"
//	@Router			/user/{uuid} [get]
func (u *UserController) get(c *gin.Context) {
	userUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		u.logger.Errorf("error parsing uuid value = %s", c.Param("uuid"))
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := u.UserCrud.Read(userUuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Errorf("User with uuid = %s not found", userUuid)
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		u.logger.Errorf("Failed to get user with uuid = %s", userUuid)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dto := dto.UserDto{}
	c.JSON(http.StatusOK, dto.FromModel(user))
}

// UserExample godoc
//
//	@Summary	Create new user
//	@Tags		user
//	@Produce	json
//	@Success	201	{object}	dto.UserDto
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Param		model	body	dto.NewUserDto	true	"Data for new user"
//	@Router		/user [post]
func (u *UserController) create(c *gin.Context) {
	var dto dto.NewUserDto
	if err := c.BindJSON(&dto); err != nil {
		u.logger.Errorf("Failed to bind error = %+v", err)
		return
	}

	newUser, err := dto.ToModel()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := u.UserCrud.Create(newUser, "Pa$$w0rd")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Log the response for debugging
	responseDto := dto.FromModel(user)
	u.logger.Infof("Response DTO: %+v", responseDto)

	c.JSON(http.StatusCreated, responseDto)
}

// UserExample godoc
//
//	@Summary	Update user with new dat
//	@Tags		user
//	@Produce	json
//	@Success	200	{object}	dto.UserDto
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Param		uuid	path	string		true	"uuid of user to be updated"
//	@Param		model	body	dto.UserDto	true	"Data for updating user"
//	@Router		/user/{uuid} [put]
func (u *UserController) update(c *gin.Context) {
	userUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		u.logger.Errorf("Error parsing UUID = %s", c.Param("uuid"))
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var dto dto.UserDto
	if err := c.BindJSON(&dto); err != nil {
		u.logger.Errorf("Failed to bind error = %+v", err)
		return
	}

	newUser, err := dto.ToModel()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := u.UserCrud.Update(userUuid, newUser)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, dto.FromModel(user))
}

// UserExample  godoc
//
//	@Summary		delete user with uuid
//	@Description	delete a user with uuid
//	@Tags			user
//	@Produce		json
//	@Success		204
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			uuid	path	string	true	"user uuid"
//	@Router			/user/{uuid} [delete]
func (u *UserController) delete(c *gin.Context) {
	userUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		u.logger.Errorf("error parsing uuid value = %s", c.Param("uuid"))
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = u.UserCrud.Delete(userUuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Errorf("User with uuid = %s not found", userUuid)
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		u.logger.Errorf("Failed to delete user with uuid = %s", userUuid)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetLoggedInUser godoc
//
//	@Summary		Get logged-in user data
//	@Description	Fetches the currently logged-in user's data based on the JWT token
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	dto.UserDto
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/user/my-data [get]
func (u *UserController) getLoggedInUser(c *gin.Context) {
	_, claims, err := auth.ParseToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		u.logger.Errorf("Failed to parse token: %v", err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	userUuid, err := uuid.Parse(claims.Uuid)
	if err != nil {
		u.logger.Errorf("Error parsing UUID = %s", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := u.UserCrud.Read(userUuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.logger.Errorf("User with uuid = %s not found", userUuid)
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		u.logger.Errorf("Failed to fetch user with uuid = %s: %v", userUuid, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	dto := dto.UserDto{}
	c.JSON(http.StatusOK, dto.FromModel(user))
}

// SearchUsersByName godoc
//
//	@Summary		Search users by name
//	@Description	Performs a fuzzy search for users by first name, last name, or full name with similarity matching
//	@Tags			user
//	@Produce		json
//	@Param			query	query	string	true	"Search query"
//	@Success		200		{array}	dto.UserDto
//	@Failure		400
//	@Failure		500
//	@Router			/user/search [get]
func (u *UserController) searchUsersByName(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		u.logger.Warn("Search query is empty")
		c.JSON(http.StatusBadRequest, "Search query is required")
		return
	}

	u.logger.Infof("Searching users with query: %s", query)

	users, err := u.UserCrud.SearchUsersByName(query)
	if err != nil {
		u.logger.Errorf("Failed to search users: %v", err)
		c.JSON(http.StatusInternalServerError, "Failed to search users")
		return
	}

	userDtos := make([]dto.UserDto, 0, len(users))
	for _, user := range users {
		dto := dto.UserDto{}
		userDtos = append(userDtos, dto.FromModel(&user))
	}

	c.JSON(http.StatusOK, userDtos)
}

// // GetUserByOIB godoc
// //
// //	@Summary		get user with oib
// //	@Description	get a user with oib
// //	@Tags			user
// //	@Produce		json
// //	@Success		200	{object}	dto.UserDto
// //	@Failure		400
// //	@Failure		404
// //	@Failure		500
// //	@Param			oib	path	string	true	"user oib"
// //	@Router			/user/oib/{oib} [get]
// func (u *UserController) getUserByOib(c *gin.Context) {
// 	oib := c.Param("oib")
// 	if oib == "" {
// 		u.logger.Errorf("OIB parameter is empty")
// 		c.AbortWithError(http.StatusBadRequest, errors.New("OIB parameter is required"))
// 		return
// 	}

// 	user, err := u.UserCrud.GetUserByOIB(oib)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			u.logger.Errorf("User with OIB = %s not found", oib)
// 			c.AbortWithError(http.StatusNotFound, err)
// 			return
// 		}

// 		u.logger.Errorf("Failed to get user with OIB = %s", oib)
// 		c.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	dto := dto.UserDto{}
// 	c.JSON(http.StatusOK, dto.FromModel(user))
// }
