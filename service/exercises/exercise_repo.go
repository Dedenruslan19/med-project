package exercises

type ExerciseRepo interface {
	GetByID(exerciseID int64) (*Exercise, error)
	Create(exercise *Exercise) (int64, error)
	Delete(exerciseID int64) error
}
