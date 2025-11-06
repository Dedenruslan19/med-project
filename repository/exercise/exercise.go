package exercise

import (
	"errors"
	"log/slog"

	service "Dedenruslan19/med-project/service/exercises"

	"gorm.io/gorm"
)

type exerciseRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewExerciseRepo(db *gorm.DB, logger *slog.Logger) service.ExerciseRepo {
	return &exerciseRepo{db: db, logger: logger}
}

func (r *exerciseRepo) GetByID(exerciseID int64) (*service.Exercise, error) {
	var exercise service.Exercise
	if err := r.db.First(&exercise, exerciseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Error("exercise not found", "exercise_id", exerciseID)
			return nil, errors.New("exercise not found")
		}
		r.logger.Error("failed to get exercise by id", "exercise_id", exerciseID, "error", err)
		return nil, err
	}
	return &exercise, nil
}

func (r *exerciseRepo) Create(exercise *service.Exercise) (int64, error) {
	newExercise := service.Exercise{
		Name:      exercise.Name,
		Equipment: exercise.Equipment,
		WorkoutID: exercise.WorkoutID,
		Sets:      exercise.Sets,
		Reps:      exercise.Reps,
	}

	if err := r.db.Create(&newExercise).Error; err != nil {
		r.logger.Error("failed to create exercise", "name", exercise.Name, "workout_id", exercise.WorkoutID, "error", err)
		return 0, err
	}
	return newExercise.ID, nil
}

func (r *exerciseRepo) Delete(exerciseID int64) error {
	result := r.db.Delete(&service.Exercise{}, exerciseID)
	if result.Error != nil {
		r.logger.Error("failed to delete exercise", "exercise_id", exerciseID, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Error("exercise not found when deleting", "exercise_id", exerciseID)
		return service.ErrExerciseNotFound
	}
	return nil
}
