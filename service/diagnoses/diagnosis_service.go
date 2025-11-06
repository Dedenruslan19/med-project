package diagnoses

import (
	"log/slog"
	"strings"
)

type service struct {
	repo   DiagnosisRepo
	logger *slog.Logger
}

type Service interface {
	Create(diagnosis *Diagnosis) (int64, error)
	GetByID(id int64) (*Diagnosis, error)
	GetByAppointmentID(appointmentID int64) (*Diagnosis, error)
	Update(diagnosis *Diagnosis) error
	CalculateTotalAmount(diagnosis *Diagnosis) float64
}

func NewService(logger *slog.Logger, repo DiagnosisRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) Create(diagnosis *Diagnosis) (int64, error) {
	id, err := s.repo.Create(diagnosis)
	if err != nil {
		s.logger.Error("Failed to create diagnosis",
			slog.Any("error", err),
			slog.Int64("appointment_id", diagnosis.AppointmentID),
		)
		return 0, err
	}
	return id, nil
}

func (s *service) GetByID(id int64) (*Diagnosis, error) {
	diagnosis, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get diagnosis by ID",
			slog.Any("error", err),
			slog.Int64("diagnosis_id", id),
		)
		return nil, err
	}
	return diagnosis, nil
}

func (s *service) GetByAppointmentID(appointmentID int64) (*Diagnosis, error) {
	diagnosis, err := s.repo.GetByAppointmentID(appointmentID)
	if err != nil {
		s.logger.Error("Failed to get diagnosis by appointment ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", appointmentID),
		)
		return nil, err
	}
	return diagnosis, nil
}

func (s *service) Update(diagnosis *Diagnosis) error {
	err := s.repo.Update(diagnosis)
	if err != nil {
		s.logger.Error("Failed to update diagnosis",
			slog.Any("error", err),
			slog.Int64("diagnosis_id", diagnosis.ID),
		)
		return err
	}
	return nil
}

// CalculateTotalAmount menghitung total biaya: 200.000 (appointment) + biaya obat
// Setiap obat yang disebutkan dalam prescribed_medications = 50.000
func (s *service) CalculateTotalAmount(diagnosis *Diagnosis) float64 {
	const appointmentFee = 200000.0
	const medicationFee = 50000.0

	totalAmount := appointmentFee

	if diagnosis.PrescribedMedications != "" {
		// Hitung jumlah obat berdasarkan koma atau newline
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
