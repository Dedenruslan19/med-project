package users

import (
	"Dedenruslan19/med-project/repository/rapidAPI/bmi"
	"Dedenruslan19/med-project/util"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo   UserRepo
	logger *slog.Logger
	bmi    bmi.Repository
}

type Service interface {
	Register(user User) (int64, error)
	Login(email, passwordhash string) (User, error)
	GetUserByID(userID int64) (User, error)
	CalculateBMI(weight, height float64) (float64, bool, error)
}

func NewService(logger *slog.Logger, repo UserRepo, bmiRepo bmi.Repository) Service {
	return &service{
		logger: logger,
		repo:   repo,
		bmi:    bmiRepo,
	}
}

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPass        = errors.New("invalid password")
	ErrHashFailed         = errors.New("failed to hash password")
	ErrJWTFailed          = errors.New("failed to generate JWT token")
	ErrEmailAlreadyExists = errors.New("email already registered")
)

func (s *service) Register(user User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password",
			slog.String("emai;", user.Email),
			slog.Any("error", err),
		)
		return 0, ErrHashFailed
	}

	user.Password = string(hashedPassword)

	id, err := s.repo.Create(user)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return 0, ErrEmailAlreadyExists
		}
		s.logger.Error("Failed to create user in DB",
			slog.String("email", user.Email),
			slog.Any("error", err),
		)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (s *service) Login(email, password string) (User, error) {
	user, err := s.repo.GetByEmail(email)

	if err != nil {
		s.logger.Warn("User not found ",
			slog.String("email", email),
			slog.Any("error", err),
		)
		return User{}, ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Warn("Invalid password attempt", slog.String("email", email))
		return User{}, ErrInvalidPass
	}

	token, err := util.GenerateJWT(user.ID, user.Email)
	if err != nil {
		s.logger.Error("Failed to generate JWT for user",
			slog.Int64("user_id", user.ID),
			slog.Any("error", err),
		)
		return User{}, ErrJWTFailed
	}
	user.Token = token

	return user, nil
}

func (s *service) GetUserByID(userID int64) (User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		s.logger.Error("Failed to fetch user by ID",
			slog.Int64("user_id", userID),
			slog.Any("error", err),
		)
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func (s *service) CalculateBMI(weight, height float64) (float64, bool, error) {
	var usedCallback bool

	bmiVal, err := s.bmi.CalculateBMI(weight, height, func(w, h float64) float64 {
		usedCallback = true
		return bmi.DefaultBMICallback(w, h)
	})

	if err != nil {
		return 0, usedCallback, err
	}

	return bmiVal, usedCallback, nil
}
