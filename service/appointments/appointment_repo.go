package appointments

type AppointmentRepo interface {
	Create(appointment *Appointment) (int64, error)
}
