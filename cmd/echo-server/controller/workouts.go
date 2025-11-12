package controller

import (
	"net/http"
	"strconv"

	"Dedenruslan19/med-project/cmd/echo-server/middleware"
	errs "Dedenruslan19/med-project/service/errors"
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
func (wc *WorkoutController) PreviewWorkout(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	var req workouts.PreviewWorkoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := wc.validate.Struct(req); err != nil {
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

	previewWorkout, err := wc.service.PreviewWorkout(userID, req)
	if err != nil {
		wc.logger.Error("Failed to preview workout", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	// Format response: {workout_name: {object}, exercises: [array]}
	response := map[string]interface{}{
		"workout_name": map[string]interface{}{
			"workout_name": previewWorkout.Workout.Name,
			"goals":        previewWorkout.Workout.Goals,
			"id":           previewWorkout.Workout.ID,
			"user_id":      previewWorkout.Workout.UserID,
		},
		"exercises": previewWorkout.Exercises,
	}

	return c.JSON(http.StatusOK, response)
}
func (wc *WorkoutController) CreateWorkout(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	var req workouts.SaveWorkoutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := wc.validate.Struct(req); err != nil {
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

	newWorkout, err := wc.service.CreateWorkout(userID, req)
	if err != nil {
		wc.logger.Error("Failed to create workout", slog.Any("error", err), slog.Int64("user_id", userID))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	// Format response sama seperti preview: {workout_name: {object}, exercises: [array]}
	response := map[string]interface{}{
		"workout_name": map[string]interface{}{
			"workout_name": newWorkout.Workout.Name,
			"goals":        newWorkout.Workout.Goals,
			"id":           newWorkout.Workout.ID,
			"user_id":      newWorkout.Workout.UserID,
		},
		"exercises": newWorkout.Exercises,
	}

	return c.JSON(http.StatusCreated, response)
}

func (wc *WorkoutController) GetWorkoutByID(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	workoutIDParam := c.Param("id")
	workoutID, err := strconv.ParseInt(workoutIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidParams)
	}

	workout, err := wc.service.GetWorkoutByID(userID, workoutID)
	if err != nil {
		switch err {
		case errs.ErrWorkoutNotFound:
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		case errs.ErrInvalidAuthor:
			return c.JSON(http.StatusForbidden, ErrUnauthorized)
		default:
			wc.logger.Error("Failed to get workout by ID", slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		}
	}

	// Response sama seperti preview
	response := map[string]interface{}{
		"workout_name": map[string]interface{}{
			"workout_name": workout.Name,
			"goals":        workout.Goals,
			"id":           workout.ID,
			"user_id":      workout.UserID,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (wc *WorkoutController) DeleteWorkout(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
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
		case errs.ErrInvalidAuthor:
			return c.JSON(http.StatusForbidden, ErrUnauthorized)
		case errs.ErrWorkoutNotFound:
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
