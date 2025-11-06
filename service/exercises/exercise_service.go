package exercises

import (
	"errors"
	"log"
	"log/slog"

	"Dedenruslan19/med-project/repository/gemini"
	"Dedenruslan19/med-project/service/workouts"
)

type service struct {
	repo       ExerciseRepo
	logger     *slog.Logger
	workouts   workouts.Service
	geminiRepo gemini.Repository
}

type Service interface {
	GetExerciseByID(exerciseID int64) (*Exercise, error)
	CreateExercise(userID, workoutID int64, input ExerciseInput) (Exercise, error)
	DeleteExercise(userID, exerciseID int64) (*Exercise, error)
}

func NewService(logger *slog.Logger, repo ExerciseRepo, workoutService workouts.Service, gemRepo gemini.Repository) Service {
	return &service{
		logger:     logger,
		repo:       repo,
		workouts:   workoutService,
		geminiRepo: gemRepo,
	}
}

type ExerciseInput struct {
	WorkoutID   int64  `json:"workout_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"required"`
	Sets        string `json:"sets"`
	Reps        string `json:"reps"`
	Equipment   string `json:"equipment"`
}

var (
	ErrExerciseNotFound = errors.New("exercise not found")
	ErrInvalidAuthor    = errors.New("invalid author")
)

func (s *service) GetExerciseByID(exerciseID int64) (*Exercise, error) {
	exercise, err := s.repo.GetByID(exerciseID)
	if err != nil {
		s.logger.Error("Failed to get id exercise",
			slog.Any("error", err),
			slog.Int64("exercise_id", exerciseID),
		)
		return &Exercise{}, ErrExerciseNotFound
	}
	return exercise, nil
}

func (s *service) CreateExercise(userID, workoutID int64, input ExerciseInput) (Exercise, error) {
	workout, err := s.workouts.GetWorkoutByID(workoutID)
	if err != nil {
		s.logger.Error("Failed to get workout by ID",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return Exercise{}, err
	}

	if workout.UserID != userID {
		return Exercise{}, ErrInvalidAuthor
	}

	exercise := Exercise{
		Name:      input.Name,
		Equipment: input.Equipment,
		WorkoutID: workoutID,
		Sets:      input.Sets,
		Reps:      input.Reps,
	}

	id, err := s.repo.Create(&exercise)
	if err != nil {
		s.logger.Error("Failed to create exercise",
			slog.Any("error", err),
			slog.Int64("exercise_id", id),
		)
		return Exercise{}, err
	}

	exercise.ID = id
	return exercise, nil
}

func (s *service) DeleteExercise(userID, exerciseID int64) (*Exercise, error) {
	exercise, err := s.repo.GetByID(exerciseID)
	if err != nil {
		return nil, ErrExerciseNotFound
	}

	workout, err := s.workouts.GetWorkoutByID(exercise.WorkoutID)
	if err != nil {
		s.logger.Error("Failed to get workout for exercise",
			slog.Any("error", err),
			slog.Int64("exercise_id", exerciseID),
		)
		return nil, err
	}

	if workout.UserID != userID {
		return nil, ErrInvalidAuthor
	}

	if err := s.repo.Delete(exerciseID); err != nil {
		log.Println("failed to delete exercise:", err)
		return nil, err
	}

	return exercise, nil
}
