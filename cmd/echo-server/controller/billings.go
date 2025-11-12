package controller

import (
	"Dedenruslan19/med-project/cmd/echo-server/middleware"
	"Dedenruslan19/med-project/service/appointments"
	"Dedenruslan19/med-project/service/billings"
	"Dedenruslan19/med-project/service/invoices"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type BillingController struct {
	service            billings.Service
	invoiceService     invoices.Service
	appointmentService appointments.Service
	validate           *validator.Validate
	logger             *slog.Logger
}

func NewBillingController(service billings.Service, invoiceService invoices.Service, appointmentService appointments.Service, logger *slog.Logger) *BillingController {
	return &BillingController{
		service:            service,
		invoiceService:     invoiceService,
		appointmentService: appointmentService,
		validate:           validator.New(),
		logger:             logger,
	}
}

type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" validate:"required,oneof=unpaid waiting_payment paid failed"`
}

type CreateInvoiceRequest struct {
	PayerEmail  string `json:"payer_email" validate:"required,email"`
	Description string `json:"description"`
}

func (bc *BillingController) GetBillingByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid billing ID",
		})
	}

	billing, err := bc.service.GetByID(id)
	if err != nil {
		bc.logger.Error("failed to get billing",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "billing not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "billing retrieved successfully",
		"data":    billing,
	})
}

func (bc *BillingController) GetBillingByAppointmentID(c echo.Context) error {
	idParam := c.Param("appointment_id")
	appointmentID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid appointment ID",
		})
	}

	billing, err := bc.service.GetByAppointmentID(appointmentID)
	if err != nil {
		bc.logger.Error("Failed to get billing by appointment ID",
			slog.Any("error", err),
			slog.Int64("appointment_id", appointmentID),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Billing not found for this appointment",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Billing retrieved successfully",
		"data":    billing,
	})
}

func (bc *BillingController) CreateInvoice(c echo.Context) error {
	idParam := c.Param("id")
	billingID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid billing ID",
		})
	}

	var req CreateInvoiceRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if validateErr := bc.validate.Struct(req); validateErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": validateErr.Error(),
		})
	}

	// Get billing details
	billing, err := bc.service.GetByID(billingID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Billing not found",
		})
	}

	// Authorization: ensure caller is the appointment's doctor
	doctorIDFromToken, ok := middleware.GetUserID(c)
	if !ok {
		bc.logger.Error("failed to get doctor id from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	appointment, err := bc.appointmentService.GetByID(billing.AppointmentID)
	if err != nil {
		bc.logger.Error("failed to get appointment for ownership check",
			slog.Any("error", err),
			slog.Int64("appointment_id", billing.AppointmentID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to validate ownership"})
	}

	if appointment.DoctorID != doctorIDFromToken {
		bc.logger.Warn("doctor attempting to create invoice for another doctor's appointment",
			slog.Int64("token_doctor_id", doctorIDFromToken),
			slog.Int64("appointment_doctor_id", appointment.DoctorID),
			slog.Int64("appointment_id", appointment.ID),
		)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You are not authorized to create invoice for this appointment"})
	}

	// Check if already paid
	if billing.PaymentStatus == "paid" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Billing already paid",
		})
	}

	// Update billing status to waiting_payment
	err = bc.service.UpdatePaymentStatus(billingID, "waiting_payment")
	if err != nil {
		bc.logger.Error("Failed to update billing status",
			slog.Any("error", err),
			slog.Int64("billing_id", billingID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update billing",
		})
	}

	bc.logger.Info("Invoice created",
		slog.Int64("billing_id", billingID),
		slog.String("payer_email", req.PayerEmail),
		slog.Float64("amount", billing.TotalAmount))

	// Create invoice record and return it so Postman shows the invoice details
	consultationFee := billing.TotalAmount
	medicationFee := 0.0
	totalAmount := consultationFee + medicationFee

	invoice, err := bc.invoiceService.CreateInvoice(billingID, consultationFee, medicationFee, totalAmount, req.PayerEmail)
	if err != nil {
		bc.logger.Error("Failed to create invoice record",
			slog.Any("error", err),
			slog.Int64("billing_id", billingID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create invoice",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Payment invoice created successfully",
		"data":    invoice,
	})
}

func (bc *BillingController) UpdatePaymentStatus(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid billing ID",
		})
	}

	var req UpdatePaymentStatusRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if validateErr := bc.validate.Struct(req); validateErr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": validateErr.Error(),
		})
	}

	bc.logger.Info("Updating payment status",
		slog.Int64("billing_id", id),
		slog.String("payment_status", req.PaymentStatus),
		slog.String("payment_status_lower", strings.ToLower(strings.TrimSpace(req.PaymentStatus))),
	)

	// Update payment status
	err = bc.service.UpdatePaymentStatus(id, req.PaymentStatus)
	if err != nil {
		bc.logger.Error("Failed to update payment status",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update payment status",
		})
	}

	// Auto-create invoice when payment status is "paid"
	if strings.ToLower(strings.TrimSpace(req.PaymentStatus)) == "paid" {
		// Get billing details
		_, err := bc.service.GetByID(id)
		if err != nil {
			bc.logger.Error("Failed to get billing for invoice creation",
				slog.Any("error", err),
				slog.Int64("billing_id", id),
			)
			// Don't fail the payment status update, just log the error
		} else {
			// In a real scenario, you would fetch the diagnosis to calculate medication fees
			consultationFee := 200000.0
			medicationFee := 0.0
			totalAmount := consultationFee + medicationFee

			// For now, use a default email. In production, fetch user email from appointment
			email := "user@example.com"

			invoice, err := bc.invoiceService.CreateInvoice(id, consultationFee, medicationFee, totalAmount, email)
			if err != nil {
				bc.logger.Error("Failed to auto-create invoice",
					slog.Any("error", err),
					slog.Int64("billing_id", id),
				)
			} else {
				bc.logger.Info("Invoice auto-created successfully",
					slog.Int64("billing_id", id),
					slog.Int64("invoice_id", invoice.ID),
					slog.String("invoice_number", invoice.InvoiceNumber),
				)
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("Payment status updated successfully%s",
			func() string {
				if strings.ToLower(strings.TrimSpace(req.PaymentStatus)) == "paid" {
					return " and invoice created"
				}
				return ""
			}()),
	})
}
