package billings

import (
	"log/slog"
	"time"
)

type service struct {
	repo   BillingRepo
	logger *slog.Logger
}

type Service interface {
	Create(billing *Billing) (int64, error)
	GetByID(id int64) (*Billing, error)
	GetByAppointmentID(appointmentID int64) (*Billing, error)
	UpdatePaymentStatus(id int64, status string, invoiceURL string) error
	MarkAsPaid(id int64) error
}

func NewService(logger *slog.Logger, repo BillingRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) Create(billing *Billing) (int64, error) {
	id, err := s.repo.Create(billing)
	if err != nil {
		s.logger.Error("Failed to create billing",
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
		s.logger.Error("Failed to get billing by ID",
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
		s.logger.Error("Failed to get billing by appointment ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", appointmentID),
		)
		return nil, err
	}
	return billing, nil
}

func (s *service) UpdatePaymentStatus(id int64, status string, invoiceURL string) error {
	err := s.repo.UpdatePaymentStatus(id, status, invoiceURL)
	if err != nil {
		s.logger.Error("Failed to update payment status",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
			slog.String("status", status),
		)
		return err
	}
	return nil
}

func (s *service) MarkAsPaid(id int64) error {
	err := s.repo.UpdatePaidAt(id)
	if err != nil {
		s.logger.Error("Failed to mark billing as paid",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
			slog.Time("paid_at", time.Now()),
		)
		return err
	}
	return nil
}
