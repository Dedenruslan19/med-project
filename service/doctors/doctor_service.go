package doctors

import (
	"Dedenruslan19/med-project/service/appointments"
	"log/slog"
)

type service struct {
	repo               ExerciseRepo
	logger             *slog.Logger
	appointmentService appointments.Service
}

type Service interface {
}

func NewService(logger *slog.Logger, repo ExerciseRepo, appointmentService appointments.Service) Service {
	return &service{
		logger:             logger,
		repo:               repo,
		appointmentService: appointmentService,
	}
}
