package diagnose

import (
	"Dedenruslan19/med-project/service/diagnoses"
	"log/slog"

	"gorm.io/gorm"
)

type diagnoseRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewDiagnoseRepo(db *gorm.DB, logger *slog.Logger) diagnoses.DiagnoseRepo {
	return &diagnoseRepo{db: db, logger: logger}
}

func (r *diagnoseRepo) Create(diagnose *diagnoses.Diagnose) (int64, error) {
	result := r.db.Create(diagnose)
	if result.Error != nil {
		r.logger.Error("failed to create diagnose",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", diagnose.AppointmentID),
		)
		return 0, result.Error
	}
	return diagnose.ID, nil
}

func (r *diagnoseRepo) GetByID(id int64) (*diagnoses.Diagnose, error) {
	var diagnose diagnoses.Diagnose
	result := r.db.Where("id = ?", id).First(&diagnose)
	if result.Error != nil {
		r.logger.Error("Failed to get diagnose by ID",
			slog.Any("error", result.Error),
			slog.Int64("diagnose_id", id),
		)
		return nil, result.Error
	}
	return &diagnose, nil
}

func (r *diagnoseRepo) GetByAppointmentID(appointmentID int64) (*diagnoses.Diagnose, error) {
	var diagnose diagnoses.Diagnose
	result := r.db.Where("appointment_id = ?", appointmentID).First(&diagnose)
	if result.Error != nil {
		r.logger.Error("failed to get diagnose by appointment ID",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", appointmentID),
		)
		return nil, result.Error
	}
	return &diagnose, nil
}

func (r *diagnoseRepo) Update(diagnose *diagnoses.Diagnose) error {
	result := r.db.Save(diagnose)
	if result.Error != nil {
		r.logger.Error("failed to update diagnose",
			slog.Any("error", result.Error),
			slog.Int64("diagnose_id", diagnose.ID),
		)
		return result.Error
	}
	return nil
}
