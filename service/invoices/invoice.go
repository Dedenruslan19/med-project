package invoices

import "time"

type Invoice struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	BillingID       int64      `json:"billing_id" gorm:"not null;unique;index"`
	InvoiceNumber   string     `json:"invoice_number" gorm:"type:varchar(100);not null;unique"`
	ConsultationFee float64    `json:"consultation_fee" gorm:"type:decimal(10,2);not null;default:200000"`
	MedicationFee   float64    `json:"medication_fee" gorm:"type:decimal(10,2);not null;default:0"`
	TotalAmount     float64    `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	SentToEmail     string     `json:"sent_to_email" gorm:"type:varchar(255);not null"`
	SentAt          *time.Time `json:"sent_at"`
	CreatedAt       time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
