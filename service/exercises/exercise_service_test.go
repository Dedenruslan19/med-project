package exercises_test

import (
	"Dedenruslan19/med-project/service/exercises"
	"Dedenruslan19/med-project/service/workouts"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetByWorkoutID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := exercises.NewMockExerciseRepo(ctrl)
	mockWorkoutService := workouts.NewMockWorkoutRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create workout service with mock repo
	workoutSvc := workouts.NewService(logger, mockWorkoutService, nil)
	service := exercises.NewService(logger, mockRepo, workoutSvc, nil)

	// Mock workout validation
	expectedWorkout := &workouts.Workout{
		ID:     1,
		UserID: 1,
		Name:   "Test Workout",
		Goals:  "Test Goals",
	}

	mockWorkoutService.EXPECT().
		GetByID(int64(1)).
		Return(expectedWorkout, nil).
		Times(1)

	expectedExercises := []exercises.Exercise{
		{
			ID:        1,
			WorkoutID: 1,
			Name:      "Push-ups",
			Sets:      "3",
			Reps:      "15",
			Equipment: "None",
		},
		{
			ID:        2,
			WorkoutID: 1,
			Name:      "Squats",
			Sets:      "3",
			Reps:      "20",
			Equipment: "None",
		},
	}

	mockRepo.EXPECT().
		GetByWorkoutID(int64(1)).
		Return(expectedExercises, nil).
		Times(1)

	result, err := service.GetExercisesByWorkoutID(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedExercises[0].Name, result[0].Name)
	assert.Equal(t, expectedExercises[1].Name, result[1].Name)
}

func TestGetByWorkoutID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := exercises.NewMockExerciseRepo(ctrl)
	mockWorkoutService := workouts.NewMockWorkoutRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	workoutSvc := workouts.NewService(logger, mockWorkoutService, nil)
	service := exercises.NewService(logger, mockRepo, workoutSvc, nil)

	// Mock workout not found
	mockWorkoutService.EXPECT().
		GetByID(int64(999)).
		Return(nil, errors.New("workout not found")).
		Times(1)

	result, err := service.GetExercisesByWorkoutID(1, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
