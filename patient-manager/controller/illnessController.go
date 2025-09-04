package controller

import (
	"PatientManager/app"
	"PatientManager/dto"
	"PatientManager/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IllnessController struct {
	illnessService service.IIllnessService
	logger         *zap.SugaredLogger
}

func NewIllnessController() *IllnessController {
	var controller *IllnessController
	app.Invoke(func(illnessService service.IIllnessService, logger *zap.SugaredLogger) {
		controller = &IllnessController{
			illnessService: illnessService,
			logger:         logger,
		}
	})
	return controller
}

func (ic *IllnessController) RegisterEndpoints(router *gin.RouterGroup) {
	illnessRoutes := router.Group("/illnesses")
	{
		illnessRoutes.POST("", ic.create)
		illnessRoutes.GET("/record/:recordUuid", ic.getAllForRecord)
		illnessRoutes.PUT("/:uuid", ic.update)
		illnessRoutes.DELETE("/:uuid", ic.delete)
	}
}

// create godoc
// @Summary		Create illness
// @Description	Creates a new illness for a patient's medical record
// @Tags			illnesses
// @Accept			json
// @Produce		json
// @Param			model	body		dto.CreateIllnessDto	true	"New Illness Data"
// @Success		201		{object}	model.Illness
// @Failure		400		{object}	gin.H
// @Failure		500		{object}	gin.H
// @Router			/illnesses [post]
func (ic *IllnessController) create(c *gin.Context) {
	var createDto dto.CreateIllnessDto
	if err := c.ShouldBindJSON(&createDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	illnessModel := createDto.ToModel()
	createdIllness, err := ic.illnessService.Create(illnessModel, createDto.MedicalRecordUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create illness"})
		return
	}

	c.JSON(http.StatusCreated, createdIllness)
}

// getAllForRecord godoc
// @Summary		Get all illnesses for a record
// @Description	Retrieves a list of illnesses for a specific medical record
// @Tags			illnesses
// @Produce		json
// @Param			recordUuid	path		string	true	"Medical Record UUID"
// @Success		200			{array}		dto.IllnessListDto
// @Failure		400			{object}	gin.H
// @Failure		500			{object}	gin.H
// @Router			/illnesses/record/{recordUuid} [get]
func (ic *IllnessController) getAllForRecord(c *gin.Context) {
	recordUuid, err := uuid.Parse(c.Param("recordUuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	illnesses, err := ic.illnessService.GetAllForRecord(recordUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve illnesses"})
		return
	}

	var responseDtos []*dto.IllnessListDto
	for _, illness := range illnesses {
		dto := (&dto.IllnessListDto{}).FromModel(&illness)
		responseDtos = append(responseDtos, dto)
	}
	c.JSON(http.StatusOK, responseDtos)
}

// update godoc
// @Summary		Update illness
// @Description	Updates an existing illness by its UUID
// @Tags			illnesses
// @Accept			json
// @Produce		json
// @Param			uuid	path		string					true	"Illness UUID"
// @Param			model	body		dto.UpdateIllnessDto	true	"Updated Illness Data"
// @Success		200		{object}	model.Illness
// @Failure		400		{object}	gin.H
// @Failure		500		{object}	gin.H
// @Router			/illnesses/{uuid} [put]
func (ic *IllnessController) update(c *gin.Context) {
	illnessUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var updateDto dto.UpdateIllnessDto
	if err := c.ShouldBindJSON(&updateDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedIllness, err := ic.illnessService.Update(illnessUuid, updateDto.ToModel())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update illness"})
		return
	}
	c.JSON(http.StatusOK, updatedIllness)
}

// delete godoc
// @Summary		Delete illness
// @Description	Deletes an illness by its UUID
// @Tags			illnesses
// @Param			uuid	path	string	true	"Illness UUID"
// @Success		204
// @Failure		400	{object}	gin.H
// @Failure		500	{object}	gin.H
// @Router			/illnesses/{uuid} [delete]
func (ic *IllnessController) delete(c *gin.Context) {
	illnessUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := ic.illnessService.Delete(illnessUuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete illness"})
		return
	}
	c.Status(http.StatusNoContent)
}
