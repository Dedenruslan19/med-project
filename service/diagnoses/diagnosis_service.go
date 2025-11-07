package diagnoses

import (
	"log/slog"
	"strings"

	"Dedenruslan19/med-project/service/appointments"
	errs "Dedenruslan19/med-project/service/errors"
)

type service struct {
	repo               DiagnoseRepo
	appointmentService appointments.Service
	logger             *slog.Logger
}

type Service interface {
	Create(diagnose *Diagnose) (int64, error)
	GetByID(id int64) (*Diagnose, error)
	GetByAppointmentID(appointmentID int64) (*Diagnose, error)
	Update(diagnose *Diagnose) error
	CalculateTotalAmount(diagnose *Diagnose) float64
}

func NewService(logger *slog.Logger, repo DiagnoseRepo, appointmentService appointments.Service) Service {
	return &service{
		logger:             logger,
		repo:               repo,
		appointmentService: appointmentService,
	}
}

func (s *service) Create(diagnose *Diagnose) (int64, error) {
	if diagnose.AppointmentID == 0 {
		return 0, errs.ErrInvalidInput
	}

	id, err := s.repo.Create(diagnose)
	if err != nil {
		s.logger.Error("failed to create diagnose",
			slog.Any("error", err),
			slog.Int64("appointment_id", diagnose.AppointmentID),
		)
		return 0, err
	}

	err = s.appointmentService.UpdateStatus(diagnose.AppointmentID, "completed")
	if err != nil {
		s.logger.Error("failed to update appointment status after diagnose",
			slog.Any("error", err),
			slog.Int64("appointment_id", diagnose.AppointmentID))
	} else {
		s.logger.Info("appointment status updated to completed, doctor is now available",
			slog.Int64("appointment_id", diagnose.AppointmentID))
	}

	return id, nil
}

func (s *service) GetByID(id int64) (*Diagnose, error) {
	diagnose, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get diagnose by ID",
			slog.Any("error", err),
			slog.Int64("diagnose_id", id),
		)
		return nil, err
	}
	return diagnose, nil
}

func (s *service) GetByAppointmentID(appointmentID int64) (*Diagnose, error) {
	diagnose, err := s.repo.GetByAppointmentID(appointmentID)
	if err != nil {
		s.logger.Error("failed to get diagnose by appointment ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", appointmentID),
		)
		return nil, err
	}
	return diagnose, nil
}

func (s *service) Update(diagnose *Diagnose) error {
	err := s.repo.Update(diagnose)
	if err != nil {
		s.logger.Error("failed to update diagnose",
			slog.Any("error", err),
			slog.Int64("diagnose_id", diagnose.ID),
		)
		return err
	}
	return nil
}

func (s *service) CalculateTotalAmount(diagnosis *Diagnose) float64 {
	const appointmentFee = 200000.0
	const medicationFee = 50000.0

	totalAmount := appointmentFee

	if diagnosis.PrescribedMedications != "" {
		medications := strings.Split(diagnosis.PrescribedMedications, ",")
		medicationCount := 0
		for _, med := range medications {
			if strings.TrimSpace(med) != "" {
				medicationCount++
			}
		}
		totalAmount += float64(medicationCount) * medicationFee
	}

	return totalAmount
}
