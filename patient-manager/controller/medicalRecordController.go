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

type MedicalRecordController struct {
	MedicalRecordService service.IMedicalRecordService
	logger               *zap.SugaredLogger
}

func NewMedicalRecordController() *MedicalRecordController {
	var controller *MedicalRecordController

	app.Invoke(func(medicalRecordService service.IMedicalRecordService, logger *zap.SugaredLogger) {
		controller = &MedicalRecordController{
			MedicalRecordService: medicalRecordService,
			logger:               logger,
		}
	})

	return controller
}

func (m *MedicalRecordController) RegisterEndpoints(router *gin.RouterGroup) {
	mr := router.Group("/medical-record")
	{
		mr.GET("/:patientOib", m.get)
		mr.PUT("/:uuid", m.update)
	}
}

// MedicalRecordExample godoc
//
//	@Summary		Get medical record by patient OIB
//	@Description	get a medical record by patient OIB
//	@Tags			medical-record
//	@Produce		json
//	@Success		200	{object}	dto.MedicalRecordDto
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			patientOib	path	string	true	"Patient OIB"
//	@Router			/medical-record/{patientOib} [get]
func (m *MedicalRecordController) get(c *gin.Context) {
	patientOib := c.Param("patientOib")
	if patientOib == "" {
		m.logger.Error("patientOib path parameter is empty")
		c.AbortWithError(http.StatusBadRequest, errors.New("patientOib cannot be empty"))
		return
	}

	record, err := m.MedicalRecordService.Read(patientOib)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			m.logger.Errorf("Medical record for patient with OIB = %s not found", patientOib)
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		m.logger.Errorf("Failed to get medical record for patient with OIB = %s", patientOib)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var dto dto.MedicalRecordDto
	c.JSON(http.StatusOK, dto.FromModel(record))
}

// MedicalRecordExample godoc
//
//	@Summary	Update medical record with new data
//	@Tags		medical-record
//	@Produce	json
//	@Success	200	{object}	dto.MedicalRecordDto
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Param		uuid	path	string					true	"uuid of medical record to be updated"
//	@Param		model	body	dto.MedicalRecordDto	true	"Data for updating medical record"
//	@Router		/medical-record/{uuid} [put]
func (m *MedicalRecordController) update(c *gin.Context) {
	recordUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		m.logger.Errorf("Error parsing UUID = %s", c.Param("uuid"))
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	data, err := m.MedicalRecordService.Read(recordUuid.String())
	data.DoctorID = 1 // TODO: read new doctor id

	record, err := m.MedicalRecordService.Update(recordUuid, data)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			m.logger.Errorf("Medical record with uuid = %s not found for update", recordUuid)
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		m.logger.Errorf("Failed to update medical record: %+v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var responseDto dto.MedicalRecordDto
	c.JSON(http.StatusOK, responseDto.FromModel(record))
}
