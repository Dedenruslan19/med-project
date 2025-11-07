package exercises

import (
	"log/slog"

	"Dedenruslan19/med-project/repository/gemini"
	errs "Dedenruslan19/med-project/service/errors"
	"Dedenruslan19/med-project/service/workouts"
)

type service struct {
	repo       ExerciseRepo
	logger     *slog.Logger
	workouts   workouts.Service
	geminiRepo gemini.Repository
}

type Service interface {
	GetExerciseByID(userID, exerciseID int64) (*Exercise, error)
	GetExercisesByWorkoutID(userID, workoutID int64) ([]Exercise, error)
	CreateExercise(userID, workoutID int64, input ExerciseInput) (Exercise, error)
	UpdateExercise(userID, exerciseID int64, input ExerciseInput) (Exercise, error)
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

func (s *service) GetExerciseByID(userID, exerciseID int64) (*Exercise, error) {
	exercise, err := s.repo.GetByID(exerciseID)
	if err != nil {
		s.logger.Error("Failed to get exercise by ID",
			slog.Any("error", err),
			slog.Int64("exercise_id", exerciseID),
		)
		return nil, errs.ErrExerciseNotFound
	}
	workout, err := s.workouts.GetWorkoutByID(userID, exercise.WorkoutID)
	if err != nil {
		return nil, err // akan mengembalikan ErrWorkoutNotFound atau ErrInvalidAuthor
	}

	if workout.UserID != userID {
		return nil, errs.ErrInvalidAuthor
	}

	return exercise, nil
}

func (s *service) GetExercisesByWorkoutID(userID, workoutID int64) ([]Exercise, error) {
	workout, err := s.workouts.GetWorkoutByID(userID, workoutID)
	if err != nil {
		return nil, err
	}

	if workout.UserID != userID {
		return nil, errs.ErrInvalidAuthor
	}

	exercises, err := s.repo.GetByWorkoutID(workoutID)
	if err != nil {
		s.logger.Error("Failed to get exercises by workout_id",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return nil, err
	}

	return exercises, nil
}

func (s *service) CreateExercise(userID, workoutID int64, input ExerciseInput) (Exercise, error) {
	workout, err := s.workouts.GetWorkoutByID(userID, workoutID)
	if err != nil {
		return Exercise{}, err
	}

	if workout.UserID != userID {
		return Exercise{}, errs.ErrInvalidAuthor
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
		s.logger.Error("failed to create exercise",
			slog.Any("error", err),
			slog.Int64("exercise_id", id),
		)
		return Exercise{}, err
	}

	exercise.ID = id
	return exercise, nil
}

func (s *service) UpdateExercise(userID, exerciseID int64, input ExerciseInput) (Exercise, error) {
	exercise, err := s.repo.GetByID(exerciseID)
	if err != nil {
		s.logger.Error("failed to get exercise by ID",
			slog.Any("error", err),
			slog.Int64("exercise_id", exerciseID),
		)
		return Exercise{}, errs.ErrExerciseNotFound
	}

	// Validasi ownership - cek apakah user adalah pemilik workout
	workout, err := s.workouts.GetWorkoutByID(userID, exercise.WorkoutID)
	if err != nil {
		s.logger.Error("failed to get workout by ID",
			slog.Any("error", err),
			slog.Int64("workout_id", exercise.WorkoutID),
		)
		return Exercise{}, err
	}

	if workout.UserID != userID {
		return Exercise{}, errs.ErrInvalidAuthor
	}
	exercise.Name = input.Name
	exercise.Sets = input.Sets
	exercise.Reps = input.Reps
	exercise.Equipment = input.Equipment

	// Save to database
	if err := s.repo.Update(exercise); err != nil {
		s.logger.Error("failed to update exercise",
			slog.Any("error", err),
			slog.Int64("exercise_id", exerciseID),
		)
		return Exercise{}, err
	}

	return *exercise, nil
}

func (s *service) DeleteExercise(userID, exerciseID int64) (*Exercise, error) {
	exercise, err := s.repo.GetByID(exerciseID)
	if err != nil {
		return nil, errs.ErrExerciseNotFound
	}

	workout, err := s.workouts.GetWorkoutByID(userID, exercise.WorkoutID)
	if err != nil {
		return nil, err
	}

	if workout.UserID != userID {
		return nil, errs.ErrInvalidAuthor
	}

	if err := s.repo.Delete(exerciseID); err != nil {
		s.logger.Error("failed to delete exercise",
			slog.Any("error", err),
			slog.Int64("exercise_id", exerciseID),
		)
		return nil, err
	}

	return exercise, nil
}
