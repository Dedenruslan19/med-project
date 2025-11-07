package doctors

import (
	"Dedenruslan19/med-project/repository/doctor"
)

type DoctorRepo interface {
	GetByID(id int64) (*doctor.Doctor, error)
	GetAll() ([]doctor.Doctor, error)
	GetByEmail(email string) (*doctor.Doctor, error)
	Create(doctor *doctor.Doctor) (int64, error)
}
