package billings

import "time"

type Billing struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	AppointmentID int64      `json:"appointment_id" gorm:"not null;index"`
	TotalAmount   float64    `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	PaymentStatus string     `json:"payment_status" gorm:"type:varchar(50);default:'unpaid'"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
