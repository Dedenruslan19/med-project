package workouts

import (
	"Dedenruslan19/med-project/repository/gemini"
	errs "Dedenruslan19/med-project/service/errors"
	"encoding/json"
	"fmt"
	"log/slog"
)

type service struct {
	repo   WorkoutRepo
	logger *slog.Logger
	gemini gemini.Repository
}

type Service interface {
	GetAllWorkouts() ([]Workout, error)
	GetWorkoutByID(userID, workoutID int64) (*Workout, error)
	CreateWorkout(userID int64, input SaveWorkoutRequest) (WorkoutWithExercises, error)
	DeleteWorkout(userID, workoutID int64) error
	PreviewWorkout(userID int64, input PreviewWorkoutRequest) (WorkoutWithExercises, error)
}

func NewService(logger *slog.Logger, repo WorkoutRepo, gemini gemini.Repository) Service {
	return &service{
		logger: logger,
		repo:   repo,
		gemini: gemini,
	}
}

func parseGeminiOutput(geminiOutput string) ([]GeneratedExercise, error) {
	var raw []map[string]interface{}

	if err := json.Unmarshal([]byte(geminiOutput), &raw); err != nil {
		return nil, err
	}

	var exercises []GeneratedExercise
	for _, r := range raw {
		equipment := "none"
		if v, ok := r["equipment"]; ok && v != nil && fmt.Sprintf("%v", v) != "" {
			equipment = fmt.Sprintf("%v", v)
		}

		exercises = append(exercises, GeneratedExercise{
			Name:      fmt.Sprintf("%v", r["name"]),
			Sets:      fmt.Sprintf("%v", r["sets"]),
			Reps:      fmt.Sprintf("%v", r["reps"]),
			Equipment: equipment,
		})
	}

	return exercises, nil
}

type GeneratedExercise struct {
	Name      string `json:"name"`
	Sets      string `json:"sets"`
	Reps      string `json:"reps"`
	Equipment string `json:"equipment"`
}

type WorkoutWithExercises struct {
	Workout   Workout             `json:"workout_name"`
	Exercises []GeneratedExercise `json:"exercises"`
}

type PreviewWorkoutRequest struct {
	WorkoutName string `json:"workout_name" validate:"required,min=2,max=255"`
	Goals       string `json:"goals" validate:"required"`
}

type SaveWorkoutRequest struct {
	WorkoutName WorkoutData         `json:"workout_name" validate:"required"`
	Exercises   []GeneratedExercise `json:"exercises" validate:"required,dive"`
}

type WorkoutData struct {
	WorkoutName string `json:"workout_name" validate:"required,min=2,max=255"`
	Goals       string `json:"goals" validate:"required"`
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
}

func (s *service) GetAllWorkouts() ([]Workout, error) {
	workouts, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("failed to get workouts",
			slog.Any("error", err),
		)
		return nil, err
	}
	return workouts, nil
}

func (s *service) GetWorkoutByID(userID, workoutID int64) (*Workout, error) {
	workout, err := s.repo.GetByID(workoutID)
	if err != nil {
		s.logger.Error("failed to get id workout",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return nil, errs.ErrWorkoutNotFound
	}

	if workout.UserID != userID {
		return nil, errs.ErrInvalidAuthor
	}

	return workout, nil
}

func (s *service) PreviewWorkout(userID int64, input PreviewWorkoutRequest) (WorkoutWithExercises, error) {
	// Preview workout, doesnt save to DB
	workout := Workout{
		Name:   input.WorkoutName,
		Goals:  input.Goals,
		UserID: userID,
		ID:     0,
	}

	// Generate exercises via Gemini
	raw, err := s.gemini.GenerateExercises(gemini.GenerateRequest{
		WorkoutID: 0,
		Target:    workout.Name,
		Goal:      workout.Goals,
		Equipment: []string{"bodyweight"},
	})
	if err != nil {
		s.logger.Error("gemini error", slog.Any("error", err))
		return WorkoutWithExercises{Workout: workout, Exercises: []GeneratedExercise{}}, nil
	}

	exercises, err := parseGeminiOutput(raw)
	if err != nil {
		s.logger.Error("failed to parse Gemini output", slog.Any("error", err))
		return WorkoutWithExercises{Workout: workout, Exercises: []GeneratedExercise{}}, nil
	}

	return WorkoutWithExercises{
		Workout:   workout,
		Exercises: exercises,
	}, nil
}

func (s *service) CreateWorkout(userID int64, input SaveWorkoutRequest) (WorkoutWithExercises, error) {
	workout := Workout{
		Name:   input.WorkoutName.WorkoutName,
		Goals:  input.WorkoutName.Goals,
		UserID: userID,
	}
	id, err := s.repo.Create(&workout)
	if err != nil {
		s.logger.Error("failed to create workout", slog.Any("error", err))
		return WorkoutWithExercises{}, err
	}
	workout.ID = id

	var exercises []GeneratedExercise
	for _, ex := range input.Exercises {
		exercise := Exercise{
			WorkoutID: workout.ID,
			Name:      ex.Name,
			Sets:      ex.Sets,
			Reps:      ex.Reps,
			Equipment: ex.Equipment,
		}
		if err := s.repo.CreateExercise(&exercise); err != nil {
			s.logger.Error("failed to create exercise", slog.Any("error", err))
			continue
		}
		exercises = append(exercises, GeneratedExercise{
			Name:      exercise.Name,
			Sets:      exercise.Sets,
			Reps:      exercise.Reps,
			Equipment: exercise.Equipment,
		})
	}

	return WorkoutWithExercises{
		Workout:   workout,
		Exercises: exercises,
	}, nil
}

func (s *service) DeleteWorkout(userID, workoutID int64) error {
	ownerID, err := s.repo.GetOwnerID(workoutID)
	if err != nil {
		s.logger.Error("failed to get workout owner",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return errs.ErrWorkoutNotFound
	}

	if ownerID != userID {
		return errs.ErrInvalidAuthor
	}

	if err := s.repo.Delete(workoutID); err != nil {
		s.logger.Error("failed to delete workout",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return err
	}

	return nil
}
