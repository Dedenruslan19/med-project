package invoices

import (
	"Dedenruslan19/med-project/repository/notification"
	"fmt"
	"log/slog"
	"time"
)

type service struct {
	repo        InvoiceRepo
	logger      *slog.Logger
	emailSender *notification.SMTPSender
}

type Service interface {
	CreateInvoice(billingID int64, consultationFee, medicationFee, totalAmount float64, email string) (*Invoice, error)
	GetByID(id int64) (*Invoice, error)
	GetByBillingID(billingID int64) (*Invoice, error)
	MarkAsSent(id int64) error
	SendInvoice(billingID int64, email string) (*Invoice, error)
}

func NewService(logger *slog.Logger, repo InvoiceRepo, emailSender *notification.SMTPSender) Service {
	return &service{
		logger:      logger,
		repo:        repo,
		emailSender: emailSender,
	}
}

func (s *service) CreateInvoice(billingID int64, consultationFee, medicationFee, totalAmount float64, email string) (*Invoice, error) {
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

func (s *service) SendInvoice(billingID int64, email string) (*Invoice, error) {
	// Create a new invoice record when SendInvoice is called
	consultationFee := 100.0
	medicationFee := 50.0
	totalAmount := consultationFee + medicationFee

	invoice, err := s.CreateInvoice(billingID, consultationFee, medicationFee, totalAmount, email)
	if err != nil {
		s.logger.Error("failed to create invoice for sending",
			slog.Any("error", err),
			slog.Int64("billing_id", billingID),
		)
		return nil, err
	}

	// Build email content
	subject := fmt.Sprintf("Invoice %s", invoice.InvoiceNumber)
	body := fmt.Sprintf("Dear Customer,\n\nPlease find your invoice details below:\n\nInvoice Number: %s\nConsultation Fee: %.2f\nMedication Fee: %.2f\nTotal Amount: %.2f\n\nThank you for your business.",
		invoice.InvoiceNumber, invoice.ConsultationFee, invoice.MedicationFee, invoice.TotalAmount)

	if s.emailSender == nil {
		s.logger.Warn("email sender not configured, invoice created but not sent",
			slog.Int64("billing_id", billingID),
			slog.Int64("invoice_id", invoice.ID),
		)
		return invoice, nil
	}

	if err := s.emailSender.Send(email, subject, body); err != nil {
		s.logger.Error("failed to send invoice email",
			slog.Any("error", err),
			slog.Int64("billing_id", billingID),
			slog.String("email", email),
		)
		return nil, err
	}

	// Mark invoice as sent
	if err := s.MarkAsSent(invoice.ID); err != nil {
		s.logger.Error("failed to mark invoice as sent after email",
			slog.Any("error", err),
			slog.Int64("invoice_id", invoice.ID),
		)
		return nil, err
	}

	return invoice, nil
}
