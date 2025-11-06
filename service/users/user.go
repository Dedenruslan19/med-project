package users

type User struct {
	ID       int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName string  `gorm:"type:varchar(255);not null" json:"full_name" validate:"required,min=2,max=255"`
	Email    string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string  `gorm:"type:varchar(255);not null" json:"-"`
	Weight   float64 `gorm:"type:decimal(5,2);not null" json:"weight" validate:"required,gt=0"`
	Height   float64 `gorm:"not null" json:"height" validate:"required,gt=0"`
	Token    string  `gorm:"-" json:"token,omitempty"`
}
