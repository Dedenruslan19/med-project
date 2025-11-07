package controller

import (
	"Dedenruslan19/med-project/service/doctors"
	"Dedenruslan19/med-project/util"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DoctorController struct {
	service  doctors.Service
	validate *validator.Validate
	logger   *slog.Logger
}

func NewDoctorController(service doctors.Service, logger *slog.Logger) *DoctorController {
	return &DoctorController{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

type DoctorRegisterRequest struct {
	FullName       string `json:"full_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	Specialization string `json:"specialization" validate:"required"`
}

func (dc *DoctorController) Register(c echo.Context) error {
	var req DoctorRegisterRequest
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
	doctor, err := dc.service.Register(req.FullName, req.Email, req.Password, req.Specialization)
	if err != nil {
		dc.logger.Error("Failed to register doctor",
			slog.Any("error", err),
			slog.String("email", req.Email),
		)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Doctor registered successfully",
		"data": map[string]interface{}{
			"id":             doctor.ID,
			"full_name":      doctor.FullName,
			"email":          doctor.Email,
			"specialization": doctor.Specialization,
		},
	})
}

type DoctorLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (dc *DoctorController) Login(c echo.Context) error {
	var req DoctorLoginRequest
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
	doctor, err := dc.service.Login(req.Email, req.Password)
	if err != nil {
		dc.logger.Error("Failed to login doctor",
			slog.Any("error", err),
			slog.String("email", req.Email),
		)
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid email or password",
		})
	}

	// Generate JWT token with role "doctor"
	token, err := util.GenerateJWT(doctor.ID, "doctor")
	if err != nil {
		dc.logger.Error("Failed to generate JWT token",
			slog.Any("error", err),
			slog.Int64("doctor_id", doctor.ID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (dc *DoctorController) GetAllDoctors(c echo.Context) error {
	doctors, err := dc.service.GetAll()
	if err != nil {
		dc.logger.Error("Failed to get all doctors", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get doctors",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Doctors retrieved successfully",
		"data":    doctors,
	})
}
