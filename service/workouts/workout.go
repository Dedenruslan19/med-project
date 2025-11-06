package workouts

type Workout struct {
	Name   string `json:"workout_name" gorm:"type:varchar(255);not null"`
	Goals  string `json:"goals" gorm:"type:text;not null"`
	ID     int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int64  `json:"user_id" gorm:"not null"`
}
