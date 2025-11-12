package invoices

type InvoiceRepo interface {
	Create(invoice *Invoice) (int64, error)
	GetByID(id int64) (*Invoice, error)
	GetByBillingID(billingID int64) (*Invoice, error)
	UpdateSentAt(id int64) error
	SendInvoiceEmail(id int64, email string) error
}
