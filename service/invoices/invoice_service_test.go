package invoices_test

import (
	"Dedenruslan19/med-project/service/invoices"
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

	mockRepo := invoices.NewMockInvoiceRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := invoices.NewService(logger, mockRepo)

	expectedInvoice := &invoices.Invoice{
		ID:              1,
		BillingID:       1,
		InvoiceNumber:   "INV-1-1234567890",
		ConsultationFee: 200000.0,
		MedicationFee:   50000.0,
		TotalAmount:     250000.0,
		SentToEmail:     "user@example.com",
	}

	mockRepo.EXPECT().
		GetByID(int64(1)).
		Return(expectedInvoice, nil).
		Times(1)

	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedInvoice.ID, result.ID)
	assert.Equal(t, expectedInvoice.InvoiceNumber, result.InvoiceNumber)
	assert.Equal(t, expectedInvoice.TotalAmount, result.TotalAmount)
}

func TestGetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := invoices.NewMockInvoiceRepo(ctrl)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	service := invoices.NewService(logger, mockRepo)

	mockRepo.EXPECT().
		GetByID(int64(999)).
		Return(nil, errors.New("invoice not found")).
		Times(1)

	result, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
}
