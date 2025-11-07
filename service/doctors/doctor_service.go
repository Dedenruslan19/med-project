package doctors

import (
	"log/slog"

	"Dedenruslan19/med-project/repository/doctor"
	errs "Dedenruslan19/med-project/service/errors"

	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo   DoctorRepo
	logger *slog.Logger
}

type Service interface {
	GetAll() ([]Doctor, error)
	GetByID(id int64) (*Doctor, error)
	Register(fullName, email, password, specialization string) (*Doctor, error)
	Login(email, password string) (*Doctor, error)
}

func NewService(logger *slog.Logger, repo DoctorRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) GetAll() ([]Doctor, error) {
	doctorsRepo, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("failed to get all doctors", slog.Any("error", err))
		return nil, err
	}

	doctors := make([]Doctor, len(doctorsRepo))
	for i, d := range doctorsRepo {
		doctors[i] = Doctor{
			ID:             d.ID,
			FullName:       d.FullName,
			Email:          d.Email,
			Specialization: d.Specialization,
			CreatedAt:      d.CreatedAt,
			IsAvailable:    d.IsAvailable,
		}
	}

	return doctors, nil
}

func (s *service) GetByID(id int64) (*Doctor, error) {
	doctorRepo, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("failed to get doctor by ID", slog.Any("error", err), slog.Int64("doctor_id", id))
		return nil, err
	}

	doctor := &Doctor{
		ID:             doctorRepo.ID,
		FullName:       doctorRepo.FullName,
		Email:          doctorRepo.Email,
		Specialization: doctorRepo.Specialization,
		CreatedAt:      doctorRepo.CreatedAt,
		IsAvailable:    doctorRepo.IsAvailable,
	}

	return doctor, nil
}

func (s *service) Login(email, password string) (*Doctor, error) {
	doctorRepo, err := s.repo.GetByEmail(email)
	if err != nil {
		s.logger.Error("doctor not found", slog.String("email", email))
		return nil, errs.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(doctorRepo.Password), []byte(password))
	if err != nil {
		s.logger.Error("invalid password", slog.String("email", email))
		return nil, errs.ErrInvalidCredentials
	}

	doctor := &Doctor{
		ID:             doctorRepo.ID,
		FullName:       doctorRepo.FullName,
		Email:          doctorRepo.Email,
		Specialization: doctorRepo.Specialization,
		CreatedAt:      doctorRepo.CreatedAt,
		IsAvailable:    doctorRepo.IsAvailable,
	}

	return doctor, nil
}

func (s *service) Register(fullName, email, password, specialization string) (*Doctor, error) {
	// Check if email already exists
	existingDoctor, _ := s.repo.GetByEmail(email)
	if existingDoctor != nil {
		s.logger.Error("email already registered", slog.String("email", email))
		return nil, errs.ErrEmailAlreadyExists
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash password", slog.Any("error", err))
		return nil, errs.ErrHashFailed
	}

	// Create doctor repository entity
	newDoctor := &doctor.Doctor{
		FullName:       fullName,
		Email:          email,
		Password:       string(hashedPassword),
		Specialization: specialization,
		IsAvailable:    true,
	}

	id, err := s.repo.Create(newDoctor)
	if err != nil {
		s.logger.Error("failed to create doctor", slog.Any("error", err))
		return nil, err
	}

	doctorResponse := &Doctor{
		ID:             id,
		FullName:       fullName,
		Email:          email,
		Specialization: specialization,
		IsAvailable:    true,
	}

	return doctorResponse, nil
}
