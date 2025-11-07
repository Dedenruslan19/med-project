package doctors

import (
	"time"
)

type Doctor struct {
	ID             int64     `json:"id"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	Specialization string    `json:"specialization"`
	CreatedAt      time.Time `json:"created_at"`
	IsAvailable    bool      `json:"is_available"`
}
