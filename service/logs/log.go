package logs

import (
	"time"
)

type ExerciseLog struct {
	ID         int64     `json:"id" gorm:"primaryKey"`
	ExerciseID int64     `json:"exercise_id" gorm:"not null"`
	UserID     int64     `json:"user_id" gorm:"not null"`
	SetCount   int64     `json:"set_count" gorm:"not null"`
	RepCount   int64     `json:"rep_count" gorm:"not null"`
	Weight     float64   `json:"weight" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
