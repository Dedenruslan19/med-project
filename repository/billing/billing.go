package billing

import (
	"Dedenruslan19/med-project/service/billings"
	"log/slog"
	"time"

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

func (r *billingRepo) UpdatePaymentStatus(id int64, status string, invoiceURL string) error {
	result := r.db.Model(&billings.Billing{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status": status,
			"invoice_url":    invoiceURL,
		})
	if result.Error != nil {
		r.logger.Error("Failed to update payment status",
			slog.Any("error", result.Error),
			slog.Int64("billing_id", id),
			slog.String("status", status),
		)
		return result.Error
	}
	return nil
}

func (r *billingRepo) UpdatePaidAt(id int64) error {
	now := time.Now()
	result := r.db.Model(&billings.Billing{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status": "paid",
			"paid_at":        now,
		})
	if result.Error != nil {
		r.logger.Error("Failed to update paid_at",
			slog.Any("error", result.Error),
			slog.Int64("billing_id", id),
		)
		return result.Error
	}
	return nil
}
