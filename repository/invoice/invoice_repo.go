package invoice

import (
	"Dedenruslan19/med-project/service/invoices"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type invoiceRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewInvoiceRepo(db *gorm.DB, logger *slog.Logger) invoices.InvoiceRepo {
	return &invoiceRepository{
		db:     db,
		logger: logger,
	}
}

func (r *invoiceRepository) Create(invoice *invoices.Invoice) (int64, error) {
	if err := r.db.Create(invoice).Error; err != nil {
		r.logger.Error("failed to create invoice", slog.Any("error", err))
		return 0, err
	}
	return invoice.ID, nil
}

func (r *invoiceRepository) GetByID(id int64) (*invoices.Invoice, error) {
	var invoice invoices.Invoice
	if err := r.db.First(&invoice, id).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) GetByBillingID(billingID int64) (*invoices.Invoice, error) {
	var invoice invoices.Invoice
	if err := r.db.Where("billing_id = ?", billingID).First(&invoice).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) UpdateSentAt(id int64) error {
	now := time.Now()
	return r.db.Model(&invoices.Invoice{}).Where("id = ?", id).Update("sent_at", now).Error
}

func (r *invoiceRepository) SendInvoiceEmail(id int64, email string) error {
	return nil
}
