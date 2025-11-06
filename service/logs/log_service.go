package logs

import (
	"errors"
	"log/slog"
)

type service struct {
	repo   LogRepo
	logger *slog.Logger
}

type Service interface {
	CreateLog(userID int64, input LogInput) (ExerciseLog, error)
	GetAllLogs() ([]ExerciseLog, error)
}

func NewService(logger *slog.Logger, repo LogRepo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

type LogInput struct {
	ExerciseID int64   `json:"exercise_id"`
	Weight     float64 `json:"weight"`
	RepCount   int64   `json:"rep_count"`
	SetCount   int64   `json:"set_count"`
}

var (
	ErrLogNotFound = errors.New("log not found")
)

func (s *service) CreateLog(userID int64, input LogInput) (ExerciseLog, error) {
	log := ExerciseLog{
		UserID:     userID,
		ExerciseID: input.ExerciseID,
		Weight:     input.Weight,
		RepCount:   input.RepCount,
		SetCount:   input.SetCount,
	}

	id, err := s.repo.Create(&log)
	if err != nil {
		s.logger.Error("Failed to create log",
			slog.Any("error", err),
			slog.Int64("user_id", userID),
		)
		return ExerciseLog{}, errors.New("failed to create log")
	}

	log.ID = id
	return log, nil
}

func (s *service) GetAllLogs() ([]ExerciseLog, error) {
	logs, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get all logs",
			slog.Any("error", err),
		)
		return nil, errors.New("failed to fetch logs")
	}

	if len(logs) == 0 {
		return nil, ErrLogNotFound
	}

	return logs, nil
}
