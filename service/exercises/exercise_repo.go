package exercises

type ExerciseRepo interface {
	GetByID(exerciseID int64) (*Exercise, error)
	GetByWorkoutID(workoutID int64) ([]Exercise, error)
	Create(exercise *Exercise) (int64, error)
	Update(exercise *Exercise) error
	Delete(exerciseID int64) error
}
