package controller

import (
	errs "Dedenruslan19/med-project/service/errors"
	"Dedenruslan19/med-project/service/users"
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService users.Service
	validate    *validator.Validate
	logger      *slog.Logger
}

func NewUserController(us users.Service, logger *slog.Logger) *UserController {
	return &UserController{
		userService: us,
		validate:    validator.New(),
		logger:      logger,
	}
}

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	ErrInvalidRequestBody = map[string]interface{}{"message": "invalid request body"}
	ErrInvalidParams      = map[string]interface{}{"message": "invalid parameter"}
	ErrInternalServer     = map[string]interface{}{"message": "internal server error"}
	ErrHashFailed         = map[string]interface{}{"message": "failed to hash password"}
	ErrUnauthorized       = map[string]interface{}{"message": "unauthorized"}
	ErrInvalidContentType = map[string]interface{}{"message": "invalid content type"}
	ErrDataNotFound       = APIResponse{"data not found", map[string]interface{}{}}
	ErrEmailExists        = APIResponse{"email already exists", map[string]interface{}{}}
)

type RegisterInput struct {
	FullName string  `json:"full_name" validate:"required,min=2,max=255"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Weight   float64 `json:"weight" validate:"required,gt=0"`
	Height   float64 `json:"height" validate:"required,gt=0"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (uc *UserController) Register(c echo.Context) error {
	var input RegisterInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := uc.validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Validation failed",
			Data:    errorsMap,
		})
	}

	id, err := uc.userService.Register(users.User{
		FullName: input.FullName,
		Email:    input.Email,
		Password: input.Password,
		Weight:   input.Weight,
		Height:   input.Height,
	})

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrEmailAlreadyExists):
			return c.JSON(http.StatusBadRequest, ErrEmailExists)
		case errors.Is(err, errs.ErrInvalidInput):
			return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
		case errors.Is(err, errs.ErrHashFailed):
			uc.logger.Error("Failed to hash password",
				slog.String("email", input.Email),
				slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		default:
			uc.logger.Error("Internal server error on register",
				slog.String("email", input.Email),
				slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		}
	}

	res := APIResponse{
		Message: "user created successfully",
		Data: map[string]interface{}{
			"id":        id,
			"full_name": input.FullName,
			"email":     input.Email,
			"weight":    input.Weight,
			"height":    input.Height,
		},
	}
	return c.JSON(http.StatusCreated, res)
}

func (uc *UserController) Login(c echo.Context) error {
	var input LoginInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, ErrInvalidRequestBody)
	}

	if err := uc.validate.Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorsMap := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, APIResponse{
			Message: "Validation failed",
			Data:    errorsMap,
		})
	}

	user, err := uc.userService.Login(input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, ErrDataNotFound)
		case errors.Is(err, errs.ErrInvalidPass):
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "password doesn't match",
			})
		default:
			uc.logger.Error("Internal server error on login",
				slog.String("email", input.Email),
				slog.Any("error", err))
			return c.JSON(http.StatusInternalServerError, ErrInternalServer)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": user.Token,
	})
}

func (uc *UserController) GetMe(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	userID, ok := userIDInterface.(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	}

	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrDataNotFound)
	}

	heightInMeters := float64(user.Height) / 100.0
	bmiVal, usedCallback, err := uc.userService.CalculateBMI(
		float64(user.Weight),
		heightInMeters,
	)
	if err != nil {
		uc.logger.Error("Failed to calculate BMI",
			slog.Int64("user_id", userID),
			slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, ErrInternalServer)
	}

	message := "success"
	if usedCallback {
		message = "success from callback"
	}

	res := APIResponse{
		Message: message,
		Data: map[string]interface{}{
			"id":        user.ID,
			"full_name": user.FullName,
			"email":     user.Email,
			"weight":    user.Weight,
			"height":    user.Height,
			"bmi":       bmiVal,
		},
	}

	return c.JSON(http.StatusOK, res)
}
