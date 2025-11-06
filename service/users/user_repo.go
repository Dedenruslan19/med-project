package users

type UserRepo interface {
	Create(user User) (int64, error)
	GetByEmail(email string) (User, error)
	FindByID(id int64) (User, error)
}
