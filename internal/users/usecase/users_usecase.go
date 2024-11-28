package usecase

import "github.com/liuhua1307/gin-read/internal/domain"

var _ domain.UserUseCase = &UserUseCase{}

// UserUseCase is a struct implementing the UserUseCase interface
type UserUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: userRepo}
}

func (u UserUseCase) CreateUser(user *domain.User) error {
	return u.userRepo.Create(user)
}

func (u UserUseCase) GetAllUsers() ([]domain.User, error) {
	return u.userRepo.FindAll()
}

func (u UserUseCase) GetUserByID(id int) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u UserUseCase) UpdateUser(user *domain.User) error {
	return u.userRepo.Update(user)
}

func (u UserUseCase) DeleteUser(id int) error {
	return u.userRepo.Delete(id)
}
