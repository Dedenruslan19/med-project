package errors

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrWorkoutNotFound    = errors.New("workout not found")
	ErrInvalidAuthor      = errors.New("invalid author")
	ErrExerciseNotFound   = errors.New("exercise not found")
	ErrLogNotFound        = errors.New("log not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPass        = errors.New("invalid password")
	ErrHashFailed         = errors.New("failed to hash password")
	ErrJWTFailed          = errors.New("failed to generate JWT token")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrDoctorBusy         = errors.New("doctor is not available at the requested time")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrDoctorOnly         = errors.New("only doctors can perform this action")
)
