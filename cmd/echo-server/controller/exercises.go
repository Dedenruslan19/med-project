package controller

import (
	"errors"
	"net/http"
	"strconv"

	"Dedenruslan19/med-project/service/exercises"

	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ExerciseController struct {
	service  exercises.Service
	validate *validator.Validate
	logger   *slog.Logger
}

func NewExerciseController(service exercises.Service, logger *slog.Logger) *ExerciseController {
	return &ExerciseController{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

func (ec *ExerciseController) CreateExercise(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok || userID == 0 {
		ec.logger.Error("invalid or missing user_id in token", slog.Any("user_id_value", userIDInterface))
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	var input exercises.ExerciseInput
	if err := c.Bind(&input); err != nil {
		ec.logger.Error("failed to bind JSON", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := ec.validate.Struct(input); err != nil {
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

	newExercise, err := ec.service.CreateExercise(userID, input.WorkoutID, input)
	if err != nil {
		if errors.Is(err, exercises.ErrInvalidAuthor) {
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}
		if errors.Is(err, exercises.ErrExerciseNotFound) {
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		}
		ec.logger.Error("failed to create exercise",
			slog.Any("error", err),
			slog.Int64("user_id", userID), slog.Int64("workout_id", input.WorkoutID))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "exercise created successfully",
		"data":    newExercise,
	})
}

func (ec *ExerciseController) DeleteExercise(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok || userID == 0 {
		ec.logger.Error("invalid or missing user_id in token", slog.Any("user_id_value", userIDInterface))
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	exerciseIDParam := c.Param("id")
	exerciseID, err := strconv.ParseInt(exerciseIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidParams)
	}

	_, err = ec.service.DeleteExercise(userID, exerciseID)
	if err != nil {
		if errors.Is(err, exercises.ErrExerciseNotFound) {
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		}
		if errors.Is(err, exercises.ErrInvalidAuthor) {
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}
		ec.logger.Error("failed to delete exercise", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "exercise deleted successfully",
	})
}
