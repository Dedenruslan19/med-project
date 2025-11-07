package workouts_test

import (
	"Dedenruslan19/med-project/service/workouts"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := workouts.NewMockWorkoutRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := workouts.NewService(logger, mockRepo, nil)

	expectedWorkout := &workouts.Workout{
		ID:     1,
		UserID: 1,
		Name:   "Morning Cardio",
		Goals:  "Cardio workout routine for better stamina",
	}

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(expectedWorkout, nil).
		Times(1)

	result, err := service.GetWorkoutByID(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedWorkout.ID, result.ID)
	assert.Equal(t, expectedWorkout.Name, result.Name)
	assert.Equal(t, expectedWorkout.Goals, result.Goals)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := workouts.NewMockWorkoutRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := workouts.NewService(logger, mockRepo, nil)

	mockRepo.EXPECT().
		GetByID(int64(999)).
		Return(nil, errors.New("workout not found")).
		Times(1)

	result, err := service.GetWorkoutByID(1, 999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
