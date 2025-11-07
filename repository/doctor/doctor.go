package doctor

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID             int64     `json:"id" gorm:"primaryKey"`
	FullName       string    `json:"full_name" gorm:"not null"`
	Email          string    `json:"email" gorm:"uniqueIndex;not null"`
	Password       string    `json:"-" gorm:"not null"`
	Specialization string    `json:"specialization" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	IsAvailable    bool      `json:"is_available" gorm:"default:true"`
}

type DoctorRepo interface {
	GetByID(id int64) (*Doctor, error)
	GetAll() ([]Doctor, error)
	GetByEmail(email string) (*Doctor, error)
	Create(doctor *Doctor) (int64, error)
}

type doctorRepository struct {
	logger *slog.Logger
	db     *gorm.DB
}

func NewDoctorRepository(logger *slog.Logger, db *gorm.DB) DoctorRepo {
	return &doctorRepository{
		logger: logger,
		db:     db,
	}
}

func (r *doctorRepository) GetByID(id int64) (*Doctor, error) {
	var doctor Doctor
	if err := r.db.First(&doctor, id).Error; err != nil {
		r.logger.Error("failed to get doctor by ID", slog.Any("error", err), slog.Int64("doctor_id", id))
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) GetAll() ([]Doctor, error) {
	var doctors []Doctor
	if err := r.db.Find(&doctors).Error; err != nil {
		r.logger.Error("failed to get all doctors", slog.Any("error", err))
		return nil, err
	}
	return doctors, nil
}

func (r *doctorRepository) GetByEmail(email string) (*Doctor, error) {
	var doctor Doctor
	if err := r.db.Where("email = ?", email).First(&doctor).Error; err != nil {
		r.logger.Error("failed to get doctor by email", slog.Any("error", err), slog.String("email", email))
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) Create(doctor *Doctor) (int64, error) {
	if err := r.db.Create(doctor).Error; err != nil {
		r.logger.Error("failed to create doctor", slog.Any("error", err), slog.String("email", doctor.Email))
		return 0, err
	}
	return doctor.ID, nil
}
