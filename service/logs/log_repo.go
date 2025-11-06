package logs

type LogRepo interface {
	Create(log *ExerciseLog) (int64, error)
	GetByUserID(userID int64) ([]ExerciseLog, error)
	GetByID(logID int64) (*ExerciseLog, error)
	GetAll() ([]ExerciseLog, error)
	Delete(logID int64) error
}
