package billings

type BillingRepo interface {
	Create(billing *Billing) (int64, error)
	GetByID(id int64) (*Billing, error)
	GetByAppointmentID(appointmentID int64) (*Billing, error)
	UpdatePaymentStatus(id int64, status string, invoiceURL string) error
	UpdatePaidAt(id int64) error
}
