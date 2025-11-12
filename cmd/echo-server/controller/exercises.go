package controller

import (
	"errors"
	"net/http"
	"strconv"

	"Dedenruslan19/med-project/cmd/echo-server/middleware"
	errs "Dedenruslan19/med-project/service/errors"
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
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		ec.logger.Error("invalid or missing user_id in token")
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
		if errors.Is(err, errs.ErrInvalidAuthor) {
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}
		if errors.Is(err, errs.ErrExerciseNotFound) {
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

func (ec *ExerciseController) GetExercisesByWorkoutID(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}
	workoutID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	exercisesList, err := ec.service.GetExercisesByWorkoutID(userID, workoutID)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidAuthor) {
			return c.JSON(http.StatusForbidden, ErrUnauthorized)
		}
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	var responseExercises []map[string]interface{}
	for _, ex := range exercisesList {
		responseExercises = append(responseExercises, map[string]interface{}{
			"id":            ex.ID,
			"exercise_name": ex.Name,
			"sets":          ex.Sets,
			"reps":          ex.Reps,
			"equipment":     ex.Equipment,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "success",
		"exercises": responseExercises,
	})
}

func (ec *ExerciseController) UpdateExercise(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		ec.logger.Error("invalid or missing user_id in token")
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	exerciseIDParam := c.Param("id")
	exerciseID, err := strconv.ParseInt(exerciseIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidParams)
	}

	var input exercises.ExerciseInput
	if bindErr := c.Bind(&input); bindErr != nil {
		ec.logger.Error("failed to bind JSON", slog.Any("error", bindErr))
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if validErr := ec.validate.Struct(input); validErr != nil {
		validationErrors := validErr.(validator.ValidationErrors)
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "validation failed",
			"errors":  errorsMap,
		})
	}

	updatedExercise, err := ec.service.UpdateExercise(userID, exerciseID, input)
	if err != nil {
		if errors.Is(err, errs.ErrExerciseNotFound) {
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		}
		if errors.Is(err, errs.ErrInvalidAuthor) {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "you are not the owner of this exercise"})
		}
		ec.logger.Error("failed to update exercise", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "exercise updated successfully",
		"data":    updatedExercise,
	})
}

func (ec *ExerciseController) DeleteExercise(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok || userID == 0 {
		ec.logger.Error("invalid or missing user_id in token")
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	exerciseIDParam := c.Param("id")
	exerciseID, err := strconv.ParseInt(exerciseIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidParams)
	}

	_, err = ec.service.DeleteExercise(userID, exerciseID)
	if err != nil {
		if errors.Is(err, errs.ErrExerciseNotFound) {
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		}
		if errors.Is(err, errs.ErrInvalidAuthor) {
			return c.JSON(http.StatusForbidden, map[string]string{"message": err.Error()})
		}
		ec.logger.Error("failed to delete exercise", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "exercise deleted successfully",
	})
}
