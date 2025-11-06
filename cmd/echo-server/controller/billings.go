package controller

import (
	"Dedenruslan19/med-project/service/billings"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type BillingController struct {
	service  billings.Service
	validate *validator.Validate
	logger   *slog.Logger
}

func NewBillingController(service billings.Service, logger *slog.Logger) *BillingController {
	return &BillingController{
		service:  service,
		validate: validator.New(),
		logger:   logger,
	}
}

type CreateBillingRequest struct {
	AppointmentID int64   `json:"appointment_id" validate:"required"`
	TotalAmount   float64 `json:"total_amount" validate:"required,gt=0"`
	ExternalID    string  `json:"external_id"`
	InvoiceURL    string  `json:"invoice_url"`
}

func (bc *BillingController) CreateBilling(c echo.Context) error {
	var req CreateBillingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if err := bc.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	billing := &billings.Billing{
		AppointmentID: req.AppointmentID,
		TotalAmount:   req.TotalAmount,
		PaymentStatus: "unpaid",
		ExternalID:    req.ExternalID,
		InvoiceURL:    req.InvoiceURL,
	}

	id, err := bc.service.Create(billing)
	if err != nil {
		bc.logger.Error("Failed to create billing",
			slog.Any("error", err),
			slog.Int64("appointment_id", req.AppointmentID),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create billing",
		})
	}

	billing.ID = id

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Billing created successfully",
		"data":    billing,
	})
}

func (bc *BillingController) GetBillingByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid billing ID",
		})
	}

	billing, err := bc.service.GetByID(id)
	if err != nil {
		bc.logger.Error("Failed to get billing",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Billing not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Billing retrieved successfully",
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

type UpdatePaymentStatusRequest struct {
	PaymentStatus string `json:"payment_status" validate:"required,oneof=unpaid waiting_payment paid failed"`
	InvoiceURL    string `json:"invoice_url"`
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

	err = bc.service.UpdatePaymentStatus(id, req.PaymentStatus, req.InvoiceURL)
	if err != nil {
		bc.logger.Error("Failed to update payment status",
			slog.Any("error", err),
			slog.Int64("billing_id", id),
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update payment status",
		})
	}

	// If status is paid, update paid_at timestamp
	if req.PaymentStatus == "paid" {
		if paidErr := bc.service.MarkAsPaid(id); paidErr != nil {
			bc.logger.Error("Failed to mark as paid",
				slog.Any("error", paidErr),
				slog.Int64("billing_id", id),
			)
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Payment status updated successfully",
	})
}
