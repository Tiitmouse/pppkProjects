package controller

import (
	"PatientManager/dto"
	"PatientManager/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PatientController struct {
	patientService service.PatientService
}

func NewPatientController(patientService service.PatientService) *PatientController {
	return &PatientController{patientService: patientService}
}

// GetAllPatients godoc
// @Summary      List all patients
// @Description  get all patients
// @Tags         patients
// @Produce      json
// @Success      200  {array}   dto.PatientDto
// @Failure      500  {object}  gin.H
// @Router       /patients [get]
func (c *PatientController) GetAllPatients(ctx *gin.Context) {
	patients, err := c.patientService.GetAllPatients()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve patients"})
		return
	}
	ctx.JSON(http.StatusOK, patients)
}

// GetPatientById godoc
// @Summary      Get a patient by ID
// @Description  get patient by ID
// @Tags         patients
// @Produce      json
// @Param        id   path      int  true  "Patient ID"
// @Success      200  {object}  dto.PatientDto
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Router       /patients/{id} [get]
func (c *PatientController) GetPatientById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	patient, err := c.patientService.GetPatientById(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	ctx.JSON(http.StatusOK, patient)
}

// CreatePatient godoc
// @Summary      Create a new patient
// @Description  add a new patient to the database
// @Tags         patients
// @Accept       json
// @Produce      json
// @Param        patient  body      dto.NewPatientDto  true  "New Patient"
// @Success      201      {object}  dto.PatientDto
// @Failure      400      {object}  gin.H
// @Failure      500      {object}  gin.H
// @Router       /patients [post]
func (c *PatientController) CreatePatient(ctx *gin.Context) {
	var newPatient dto.NewPatientDto
	if err := ctx.ShouldBindJSON(&newPatient); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPatient, err := c.patientService.CreatePatient(newPatient)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	ctx.JSON(http.StatusCreated, createdPatient)
}

// UpdatePatient godoc
// @Summary      Update an existing patient
// @Description  update patient details by ID
// @Tags         patients
// @Accept       json
// @Produce      json
// @Param        id       path      int            true  "Patient ID"
// @Param        patient  body      dto.PatientDto  true  "Patient Data"
// @Success      200      {object}  dto.PatientDto
// @Failure      400      {object}  gin.H
// @Failure      500      {object}  gin.H
// @Router       /patients/{id} [put]
func (c *PatientController) UpdatePatient(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	var patientDto dto.PatientDto
	if err := ctx.ShouldBindJSON(&patientDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPatient, err := c.patientService.UpdatePatient(uint(id), patientDto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	ctx.JSON(http.StatusOK, updatedPatient)
}

// DeletePatient godoc
// @Summary      Delete a patient
// @Description  delete a patient by ID
// @Tags         patients
// @Param        id   path      int  true  "Patient ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      500  {object}  gin.H
// @Router       /patients/{id} [delete]
func (c *PatientController) DeletePatient(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	if err := c.patientService.DeletePatient(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
