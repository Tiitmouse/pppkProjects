package controller

import (
	"PatientManager/app"
	"PatientManager/dto"
	"PatientManager/service"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CheckupController struct {
	checkupService service.ICheckupService
	logger         *zap.SugaredLogger
}

func NewCheckupController() *CheckupController {
	var controller *CheckupController

	app.Invoke(func(checkupService service.ICheckupService, logger *zap.SugaredLogger) {
		controller = &CheckupController{
			checkupService: checkupService,
			logger:         logger,
		}
	})

	return controller
}

func (cc *CheckupController) RegisterEndpoints(router *gin.RouterGroup) {
	checkupRoutes := router.Group("/checkup")
	{
		checkupRoutes.POST("", cc.create)
		checkupRoutes.PUT("/:uuid", cc.update)
		checkupRoutes.DELETE("/:uuid", cc.delete)
	}
}

// create godoc
//
//	@Summary		Create a new checkup
//	@Description	Creates a new checkup associated with a medical record.
//	@Tags			checkup
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	dto.CheckupDto
//	@Failure		400
//	@Failure		500
//	@Param			model	body	dto.CreateCheckupDto	true	"Data for creating a new checkup"
//	@Router			/checkup [post]
func (cc *CheckupController) create(c *gin.Context) {
	var createDto dto.CreateCheckupDto
	if err := c.ShouldBindJSON(&createDto); err != nil {
		cc.logger.Errorf("Error binding JSON for create checkup: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	checkupModel, err := createDto.ToModel()
	if err != nil {
		cc.logger.Errorf("Error converting DTO to model for create checkup: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	createdCheckup, err := cc.checkupService.Create(checkupModel, createDto.MedicalRecordUuid)
	if err != nil {
		cc.logger.Errorf("Failed to create checkup: %+v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var responseDto dto.CheckupDto
	c.JSON(http.StatusCreated, responseDto.FromModel(createdCheckup))
}

// update godoc
//
//	@Summary		Update an existing checkup
//	@Description	Updates the details of a specific checkup by its UUID.
//	@Tags			checkup
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.CheckupDto
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			uuid	path	string			true	"UUID of the checkup to be updated"
//	@Param			model	body	dto.CheckupDto	true	"Data for updating the checkup"
//	@Router			/checkup/{uuid} [put]
func (cc *CheckupController) update(c *gin.Context) {
	checkupUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		cc.logger.Errorf("Error parsing UUID '%s': %v", c.Param("uuid"), err)
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid UUID format"))
		return
	}

	var updateDto dto.CheckupDto
	if err := c.ShouldBindJSON(&updateDto); err != nil {
		cc.logger.Errorf("Error binding JSON for update checkup: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updateData, err := updateDto.ToModel()
	if err != nil {
		cc.logger.Errorf("Error converting DTO to model for update checkup: %v", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	updatedCheckup, err := cc.checkupService.Update(checkupUuid, updateData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cc.logger.Warnf("Checkup with UUID %s not found for update", checkupUuid)
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		cc.logger.Errorf("Failed to update checkup with UUID %s: %+v", checkupUuid, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var responseDto dto.CheckupDto
	c.JSON(http.StatusOK, responseDto.FromModel(updatedCheckup))
}

// delete godoc
//
//	@Summary		Delete a checkup
//	@Description	Deletes a checkup by its UUID.
//	@Tags			checkup
//	@Success		204
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			uuid	path	string	true	"UUID of the checkup to be deleted"
//	@Router			/checkup/{uuid} [delete]
func (cc *CheckupController) delete(c *gin.Context) {
	checkupUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		cc.logger.Errorf("Error parsing UUID '%s': %v", c.Param("uuid"), err)
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid UUID format"))
		return
	}

	err = cc.checkupService.Delete(checkupUuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cc.logger.Warnf("Checkup with UUID %s not found for deletion", checkupUuid)
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		cc.logger.Errorf("Failed to delete checkup with UUID %s: %+v", checkupUuid, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
