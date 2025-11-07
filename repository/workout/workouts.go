package workout

import (
	"errors"
	"fmt"
	"log/slog"

	errs "Dedenruslan19/med-project/service/errors"
	"Dedenruslan19/med-project/service/workouts"

	"gorm.io/gorm"
)

type workoutRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewWorkoutRepo(db *gorm.DB, logger *slog.Logger) workouts.WorkoutRepo {
	return &workoutRepo{db: db, logger: logger}
}

func (r *workoutRepo) GetAll() ([]workouts.Workout, error) {
	var workoutList []workouts.Workout
	err := r.db.Find(&workoutList).Error
	if err != nil {
		r.logger.Error("failed to get all workouts",
			slog.Any("error", err))
		return nil, err
	}
	return workoutList, nil
}

func (r *workoutRepo) GetByID(workoutID int64) (*workouts.Workout, error) {
	var workout workouts.Workout
	err := r.db.First(&workout, workoutID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("workout not found")
		}
		r.logger.Error("failed to get workout by ID",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID))
		return nil, err
	}
	return &workout, nil
}

func (r *workoutRepo) Create(workout *workouts.Workout) (int64, error) {
	if err := r.db.Create(workout).Error; err != nil {
		r.logger.Error("failed to create workout",
			slog.Any("error", err), slog.Int64("user_id", workout.UserID))
		return 0, err
	}
	return workout.ID, nil
}

func (r *workoutRepo) CreateExercise(exercise *workouts.Exercise) error {
	if err := r.db.Create(exercise).Error; err != nil {
		r.logger.Error("failed to create exercise",
			slog.Any("error", err),
			slog.Int64("workout_id", exercise.WorkoutID),
			slog.String("exercise_name", exercise.Name))
		return err
	}
	return nil
}

func (r *workoutRepo) CreateExercises(exercises []workouts.Exercise, workoutID int64) error {
	tx := r.db.Begin()
	for _, ex := range exercises {
		ex.WorkoutID = workoutID
		if err := tx.Create(&ex).Error; err != nil {
			tx.Rollback()
			r.logger.Error("failed to create exercises",
				slog.Any("error", err))
			return err
		}
	}
	return tx.Commit().Error
}

func (r *workoutRepo) Update(workout *workouts.Workout) error {
	result := r.db.Model(&workouts.Workout{}).
		Where("id = ? AND user_id = ?", workout.ID, workout.UserID).
		Updates(map[string]interface{}{
			"name":  workout.Name,
			"goals": workout.Goals,
		})

	if result.Error != nil {
		r.logger.Error("failed to update workout",
			slog.Any("error", result.Error),
			slog.Int64("workout_id", workout.ID))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("workout not found or not owned by user")
	}
	return nil
}

func (r *workoutRepo) GetOwnerID(workoutID int64) (int64, error) {
	var workout workouts.Workout
	err := r.db.Select("user_id").Where("id = ?", workoutID).First(&workout).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errs.ErrWorkoutNotFound
		}
		r.logger.Error("failed to get workout owner",
			slog.Any("error", err),
			slog.Int64("workout_id", workoutID))
		return 0, err
	}
	return workout.UserID, nil
}

func (r *workoutRepo) Delete(workoutID int64) error {
	result := r.db.Delete(&workouts.Workout{}, workoutID)
	if result.Error != nil {
		r.logger.Error("failed to delete workout",
			slog.Any("error", result.Error),
			slog.Int64("workout_id", workoutID))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrWorkoutNotFound
	}
	return nil
}
