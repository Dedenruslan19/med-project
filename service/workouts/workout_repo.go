package workouts

type WorkoutRepo interface {
	GetAll() ([]Workout, error)
	GetByID(workoutID int64) (*Workout, error)
	Create(workout *Workout) (int64, error)
	CreateExercise(exercise *Exercise) error
	Update(workout *Workout) error
	GetOwnerID(workoutID int64) (int64, error)
	Delete(workoutID int64) error
}
