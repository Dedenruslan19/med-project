package controller

import (
	"Dedenruslan19/med-project/service/billings"
	"Dedenruslan19/med-project/service/diagnoses"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DiagnosisController struct {
	service        diagnoses.Service
	billingService billings.Service
	validate       *validator.Validate
	logger         *slog.Logger
}

func NewDiagnosisController(service diagnoses.Service, billingService billings.Service, logger *slog.Logger) *DiagnosisController {
	return &DiagnosisController{
		service:        service,
		billingService: billingService,
		validate:       validator.New(),
		logger:         logger,
	}
}

type CreateDiagnosisRequest struct {
	AppointmentID         int64  `json:"appointment_id" validate:"required"`
	DoctorID              int64  `json:"doctor_id" validate:"required"`
	Notes                 string `json:"notes" validate:"required"`
	PrescribedMedications string `json:"prescribed_medications"`
}

func (dc *DiagnosisController) CreateDiagnosis(c echo.Context) error {
	var req CreateDiagnosisRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := dc.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	diagnosis := &diagnoses.Diagnosis{
		AppointmentID:         req.AppointmentID,
		DoctorID:              req.DoctorID,
		Notes:                 req.Notes,
		PrescribedMedications: req.PrescribedMedications,
	}

	id, err := dc.service.Create(diagnosis)
	if err != nil {
		dc.logger.Error("Failed to create diagnosis",
			slog.Any("error", err),
			slog.Int64("appointment_id", req.AppointmentID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create diagnosis",
		})
	}

	diagnosis.ID = id

	// Hitung total amount dan buat billing otomatis
	totalAmount := dc.service.CalculateTotalAmount(diagnosis)

	billing := &billings.Billing{
		AppointmentID: req.AppointmentID,
		TotalAmount:   totalAmount,
		PaymentStatus: "unpaid",
	}

	billingID, err := dc.billingService.Create(billing)
	if err != nil {
		dc.logger.Error("Failed to create billing after diagnosis",
			slog.Any("error", err),
			slog.Int64("appointment_id", req.AppointmentID),
		)
		// Diagnosis sudah dibuat, jadi tetap return success
	} else {
		billing.ID = billingID
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Diagnosis created successfully",
		"data": map[string]interface{}{
			"diagnosis":    diagnosis,
			"billing":      billing,
			"total_amount": totalAmount,
		},
	})
}

func (dc *DiagnosisController) GetDiagnosisByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid diagnosis ID",
		})
	}

	diagnosis, err := dc.service.GetByID(id)
	if err != nil {
		dc.logger.Error("Failed to get diagnosis",
			slog.Any("error", err),
			slog.Int64("diagnosis_id", id),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Diagnosis not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Diagnosis retrieved successfully",
		"data":    diagnosis,
	})
}

func (dc *DiagnosisController) GetDiagnosisByAppointmentID(c echo.Context) error {
	idParam := c.Param("appointment_id")
	appointmentID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment ID",
		})
	}

	diagnosis, err := dc.service.GetByAppointmentID(appointmentID)
	if err != nil {
		dc.logger.Error("Failed to get diagnosis by appointment ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", appointmentID),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Diagnosis not found for this appointment",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Diagnosis retrieved successfully",
		"data":    diagnosis,
	})
}

type UpdateDiagnosisRequest struct {
	Notes                 string `json:"notes"`
	PrescribedMedications string `json:"prescribed_medications"`
}

func (dc *DiagnosisController) UpdateDiagnosis(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid diagnosis ID",
		})
	}

	// Get existing diagnosis
	diagnosis, err := dc.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Diagnosis not found",
		})
	}

	var req UpdateDiagnosisRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if req.Notes != "" {
		diagnosis.Notes = req.Notes
	}
	if req.PrescribedMedications != "" {
		diagnosis.PrescribedMedications = req.PrescribedMedications
	}

	err = dc.service.Update(diagnosis)
	if err != nil {
		dc.logger.Error("Failed to update diagnosis",
			slog.Any("error", err),
			slog.Int64("diagnosis_id", id),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update diagnosis",
		})
	}

	// Update billing amount jika ada perubahan obat
	newTotalAmount := dc.service.CalculateTotalAmount(diagnosis)
	billing, err := dc.billingService.GetByAppointmentID(diagnosis.AppointmentID)
	if err == nil && billing != nil {
		dc.logger.Info("Billing amount should be updated",
			slog.Float64("new_amount", newTotalAmount),
			slog.Int64("billing_id", billing.ID),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Diagnosis updated successfully",
		"data":         diagnosis,
		"total_amount": newTotalAmount,
	})
}
