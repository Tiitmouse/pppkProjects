package controller

import (
	"PatientManager/app"
	"PatientManager/dto"
	"PatientManager/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PrescriptionController struct {
	prescriptionService service.IPrescriptionService
	logger              *zap.SugaredLogger
}

func NewPrescriptionController() *PrescriptionController {
	var controller *PrescriptionController
	app.Invoke(func(prescriptionService service.IPrescriptionService, logger *zap.SugaredLogger) {
		controller = &PrescriptionController{
			prescriptionService: prescriptionService,
			logger:              logger,
		}
	})
	return controller
}

func (pc *PrescriptionController) RegisterEndpoints(router *gin.RouterGroup) {
	prescriptionRoutes := router.Group("/prescriptions")
	{
		prescriptionRoutes.POST("", pc.create)
		prescriptionRoutes.GET("/illness/:illnessId", pc.getAllForIllness)
		prescriptionRoutes.DELETE("/:uuid", pc.delete)
	}
}

// create godoc
// @Summary		Create a new prescription
// @Description	Creates a new prescription for an illness and links it to a list of medications.
// @Tags			prescriptions
// @Accept			json
// @Produce		json
// @Param			model	body		dto.CreatePrescriptionDto	true	"Data for new prescription"
// @Success		201		{object}	dto.PrescriptionListDto
// @Failure		400		{object}	gin.H
// @Failure		500		{object}	gin.H
// @Router			/prescriptions [post]
func (pc *PrescriptionController) create(c *gin.Context) {
	var createDto dto.CreatePrescriptionDto
	if err := c.ShouldBindJSON(&createDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prescriptionModel := createDto.ToModel()
	createdPrescription, err := pc.prescriptionService.Create(prescriptionModel, createDto.MedicationUuids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prescription"})
		return
	}

	responseDto := (&dto.PrescriptionListDto{}).FromModel(createdPrescription)
	c.JSON(http.StatusCreated, responseDto)
}

// getAllForIllness godoc
// @Summary		Get all prescriptions for an illness
// @Description	Retrieves a list of all prescriptions associated with a specific illness ID.
// @Tags			prescriptions
// @Produce		json
// @Param			illnessId	path		int	true	"Illness ID"
// @Success		200			{array}		dto.PrescriptionListDto
// @Failure		400			{object}	gin.H
// @Failure		500			{object}	gin.H
// @Router			/prescriptions/illness/{illnessId} [get]
func (pc *PrescriptionController) getAllForIllness(c *gin.Context) {
	illnessId, err := strconv.ParseUint(c.Param("illnessId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid illness ID"})
		return
	}

	prescriptions, err := pc.prescriptionService.GetAllForIllness(uint(illnessId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve prescriptions"})
		return
	}

	var responseDtos []*dto.PrescriptionListDto
	for _, p := range prescriptions {
		responseDtos = append(responseDtos, (&dto.PrescriptionListDto{}).FromModel(&p))
	}
	c.JSON(http.StatusOK, responseDtos)
}

// delete godoc
// @Summary		Delete a prescription
// @Description	Deletes a prescription by its UUID and disassociates its medications.
// @Tags			prescriptions
// @Param			uuid	path	string	true	"Prescription UUID"
// @Success		204
// @Failure		400	{object}	gin.H
// @Failure		500	{object}	gin.H
// @Router			/prescriptions/{uuid} [delete]
func (pc *PrescriptionController) delete(c *gin.Context) {
	prescriptionUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := pc.prescriptionService.Delete(prescriptionUuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete prescription"})
		return
	}
	c.Status(http.StatusNoContent)
}
