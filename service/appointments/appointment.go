package appointments

import "time"

type Appointment struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int64     `json:"user_id" gorm:"not null;index"`
	DoctorID        int64     `json:"doctor_id" gorm:"not null;index"`
	Date            string    `json:"date" gorm:"type:date;not null"`
	AppointmentDate time.Time `json:"appointment_date" gorm:"not null"`
	Status          string    `json:"status" gorm:"type:varchar(50);not null"`
	Notes           string    `json:"notes" gorm:"type:text"`
}
