package controller

import (
	"Dedenruslan19/med-project/cmd/echo-server/middleware"
	errs "Dedenruslan19/med-project/service/errors"
	"Dedenruslan19/med-project/service/logs"
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type LogController struct {
	service  logs.Service
	validate *validator.Validate
	logger   *slog.Logger
}

func NewLogController(service logs.Service, logger *slog.Logger) *LogController {
	return &LogController{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

func (lc *LogController) CreateLog(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		lc.logger.Warn("user_id not found or invalid in context")
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	var input logs.LogInput
	if err := c.Bind(&input); err != nil {
		lc.logger.Warn("Failed to bind input", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := lc.validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "validation failed",
			"errors":  errorsMap,
		})
	}

	newLog, err := lc.service.CreateLog(userID, input)
	if err != nil {
		lc.logger.Error("Failed to create log", slog.Any("error", err), slog.Int64("user_id", userID))
		if errors.Is(err, errs.ErrLogNotFound) {
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "log created successfully",
		"data":    newLog,
	})
}

func (lc *LogController) GetAllLogs(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		lc.logger.Warn("user_id not found or invalid in context")
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}
	logsList, err := lc.service.GetAllLogs(userID)
	if err != nil {
		lc.logger.Error("Failed to get all logs", slog.Any("error", err))
		if errors.Is(err, errs.ErrLogNotFound) {
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logs retrieved successfully",
		"data":    logsList,
	})
}
