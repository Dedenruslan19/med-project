package workouts

type Workout struct {
	Name   string `json:"workout_name" gorm:"type:varchar(255);not null"`
	Goals  string `json:"goals" gorm:"type:text;not null"`
	ID     int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int64  `json:"user_id" gorm:"not null"`
}

type Exercise struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	WorkoutID int64  `json:"workout_id" gorm:"not null;index"`
	Name      string `json:"exercise_name" gorm:"type:varchar(255);not null"`
	Sets      string `json:"sets" gorm:"default:0"`
	Reps      string `json:"reps" gorm:"default:0"`
	Equipment string `json:"equipment" gorm:"type:varchar(255)"`
}
