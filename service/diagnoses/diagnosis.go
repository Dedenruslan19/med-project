package diagnoses

import "time"

type Diagnosis struct {
	ID                    int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AppointmentID         int64     `json:"appointment_id" gorm:"not null;index"`
	DoctorID              int64     `json:"doctor_id" gorm:"not null"`
	Notes                 string    `json:"notes" gorm:"type:text;not null"`
	PrescribedMedications string    `json:"prescribed_medications" gorm:"type:text"`
	CreatedAt             time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
