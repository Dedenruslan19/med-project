package appointments

import "time"

type AppointmentRepo interface {
	Create(appointment *Appointment) (int64, error)
	GetByID(id int64) (*Appointment, error)
	GetByUserID(userID int64) ([]Appointment, error)
	UpdateStatus(id int64, status string) error
	IsDoctorAvailable(doctorID int64, appointmentDate time.Time) (bool, error)
}
