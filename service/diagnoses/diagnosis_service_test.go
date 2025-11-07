package diagnoses_test

import (
	"Dedenruslan19/med-project/service/diagnoses"
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

	mockRepo := diagnoses.NewMockDiagnoseRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := diagnoses.NewService(logger, mockRepo, nil)

	expectedDiagnosis := &diagnoses.Diagnose{
		ID:                    1,
		AppointmentID:         1,
		DoctorID:              1,
		Notes:                 "Common cold with fever",
		PrescribedMedications: "Paracetamol, Cough syrup",
	}

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(expectedDiagnosis, nil).
		Times(1)

	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedDiagnosis.ID, result.ID)
	assert.Equal(t, expectedDiagnosis.Notes, result.Notes)
	assert.Equal(t, expectedDiagnosis.PrescribedMedications, result.PrescribedMedications)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := diagnoses.NewMockDiagnoseRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := diagnoses.NewService(logger, mockRepo, nil)

	mockRepo.EXPECT().
		GetByID(int64(999)).
		Return(nil, errors.New("diagnosis not found")).
		Times(1)

	result, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
