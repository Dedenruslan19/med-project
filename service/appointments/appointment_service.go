package appointments

import (
	errs "Dedenruslan19/med-project/service/errors"
	"log/slog"
)

type service struct {
	repo   AppointmentRepo
	logger *slog.Logger
}

type Service interface {
	Create(appointment *Appointment) (int64, error)
	GetByID(id int64) (*Appointment, error)
	GetByUserID(userID int64) ([]Appointment, error)
	UpdateStatus(id int64, status string) error
}

func NewService(logger *slog.Logger, repo AppointmentRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) Create(appointment *Appointment) (int64, error) {
	available, err := s.repo.IsDoctorAvailable(appointment.DoctorID, appointment.AppointmentDate)
	if err != nil {
		s.logger.Error("failed to check doctor availabilities",
			slog.Any("error", err),
			slog.Int64("doctor_id", appointment.DoctorID),
		)
		return 0, err
	}

	if !available {
		s.logger.Warn("doctor is busy at the requested time",
			slog.Int64("doctor_id", appointment.DoctorID),
			slog.Time("appointment_date", appointment.AppointmentDate),
		)
		return 0, errs.ErrDoctorBusy
	}

	id, err := s.repo.Create(appointment)
	if err != nil {
		s.logger.Error("failed to create appointment",
			slog.Any("error", err),
			slog.Int64("appointment_id", id),
		)
		return 0, err
	}
	return id, nil
}

func (s *service) GetByID(id int64) (*Appointment, error) {
	appointment, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get appointment by ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", id),
		)
		return nil, err
	}
	return appointment, nil
}

func (s *service) GetByUserID(userID int64) ([]Appointment, error) {
	appointments, err := s.repo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("failed to get appointments by user ID",
			slog.Any("error", err),
			slog.Int64("user_id", userID),
		)
		return nil, err
	}
	return appointments, nil
}

func (s *service) UpdateStatus(id int64, status string) error {
	err := s.repo.UpdateStatus(id, status)
	if err != nil {
		s.logger.Error("failed to update appointment status",
			slog.Any("error", err),
			slog.Int64("appointment_id", id),
			slog.String("status", status),
		)
		return err
	}
	return nil
}
