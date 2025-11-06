package exercises

type Exercise struct {
	Sets      string `json:"sets" gorm:"default:0"`
	Reps      string `json:"reps" gorm:"default:0"`
	Name      string `json:"exercise_name" gorm:"type:varchar(255);not null"`
	Equipment string `json:"equipment" gorm:"type:varchar(255)"`
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	WorkoutID int64  `json:"workout_id" gorm:"not null;index"`
}
