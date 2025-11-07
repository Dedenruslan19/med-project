package controller

import (
	"Dedenruslan19/med-project/service/appointments"
	"Dedenruslan19/med-project/service/billings"
	"Dedenruslan19/med-project/service/diagnoses"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DiagnoseController struct {
	service            diagnoses.Service
	appointmentService appointments.Service
	billingService     billings.Service
	validate           *validator.Validate
	logger             *slog.Logger
}

func NewDiagnoseController(service diagnoses.Service, appointmentService appointments.Service, billingService billings.Service, logger *slog.Logger) *DiagnoseController {
	return &DiagnoseController{
		service:            service,
		appointmentService: appointmentService,
		billingService:     billingService,
		validate:           validator.New(),
		logger:             logger,
	}
}

type CreateDiagnoseRequest struct {
	AppointmentID         int64  `json:"appointment_id" validate:"required"`
	DoctorID              int64  `json:"doctor_id" validate:"required"`
	Notes                 string `json:"notes" validate:"required"`
	PrescribedMedications string `json:"prescribed_medications"`
}

type UpdateDiagnoseRequest struct {
	Notes                 string `json:"notes"`
	PrescribedMedications string `json:"prescribed_medications"`
}

func (dc *DiagnoseController) CreateDiagnose(c echo.Context) error {
	var req CreateDiagnoseRequest
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

	doctorIDFromToken, ok := c.Get("user_id").(int64)
	if !ok {
		dc.logger.Error("Failed to get doctor ID from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized: invalid token",
		})
	}

	appointment, err := dc.appointmentService.GetByID(req.AppointmentID)
	if err != nil {
		dc.logger.Error("Failed to get appointment",
			slog.Any("error", err),
			slog.Int64("appointment_id", req.AppointmentID),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Appointment not found",
		})
	}

	if appointment.DoctorID != doctorIDFromToken {
		dc.logger.Warn("Doctor attempting to create diagnose for another doctor's appointment",
			slog.Int64("token_doctor_id", doctorIDFromToken),
			slog.Int64("appointment_doctor_id", appointment.DoctorID),
			slog.Int64("appointment_id", req.AppointmentID),
		)
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "You are not authorized to create diagnose for this appointment",
		})
	}

	diagnose := &diagnoses.Diagnose{
		AppointmentID:         req.AppointmentID,
		DoctorID:              req.DoctorID,
		Notes:                 req.Notes,
		PrescribedMedications: req.PrescribedMedications,
	}

	id, err := dc.service.Create(diagnose)
	if err != nil {
		dc.logger.Error("Failed to create diagnose",
			slog.Any("error", err),
			slog.Int64("appointment_id", req.AppointmentID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create diagnose",
		})
	}

	diagnose.ID = id

	totalAmount := dc.service.CalculateTotalAmount(diagnose)

	billing := &billings.Billing{
		AppointmentID: req.AppointmentID,
		TotalAmount:   totalAmount,
		PaymentStatus: "unpaid",
	}

	billingID, err := dc.billingService.Create(billing)
	if err != nil {
		dc.logger.Error("Failed to create billing after diagnose",
			slog.Any("error", err),
			slog.Int64("appointment_id", req.AppointmentID),
		)
	} else {
		billing.ID = billingID
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Diagnose created successfully",
		"diagnose": diagnose,
	})
}

func (dc *DiagnoseController) GetDiagnoseByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid diagnose ID",
		})
	}

	diagnose, err := dc.service.GetByID(id)
	if err != nil {
		dc.logger.Error("Failed to get diagnose",
			slog.Any("error", err),
			slog.Int64("diagnose_id", id),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "diagnose not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "diagnose retrieved successfully",
		"data":    diagnose,
	})
}

func (dc *DiagnoseController) GetDiagnoseByAppointmentID(c echo.Context) error {
	idParam := c.Param("appointment_id")
	appointmentID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment ID",
		})
	}

	diagnose, err := dc.service.GetByAppointmentID(appointmentID)
	if err != nil {
		dc.logger.Error("Failed to get diagnose by appointment ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", appointmentID),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "diagnose not found for this appointment",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "diagnose retrieved successfully",
		"data":    diagnose,
	})
}

func (dc *DiagnoseController) UpdateDiagnose(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid diagnose ID",
		})
	}

	// Get existing diagnose
	diagnose, err := dc.service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "diagnose not found",
		})
	}

	var req UpdateDiagnoseRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Notes != "" {
		diagnose.Notes = req.Notes
	}
	if req.PrescribedMedications != "" {
		diagnose.PrescribedMedications = req.PrescribedMedications
	}

	err = dc.service.Update(diagnose)
	if err != nil {
		dc.logger.Error("Failed to update diagnose",
			slog.Any("error", err),
			slog.Int64("diagnose_id", id),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update diagnose",
		})
	}

	newTotalAmount := dc.service.CalculateTotalAmount(diagnose)
	billing, err := dc.billingService.GetByAppointmentID(diagnose.AppointmentID)
	if err == nil && billing != nil {
		dc.logger.Info("Billing amount should be updated",
			slog.Float64("new_amount", newTotalAmount),
			slog.Int64("billing_id", billing.ID),
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "diagnose updated successfully",
		"data":         diagnose,
		"total_amount": newTotalAmount,
	})
}
