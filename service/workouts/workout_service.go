package workouts

import (
	"Dedenruslan19/med-project/repository/gemini"
	"encoding/json"
	"errors"
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
	GetWorkoutByID(workoutID int64) (*Workout, error)
	CreateWorkout(userID int64, input WorkoutInput) (WorkoutWithExercises, error)
	UpdateWorkout(userID, workoutID int64, input WorkoutInput) (Workout, error)
	DeleteWorkout(userID, workoutID int64) error
}

func NewService(logger *slog.Logger, repo WorkoutRepo, gemini gemini.Repository) Service {
	return &service{
		logger: logger,
		repo:   repo,
		gemini: gemini,
	}
}

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrWorkoutNotFound = errors.New("workout not found")
	ErrInvalidAuthor   = errors.New("invalid author")
)

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
	Workout   Workout             `json:"workout"`
	Exercises []GeneratedExercise `json:"exercises"`
}

type WorkoutInput struct {
	Name      string              `json:"name" validate:"required,min=2,max=255"`
	Goals     string              `json:"goals" validate:"required"`
	Exercises []GeneratedExercise `json:"exercises"`
}

func (s *service) GetAllWorkouts() ([]Workout, error) {
	workouts, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Failed to get workouts",
			slog.Any("error", err),
		)
		return nil, err
	}
	return workouts, nil
}

func (s *service) GetWorkoutByID(workoutID int64) (*Workout, error) {
	workout, err := s.repo.GetByID(workoutID)
	if err != nil {
		s.logger.Error("Failed to get id workout",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return nil, ErrWorkoutNotFound
	}
	return workout, nil
}

func (s *service) CreateWorkout(userID int64, input WorkoutInput) (WorkoutWithExercises, error) {
	workout := Workout{
		Name:   input.Name,
		Goals:  input.Goals,
		UserID: userID,
	}

	id, err := s.repo.Create(&workout)
	if err != nil {
		s.logger.Error("Failed to create workout",
			slog.Any("error", err),
		)
		return WorkoutWithExercises{}, err
	}
	workout.ID = id

	// Generate exercises using Gemini
	generateReq := gemini.GenerateRequest{
		WorkoutID: workout.ID,
		Target:    workout.Name,
		Goal:      workout.Goals,
		Equipment: []string{"bodyweight"},
	}

	raw, err := s.gemini.GenerateExercises(generateReq)
	if err != nil {
		s.logger.Error("Gemini error", slog.Any("error", err))
		return WorkoutWithExercises{
			Workout:   workout,
			Exercises: []GeneratedExercise{},
		}, nil
	}

	exercises, err := parseGeminiOutput(raw)
	if err != nil {
		s.logger.Error("Failed to parse Gemini output", slog.Any("error", err))
		return WorkoutWithExercises{
			Workout:   workout,
			Exercises: []GeneratedExercise{},
		}, nil
	}

	return WorkoutWithExercises{
		Workout:   workout,
		Exercises: exercises,
	}, nil
}

func (s *service) UpdateWorkout(userID, workoutID int64, input WorkoutInput) (Workout, error) {
	ownerID, err := s.repo.GetOwnerID(workoutID)
	if err != nil {
		s.logger.Error("Failed to get workout owner",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return Workout{}, ErrWorkoutNotFound
	}

	if ownerID != userID {
		return Workout{}, ErrInvalidAuthor
	}

	workout := Workout{
		ID:     workoutID,
		Name:   input.Name,
		Goals:  input.Goals,
		UserID: userID,
	}

	if err = s.repo.Update(&workout); err != nil {
		s.logger.Error("Failed to update workout",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return Workout{}, err
	}

	updatedWorkout, err := s.repo.GetByID(workoutID)
	if err != nil {
		s.logger.Error("Failed to get updated workout",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return Workout{}, err
	}

	return *updatedWorkout, nil
}

func (s *service) DeleteWorkout(userID, workoutID int64) error {
	ownerID, err := s.repo.GetOwnerID(workoutID)
	if err != nil {
		s.logger.Error("Failed to get workout owner",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return ErrWorkoutNotFound
	}

	if ownerID != userID {
		return ErrInvalidAuthor
	}

	if err := s.repo.Delete(workoutID); err != nil {
		s.logger.Error("Failed to delete workout",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID),
		)
		return err
	}

	return nil
}
