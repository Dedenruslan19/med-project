package appointment

import (
	"Dedenruslan19/med-project/service/appointments"
	"log/slog"

	"gorm.io/gorm"
)

type appointmentRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewAppointmentRepo(db *gorm.DB, logger *slog.Logger) appointments.AppointmentRepo {
	return &appointmentRepo{db: db, logger: logger}
}

func (r *appointmentRepo) Create(appointment *appointments.Appointment) (int64, error) {
	result := r.db.Create(appointment)
	if result.Error != nil {
		r.logger.Error("Failed to create appointment",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", appointment.ID),
		)
		return 0, result.Error
	}
	return appointment.ID, nil
}
