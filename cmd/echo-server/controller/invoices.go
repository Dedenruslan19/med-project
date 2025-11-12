package controller

import (
	"Dedenruslan19/med-project/cmd/echo-server/middleware"
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

	// Authorization: ensure caller is the appointment's doctor
	doctorIDFromToken, ok := middleware.GetUserID(c)
	if !ok {
		ic.logger.Error("Failed to get doctor ID from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
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

	if appointment.DoctorID != doctorIDFromToken {
		ic.logger.Warn("doctor attempting to access invoice for another doctor's appointment",
			slog.Int64("token_doctor_id", doctorIDFromToken),
			slog.Int64("appointment_doctor_id", appointment.DoctorID),
			slog.Int64("appointment_id", appointment.ID),
		)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not authorized to view this invoice"})
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

func (ic *InvoiceController) SendInvoice(c echo.Context) error {
	type SendInvoiceRequest struct {
		BillingID int64  `json:"billing_id" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
	}

	var req SendInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	// Authorization: ensure caller is the appointment's doctor for this billing
	doctorIDFromToken, ok := middleware.GetUserID(c)
	if !ok {
		ic.logger.Error("Failed to get doctor ID from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	billing, err := ic.billingService.GetByID(req.BillingID)
	if err != nil {
		ic.logger.Error("Failed to get billing for ownership check",
			slog.Any("error", err),
			slog.Int64("billing_id", req.BillingID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get billing"})
	}

	appointment, err := ic.appointmentService.GetByID(billing.AppointmentID)
	if err != nil {
		ic.logger.Error("Failed to get appointment for ownership check",
			slog.Any("error", err),
			slog.Int64("appointment_id", billing.AppointmentID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get appointment"})
	}

	if appointment.DoctorID != doctorIDFromToken {
		ic.logger.Warn("doctor attempting to send invoice for another doctor's appointment",
			slog.Int64("token_doctor_id", doctorIDFromToken),
			slog.Int64("appointment_doctor_id", appointment.DoctorID),
			slog.Int64("billing_id", req.BillingID),
		)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not authorized to send invoice for this billing"})
	}

	invoice, err := ic.invoiceService.SendInvoice(req.BillingID, req.Email)
	if err != nil {
		ic.logger.Error("Failed to send invoice",
			slog.Any("error", err),
			slog.Int64("billing_id", req.BillingID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Invoice sent successfully",
		"data":    invoice,
	})
}
