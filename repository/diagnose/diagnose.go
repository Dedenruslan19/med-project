package diagnose

import (
	"Dedenruslan19/med-project/service/diagnoses"
	"log/slog"

	"gorm.io/gorm"
)

type diagnosisRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewDiagnosisRepo(db *gorm.DB, logger *slog.Logger) diagnoses.DiagnosisRepo {
	return &diagnosisRepo{db: db, logger: logger}
}

func (r *diagnosisRepo) Create(diagnosis *diagnoses.Diagnosis) (int64, error) {
	result := r.db.Create(diagnosis)
	if result.Error != nil {
		r.logger.Error("Failed to create diagnosis",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", diagnosis.AppointmentID),
		)
		return 0, result.Error
	}
	return diagnosis.ID, nil
}

func (r *diagnosisRepo) GetByID(id int64) (*diagnoses.Diagnosis, error) {
	var diagnosis diagnoses.Diagnosis
	result := r.db.Where("id = ?", id).First(&diagnosis)
	if result.Error != nil {
		r.logger.Error("Failed to get diagnosis by ID",
			slog.Any("error", result.Error),
			slog.Int64("diagnosis_id", id),
		)
		return nil, result.Error
	}
	return &diagnosis, nil
}

func (r *diagnosisRepo) GetByAppointmentID(appointmentID int64) (*diagnoses.Diagnosis, error) {
	var diagnosis diagnoses.Diagnosis
	result := r.db.Where("appointment_id = ?", appointmentID).First(&diagnosis)
	if result.Error != nil {
		r.logger.Error("Failed to get diagnosis by appointment ID",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", appointmentID),
		)
		return nil, result.Error
	}
	return &diagnosis, nil
}

func (r *diagnosisRepo) Update(diagnosis *diagnoses.Diagnosis) error {
	result := r.db.Save(diagnosis)
	if result.Error != nil {
		r.logger.Error("Failed to update diagnosis",
			slog.Any("error", result.Error),
			slog.Int64("diagnosis_id", diagnosis.ID),
		)
		return result.Error
	}
	return nil
}
