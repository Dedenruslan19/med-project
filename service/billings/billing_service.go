package billings

import (
	"Dedenruslan19/med-project/service/invoices"
	"log/slog"
	"time"
)

type service struct {
	repo           BillingRepo
	invoiceService invoices.Service
	logger         *slog.Logger
}

type Service interface {
	Create(billing *Billing) (int64, error)
	GetByID(id int64) (*Billing, error)
	GetByAppointmentID(appointmentID int64) (*Billing, error)
	UpdatePaymentStatus(id int64, status string) error
	SetInvoiceService(invoiceService invoices.Service)
}

func NewService(logger *slog.Logger, repo BillingRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) SetInvoiceService(invoiceService invoices.Service) {
	s.invoiceService = invoiceService
}

func (s *service) Create(billing *Billing) (int64, error) {
	id, err := s.repo.Create(billing)
	if err != nil {
		s.logger.Error("failed to create billing",
			slog.Any("error", err),
			slog.Int64("appointment_id", billing.AppointmentID),
		)
		return 0, err
	}
	return id, nil
}

func (s *service) GetByID(id int64) (*Billing, error) {
	billing, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get billing by ID",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
		)
		return nil, err
	}
	return billing, nil
}

func (s *service) GetByAppointmentID(appointmentID int64) (*Billing, error) {
	billing, err := s.repo.GetByAppointmentID(appointmentID)
	if err != nil {
		s.logger.Error("failed to get billing by appointment ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", appointmentID),
		)
		return nil, err
	}
	return billing, nil
}

func (s *service) UpdatePaymentStatus(id int64, status string) error {
	billing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	billing.PaymentStatus = status

	// If status is paid, set paid_at timestamp
	if status == "paid" {
		now := time.Now()
		billing.PaidAt = &now
	}

	err = s.repo.Update(billing)
	if err != nil {
		s.logger.Error("failed to update payment status",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
			slog.String("status", status),
		)
		return err
	}
	return nil
}
