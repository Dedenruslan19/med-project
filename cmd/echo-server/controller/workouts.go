package controller

import (
	"net/http"
	"strconv"

	"Dedenruslan19/med-project/service/workouts"

	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type WorkoutController struct {
	service  workouts.Service
	validate *validator.Validate
	logger   *slog.Logger
}

func NewWorkoutController(service workouts.Service, logger *slog.Logger) *WorkoutController {
	return &WorkoutController{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

func (wc *WorkoutController) GetAllWorkouts(c echo.Context) error {
	data, err := wc.service.GetAllWorkouts()
	if err != nil {
		wc.logger.Error("Failed to get all workouts", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "success",
		"workouts": data,
	})
}

func (wc *WorkoutController) CreateWorkout(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	var input workouts.WorkoutInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := wc.validate.Struct(input); err != nil {
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

	newWorkout, err := wc.service.CreateWorkout(userID, input)
	if err != nil {
		wc.logger.Error("Failed to create workout", slog.Any("error", err), slog.Int64("user_id", userID))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "workout created successfully",
		"data":    newWorkout,
	})
}

func (wc *WorkoutController) GetWorkoutByID(c echo.Context) error {
	workoutIDParam := c.Param("id")
	workoutID, err := strconv.ParseInt(workoutIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidParams)
	}

	workout, err := wc.service.GetWorkoutByID(workoutID)
	if err != nil {
		if err == workouts.ErrWorkoutNotFound {
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		}
		wc.logger.Error("Failed to get workout by ID", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    workout,
	})
}

func (wc *WorkoutController) UpdateWorkout(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	workoutIDParam := c.Param("id")
	workoutID, err := strconv.ParseInt(workoutIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidParams)
	}

	var input workouts.WorkoutInput
	if bindErr := c.Bind(&input); bindErr != nil {
		return c.JSON(http.StatusBadRequest, workouts.ErrInvalidInput)
	}

	if validateErr := wc.validate.Struct(input); validateErr != nil {
		validationErrors := validateErr.(validator.ValidationErrors)
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "validation failed",
			"errors":  errorsMap,
		})
	}

	updatedWorkout, err := wc.service.UpdateWorkout(userID, workoutID, input)
	if err != nil {
		switch err {
		case workouts.ErrWorkoutNotFound:
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		case workouts.ErrInvalidAuthor:
			return c.JSON(http.StatusForbidden, ErrUnauthorized)
		default:
			wc.logger.Error("Failed to update workout", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "workout updated successfully",
		"workout": updatedWorkout,
	})
}

func (wc *WorkoutController) DeleteWorkout(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	workoutIDParam := c.Param("id")
	workoutID, err := strconv.ParseInt(workoutIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidParams)
	}

	err = wc.service.DeleteWorkout(userID, workoutID)
	if err != nil {
		switch err {
		case workouts.ErrInvalidAuthor:
			return c.JSON(http.StatusForbidden, ErrUnauthorized)
		case workouts.ErrWorkoutNotFound:
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		default:
			wc.logger.Error("Failed to delete workout", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "workout deleted successfully",
	})
}
