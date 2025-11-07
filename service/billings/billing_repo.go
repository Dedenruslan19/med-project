package billings

type BillingRepo interface {
	Create(billing *Billing) (int64, error)
	GetByID(id int64) (*Billing, error)
	GetByAppointmentID(appointmentID int64) (*Billing, error)
	Update(billing *Billing) error
}
