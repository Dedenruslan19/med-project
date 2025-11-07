package appointments

import "time"

type Appointment struct {
	Status          string    `json:"status" gorm:"default:'pending'"`
	Notes           string    `json:"notes"`
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int64     `json:"user_id" gorm:"not null;index"`
	DoctorID        int64     `json:"doctor_id" gorm:"not null;index"`
	AppointmentDate time.Time `json:"appointment_date" gorm:"not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}
