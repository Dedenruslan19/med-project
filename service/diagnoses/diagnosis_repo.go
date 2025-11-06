package diagnoses

type DiagnosisRepo interface {
	Create(diagnosis *Diagnosis) (int64, error)
	GetByID(id int64) (*Diagnosis, error)
	GetByAppointmentID(appointmentID int64) (*Diagnosis, error)
	Update(diagnosis *Diagnosis) error
}
