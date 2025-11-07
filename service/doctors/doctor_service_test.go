package doctors_test

import (
	"Dedenruslan19/med-project/repository/doctor"
	"Dedenruslan19/med-project/service/doctors"
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

	mockRepo := doctors.NewMockDoctorRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := doctors.NewService(logger, mockRepo)

	expectedDoctor := &doctor.Doctor{
		ID:             1,
		FullName:       "Dr. John Smith",
		Email:          "john.smith@doctor.com",
		Specialization: "Cardiologist",
		IsAvailable:    true,
	}

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(expectedDoctor, nil).
		Times(1)

	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedDoctor.ID, result.ID)
	assert.Equal(t, expectedDoctor.FullName, result.FullName)
	assert.Equal(t, expectedDoctor.Specialization, result.Specialization)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := doctors.NewMockDoctorRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := doctors.NewService(logger, mockRepo)

	mockRepo.EXPECT().
		GetByID(int64(999)).
		Return(nil, errors.New("doctor not found")).
		Times(1)

	result, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
