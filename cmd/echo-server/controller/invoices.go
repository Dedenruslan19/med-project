package controller

import (
	"Dedenruslan19/med-project/service/appointments"
	"Dedenruslan19/med-project/service/billings"
	"Dedenruslan19/med-project/service/diagnoses"
	"Dedenruslan19/med-project/service/doctors"
	"Dedenruslan19/med-project/service/invoices"
	"Dedenruslan19/med-project/service/users"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type InvoiceController struct {
	invoiceService     invoices.Service
	billingService     billings.Service
	appointmentService appointments.Service
	diagnoseService    diagnoses.Service
	userService        users.Service
	doctorService      doctors.Service
	logger             *slog.Logger
}

func NewInvoiceController(
	invoiceService invoices.Service,
	billingService billings.Service,
	appointmentService appointments.Service,
	diagnoseService diagnoses.Service,
	userService users.Service,
	doctorService doctors.Service,
	logger *slog.Logger,
) *InvoiceController {
	return &InvoiceController{
		invoiceService:     invoiceService,
		billingService:     billingService,
		appointmentService: appointmentService,
		diagnoseService:    diagnoseService,
		userService:        userService,
		doctorService:      doctorService,
		logger:             logger,
	}
}

// GetInvoiceByBillingID - Get invoice with complete details
func (ic *InvoiceController) GetInvoiceByBillingID(c echo.Context) error {
	billingIDParam := c.Param("billing_id")
	billingID, err := strconv.ParseInt(billingIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid billing ID",
		})
	}
	invoice, err := ic.invoiceService.GetByBillingID(billingID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Invoice not found",
		})
	}
	billing, err := ic.billingService.GetByID(invoice.BillingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get billing details",
		})
	}
	appointment, err := ic.appointmentService.GetByID(billing.AppointmentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get appointment details",
		})
	}
	user, err := ic.userService.GetUserByID(appointment.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get user details",
		})
	}
	doctor, err := ic.doctorService.GetByID(appointment.DoctorID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get doctor details",
		})
	}
	diagnosis, err := ic.diagnoseService.GetByAppointmentID(appointment.ID)
	if err != nil {
		ic.logger.Warn("Diagnosis not found", slog.Int64("appointment_id", appointment.ID))
	}
	response := map[string]interface{}{
		"invoice_number":   invoice.InvoiceNumber,
		"invoice_date":     invoice.CreatedAt,
		"consultation_fee": invoice.ConsultationFee,
		"medication_fee":   invoice.MedicationFee,
		"total_amount":     invoice.TotalAmount,
		"payment_status":   billing.PaymentStatus,
		"paid_at":          billing.PaidAt,
		"sent_at":          invoice.SentAt,
		"patient": map[string]interface{}{
			"name":  user.FullName,
			"email": user.Email,
		},
		"doctor": map[string]interface{}{
			"name":           doctor.FullName,
			"specialization": doctor.Specialization,
		},
		"appointment": map[string]interface{}{
			"date":   appointment.AppointmentDate,
			"status": appointment.Status,
		},
	}

	if diagnosis != nil {
		medicationCount := 0
		if diagnosis.PrescribedMedications != "" {
			medications := strings.Split(diagnosis.PrescribedMedications, ",")
			medicationCount = len(medications)
		}

		response["diagnosis"] = map[string]interface{}{
			"notes":                  diagnosis.Notes,
			"prescribed_medications": diagnosis.PrescribedMedications,
			"medication_count":       medicationCount,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invoice retrieved successfully",
		"data":    response,
	})
}

// MarkInvoiceAsSent - Mark invoice as sent to email
func (ic *InvoiceController) MarkInvoiceAsSent(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid invoice ID",
		})
	}

	err = ic.invoiceService.MarkAsSent(id)
	if err != nil {
		ic.logger.Error("Failed to mark invoice as sent",
			slog.Any("error", err),
			slog.Int64("invoice_id", id),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to mark invoice as sent",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Invoice marked as sent successfully",
	})
}
