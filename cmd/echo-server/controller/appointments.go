package controller

import (
	"Dedenruslan19/med-project/service/appointments"
	errs "Dedenruslan19/med-project/service/errors"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AppointmentController struct {
	service  appointments.Service
	validate *validator.Validate
	logger   *slog.Logger
}

func NewAppointmentController(service appointments.Service, logger *slog.Logger) *AppointmentController {
	return &AppointmentController{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

type CreateAppointmentRequest struct {
	DoctorID        int64  `json:"doctor_id" validate:"required"`
	AppointmentDate string `json:"appointment_date" validate:"required"`
	Notes           string `json:"notes"`
}

func (ac *AppointmentController) CreateAppointment(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	var req CreateAppointmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := ac.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	appointmentDate, err := time.Parse(time.RFC3339, req.AppointmentDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid date format. Use RFC3339 format",
		})
	}

	appointment := &appointments.Appointment{
		UserID:          userID,
		DoctorID:        req.DoctorID,
		AppointmentDate: appointmentDate,
		Status:          "pending",
		Notes:           req.Notes,
	}

	id, err := ac.service.Create(appointment)
	if err != nil {
		if errors.Is(err, errs.ErrDoctorBusy) {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "Doctor is not available at the requested time. Please choose another time or doctor.",
			})
		}

		ac.logger.Error("Failed to create appointment",
			slog.Any("error", err),
			slog.Int64("user_id", userID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create appointment",
		})
	}

	appointment.ID = id

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Appointment created successfully",
		"data":    appointment,
	})
}

func (ac *AppointmentController) GetAppointmentsByUser(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	appointmentList, err := ac.service.GetByUserID(userID)
	if err != nil {
		ac.logger.Error("Failed to get appointments",
			slog.Any("error", err),
			slog.Int64("user_id", userID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get appointments",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Appointments retrieved successfully",
		"data":    appointmentList,
	})
}

func (ac *AppointmentController) GetAppointmentByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment ID",
		})
	}

	appointment, err := ac.service.GetByID(id)
	if err != nil {
		ac.logger.Error("Failed to get appointment",
			slog.Any("error", err),
			slog.Int64("appointment_id", id),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Appointment not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Appointment retrieved successfully",
		"data":    appointment,
	})
}
