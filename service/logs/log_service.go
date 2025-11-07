package logs

import (
	"log/slog"

	errs "Dedenruslan19/med-project/service/errors"
)

type service struct {
	repo   LogRepo
	logger *slog.Logger
}

type Service interface {
	CreateLog(userID int64, input LogInput) (ExerciseLog, error)
	GetAllLogs(userID int64) ([]ExerciseLog, error)
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
		s.logger.Error("failed to create log",
			slog.Any("error", err),
			slog.Int64("user_id", userID),
		)
		return ExerciseLog{}, err
	}

	log.ID = id
	return log, nil
}

func (s *service) GetAllLogs(userID int64) ([]ExerciseLog, error) {
	logs, err := s.repo.GetByUserID(userID)
	if err != nil {
		s.logger.Error("failed to get all logs",
			slog.Any("error", err),
		)
		return nil, err
	}

	if len(logs) == 0 {
		return nil, errs.ErrLogNotFound
	}

	return logs, nil
}
