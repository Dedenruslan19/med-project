package billing

import (
	"Dedenruslan19/med-project/service/billings"
	"log/slog"

	"gorm.io/gorm"
)

type billingRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewBillingRepo(db *gorm.DB, logger *slog.Logger) billings.BillingRepo {
	return &billingRepo{db: db, logger: logger}
}

func (r *billingRepo) Create(billing *billings.Billing) (int64, error) {
	result := r.db.Create(billing)
	if result.Error != nil {
		r.logger.Error("Failed to create billing",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", billing.AppointmentID),
		)
		return 0, result.Error
	}
	return billing.ID, nil
}

func (r *billingRepo) GetByID(id int64) (*billings.Billing, error) {
	var billing billings.Billing
	result := r.db.Where("id = ?", id).First(&billing)
	if result.Error != nil {
		r.logger.Error("Failed to get billing by ID",
			slog.Any("error", result.Error),
			slog.Int64("billing_id", id),
		)
		return nil, result.Error
	}
	return &billing, nil
}

func (r *billingRepo) GetByAppointmentID(appointmentID int64) (*billings.Billing, error) {
	var billing billings.Billing
	result := r.db.Where("appointment_id = ?", appointmentID).First(&billing)
	if result.Error != nil {
		r.logger.Error("Failed to get billing by appointment ID",
			slog.Any("error", result.Error),
			slog.Int64("appointment_id", appointmentID),
		)
		return nil, result.Error
	}
	return &billing, nil
}

func (r *billingRepo) Update(billing *billings.Billing) error {
	result := r.db.Save(billing)
	if result.Error != nil {
		r.logger.Error("Failed to update billing",
			slog.Any("error", result.Error),
			slog.Int64("billing_id", billing.ID),
		)
		return result.Error
	}
	return nil
}
