package logs_test

import (
	"Dedenruslan19/med-project/service/logs"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateLog_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := logs.NewMockLogRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := logs.NewService(logger, mockRepo)

	input := logs.LogInput{
		ExerciseID: 1,
		Weight:     70.5,
		RepCount:   15,
		SetCount:   3,
	}

	mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(int64(1), nil).
		Times(1)

	result, err := service.CreateLog(1, input)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, input.Weight, result.Weight)
	assert.Equal(t, input.ExerciseID, result.ExerciseID)
	assert.Equal(t, int64(1), result.UserID)
	assert.Equal(t, input.RepCount, result.RepCount)
	assert.Equal(t, input.SetCount, result.SetCount)
}

func TestCreateLog_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := logs.NewMockLogRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := logs.NewService(logger, mockRepo)

	input := logs.LogInput{
		ExerciseID: 999,
		Weight:     70.5,
		RepCount:   15,
		SetCount:   3,
	}

	mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(int64(0), errors.New("failed to create log")).
		Times(1)

	result, err := service.CreateLog(1, input)

	assert.Error(t, err)
	assert.Equal(t, int64(0), result.ID)
}
