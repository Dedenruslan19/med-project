package appointments

import (
	"log/slog"
)

type service struct {
	repo   AppointmentRepo
	logger *slog.Logger
}

type Service interface {
	Create(appointment *Appointment) (int64, error)
}

func NewService(logger *slog.Logger, repo AppointmentRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) Create(appointment *Appointment) (int64, error) {
	id, err := s.repo.Create(appointment)
	if err != nil {
		s.logger.Error("Failed to create appointment",
			slog.Any("error", err),
			slog.Int64("appointment_id", id),
		)
		return 0, err
	}
	return id, nil
}
