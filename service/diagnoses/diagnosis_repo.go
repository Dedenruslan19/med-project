package diagnoses

type DiagnoseRepo interface {
	Create(diagnose *Diagnose) (int64, error)
	GetByID(id int64) (*Diagnose, error)
	GetByAppointmentID(appointmentID int64) (*Diagnose, error)
	Update(diagnose *Diagnose) error
}
