package billings_test

import (
	"Dedenruslan19/med-project/service/billings"
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

	mockRepo := billings.NewMockBillingRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := billings.NewService(logger, mockRepo)

	expectedBilling := &billings.Billing{
		ID:            1,
		AppointmentID: 1,
		TotalAmount:   250000.0,
		PaymentStatus: "waiting_payment",
	}

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(expectedBilling, nil).
		Times(1)

	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBilling.ID, result.ID)
	assert.Equal(t, expectedBilling.TotalAmount, result.TotalAmount)
	assert.Equal(t, expectedBilling.PaymentStatus, result.PaymentStatus)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := billings.NewMockBillingRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := billings.NewService(logger, mockRepo)

	mockRepo.EXPECT().
		GetByID(int64(999)).
		Return(nil, errors.New("billing not found")).
		Times(1)

	result, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
