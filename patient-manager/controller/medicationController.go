package controller

import (
	"PatientManager/app"
	"PatientManager/dto"
	"PatientManager/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MedicationController struct {
	medicationService service.IMedicationService
	logger            *zap.SugaredLogger
}

func NewMedicationController() *MedicationController {
	var controller *MedicationController
	app.Invoke(func(medicationService service.IMedicationService, logger *zap.SugaredLogger) {
		controller = &MedicationController{
			medicationService: medicationService,
			logger:            logger,
		}
	})
	return controller
}

func (mc *MedicationController) RegisterEndpoints(router *gin.RouterGroup) {
	medicationRoutes := router.Group("/medications")
	{
		medicationRoutes.GET("", mc.getAll)
	}
}

// getAll godoc
// @Summary		List all medications
// @Description	get all medications available in the system
// @Tags			medications
// @Produce		json
// @Success		200	{array}		dto.MedicationListDto
// @Failure		500
// @Router			/medications [get]
func (mc *MedicationController) getAll(c *gin.Context) {
	medications, err := mc.medicationService.GetAll()
	if err != nil {
		mc.logger.Errorf("Failed to get all medications: %+v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var responseDtos []*dto.MedicationListDto
	for _, medication := range medications {
		dto := (&dto.MedicationListDto{}).FromModel(&medication)
		responseDtos = append(responseDtos, dto)
	}

	c.JSON(http.StatusOK, responseDtos)
}
