package users_test

import (
	"Dedenruslan19/med-project/service/users"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := users.NewMockUserRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := users.NewService(logger, mockRepo, nil)

	expectedUser := users.User{
		ID:       1,
		Email:    "john@example.com",
		Height:   175.0,
		Weight:   70.0,
		Password: "hashedpassword",
	}

	mockRepo.EXPECT().
		FindByID(int64(1)).
		Return(expectedUser, nil).
		Times(1)

	result, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
}

func TestGetUserByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := users.NewMockUserRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := users.NewService(logger, mockRepo, nil)

	mockRepo.EXPECT().
		FindByID(int64(999)).
		Return(users.User{}, errors.New("user not found")).
		Times(1)

	result, err := service.GetUserByID(999)

	assert.Error(t, err)
	assert.Equal(t, int64(0), result.ID)
}
