package user

import (
	"errors"
	"log/slog"
	"strings"

	service "Dedenruslan19/med-project/service/users"

	"gorm.io/gorm"
)

type userRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepo(db *gorm.DB, logger *slog.Logger) service.UserRepo {
	return &userRepo{db: db, logger: logger}
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
)

func (r *userRepo) Create(u service.User) (int64, error) {
	if err := r.db.Create(&u).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "Duplicate entry") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {

			r.logger.Error("failed to create user - email already exists",
				"email", u.Email,
				"error", err)
			return 0, service.ErrEmailAlreadyExists
		}

		r.logger.Error("failed to create user",
			"email", u.Email,
			"error", err)
		return 0, err
	}

	return u.ID, nil
}

func (r *userRepo) GetByEmail(email string) (service.User, error) {
	var u service.User
	err := r.db.Where("email = ?", email).First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error("failed to get user by email - not found",
			"email", email,
			"error", err)
		return u, ErrUserNotFound
	}

	if err != nil {
		r.logger.Error("failed to get user by email",
			"email", email,
			"error", err)
	}

	return u, err
}

func (r *userRepo) FindByID(id int64) (service.User, error) {
	var u service.User
	if err := r.db.First(&u, id).Error; err != nil {
		r.logger.Error("failed to find user by id",
			"user_id", id,
			"error", err)
		return service.User{}, err
	}
	return u, nil
}
