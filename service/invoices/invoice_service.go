package invoices

import (
	"fmt"
	"log/slog"
	"time"
)

type service struct {
	repo   InvoiceRepo
	logger *slog.Logger
}

type Service interface {
	CreateInvoice(billingID int64, consultationFee, medicationFee, totalAmount float64, email string) (*Invoice, error)
	GetByID(id int64) (*Invoice, error)
	GetByBillingID(billingID int64) (*Invoice, error)
	MarkAsSent(id int64) error
}

func NewService(logger *slog.Logger, repo InvoiceRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) CreateInvoice(billingID int64, consultationFee, medicationFee, totalAmount float64, email string) (*Invoice, error) {
	// Generate invoice number
	invoiceNumber := fmt.Sprintf("INV-%d-%d", billingID, time.Now().Unix())

	invoice := &Invoice{
		BillingID:       billingID,
		InvoiceNumber:   invoiceNumber,
		ConsultationFee: consultationFee,
		MedicationFee:   medicationFee,
		TotalAmount:     totalAmount,
		SentToEmail:     email,
	}

	id, err := s.repo.Create(invoice)
	if err != nil {
		s.logger.Error("failed to create invoice",
			slog.Any("error", err),
			slog.Int64("billing_id", billingID),
		)
		return nil, err
	}

	invoice.ID = id
	return invoice, nil
}

func (s *service) GetByID(id int64) (*Invoice, error) {
	invoice, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get invoice by ID",
			slog.Any("error", err),
			slog.Int64("invoice_id", id),
		)
		return nil, err
	}
	return invoice, nil
}

func (s *service) GetByBillingID(billingID int64) (*Invoice, error) {
	invoice, err := s.repo.GetByBillingID(billingID)
	if err != nil {
		s.logger.Error("failed to get invoice by billing ID",
			slog.Any("error", err),
			slog.Int64("billing_id", billingID),
		)
		return nil, err
	}
	return invoice, nil
}

func (s *service) MarkAsSent(id int64) error {
	err := s.repo.UpdateSentAt(id)
	if err != nil {
		s.logger.Error("failed to mark invoice as sent",
			slog.Any("error", err),
			slog.Int64("invoice_id", id),
		)
		return err
	}
	return nil
}
