package repository

import (
	"github.com/liuhua1307/gin-read/internal/configs"
	"github.com/liuhua1307/gin-read/internal/domain"
	"github.com/liuhua1307/gin-read/internal/pkg/log"
	log2 "github.com/liuhua1307/gin-read/pkg/log"
	"gorm.io/gorm"
)

var _ domain.UserRepository = &UserMySQLRepository{}

// UserMySQLRepository is a struct implementing the UserRepository interface
type UserMySQLRepository struct {
	d *configs.DataInstance
}

func NewUserMySQLRepository(d *configs.DataInstance) *UserMySQLRepository {
	return &UserMySQLRepository{d: d}
}

func (u UserMySQLRepository) Create(user *domain.User) error {
	return u.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			log.Log().Error("Create user failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "User", Value: user.Name})
			return err
		}
		return nil
	})
}

func (u UserMySQLRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	if err := u.d.DB.Find(&users).Error; err != nil {
		log.Log().Error("Find all users failed ", log2.Field{Key: "error", Value: err})
		return nil, err
	}
	return users, nil
}

func (u UserMySQLRepository) FindByID(id int) (*domain.User, error) {
	var user domain.User
	if err := u.d.DB.First(&user, id).Error; err != nil {
		log.Log().Error("Find user by id failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "ID", Value: id})
		return nil, err
	}
	return &user, nil
}

func (u UserMySQLRepository) Update(user *domain.User) error {
	return u.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			log.Log().Error("Update user failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "User", Value: user.Name})
			return err
		}
		return nil
	})
}

func (u UserMySQLRepository) Delete(id int) error {
	return u.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.User{}, id).Error; err != nil {
			log.Log().Error("Delete user failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "ID", Value: id})
			return err
		}
		return nil
	})
}
