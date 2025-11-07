package appointments_test

import (
	"Dedenruslan19/med-project/service/appointments"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := appointments.NewMockAppointmentRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := appointments.NewService(logger, mockRepo)

	appointmentDate := time.Date(2025, 11, 10, 14, 30, 0, 0, time.UTC)
	expectedAppointment := &appointments.Appointment{
		ID:              1,
		UserID:          1,
		DoctorID:        1,
		AppointmentDate: appointmentDate,
		Status:          "pending",
		Notes:           "Regular checkup",
	}

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(expectedAppointment, nil).
		Times(1)

	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAppointment.ID, result.ID)
	assert.Equal(t, expectedAppointment.UserID, result.UserID)
	assert.Equal(t, expectedAppointment.Status, result.Status)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := appointments.NewMockAppointmentRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := appointments.NewService(logger, mockRepo)

	mockRepo.EXPECT().
		GetByID(int64(999)).
		Return(nil, errors.New("appointment not found")).
		Times(1)

	result, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
