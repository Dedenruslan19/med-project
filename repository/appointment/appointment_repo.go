package appointment

import (
	"Dedenruslan19/med-project/service/appointments"
	"log/slog"
	"time"

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
		r.logger.Error("failed to create appointment",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", appointment.ID),
		)
		return 0, result.Error
	}
	return appointment.ID, nil
}

func (r *appointmentRepo) GetByID(id int64) (*appointments.Appointment, error) {
	var appointment appointments.Appointment
	result := r.db.Where("id = ?", id).First(&appointment)
	if result.Error != nil {
		r.logger.Error("failed to get appointment by ID",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", id),
		)
		return nil, result.Error
	}
	return &appointment, nil
}

func (r *appointmentRepo) GetByUserID(userID int64) ([]appointments.Appointment, error) {
	var appointmentList []appointments.Appointment
	result := r.db.Where("user_id = ?", userID).Find(&appointmentList)
	if result.Error != nil {
		r.logger.Error("failed to get appointments by user ID",
			slog.Any("error", result.Error),
			slog.Int64("user_id", userID),
		)
		return nil, result.Error
	}
	return appointmentList, nil
}

func (r *appointmentRepo) UpdateStatus(id int64, status string) error {
	result := r.db.Model(&appointments.Appointment{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		r.logger.Error("failed to update appointment status",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", id),
			slog.String("status", status),
		)
		return result.Error
	}
	return nil
}

func (r *appointmentRepo) IsDoctorAvailable(doctorID int64, appointmentDate time.Time) (bool, error) {
	var count int64

	err := r.db.Model(&appointments.Appointment{}).
		Joins("LEFT JOIN diagnoses d ON d.appointment_id = appointments.id").
		Where("appointments.doctor_id = ? AND appointments.appointment_date = ? AND d.id IS NULL",
			doctorID, appointmentDate).
		Count(&count).Error
	if err != nil {
		r.logger.Error("failed to check doctor availability",
			slog.Any("error", err),
			slog.Int64("doctor_id", doctorID),
			slog.Time("appointment_date", appointmentDate),
		)
		return false, err
	}

	return count == 0, nil
}
