package controller

import (
	"Dedenruslan19/med-project/service/appointments"
	"log/slog"

	"github.com/go-playground/validator/v10"
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

// func (ac *AppointmentController) CreateAppointment(c echo.Context) error {

// }
