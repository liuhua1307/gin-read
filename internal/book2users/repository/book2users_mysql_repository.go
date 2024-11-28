package repository

import (
	"github.com/liuhua1307/gin-read/internal/configs"
	"github.com/liuhua1307/gin-read/internal/domain"
	"github.com/liuhua1307/gin-read/internal/pkg/log"
	log2 "github.com/liuhua1307/gin-read/pkg/log"
	"gorm.io/gorm"
)

var _ domain.Book2UsersRepository = &Book2UsersMySQLRepository{}

type Book2UsersMySQLRepository struct {
	d *configs.DataInstance
}

func NewBook2UsersMySQLRepository(d *configs.DataInstance) *Book2UsersMySQLRepository {
	return &Book2UsersMySQLRepository{d: d}
}

func (b *Book2UsersMySQLRepository) StartTransaction() *gorm.DB {
	return b.d.DB.Begin()
}

func (b *Book2UsersMySQLRepository) FindByUserID(userId int) (*domain.BooksLeaned, error) {
	var book2users domain.BooksLeaned
	if err := b.d.DB.Where("user_id = ?", userId).First(&book2users).Error; err != nil {
		log.Log().Error("Find book2users by user id failed", log2.Field{
			Key:   "err",
			Value: err,
		}, log2.Field{
			Key:   "userId",
			Value: userId,
		})
		return nil, err
	}
	return &book2users, nil
}

func (b *Book2UsersMySQLRepository) Create(book2users *domain.BooksLeaned) error {
	return b.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(book2users).Error; err != nil {
			log.Log().Error("Create book2users failed", log2.Field{
				Key:   "err",
				Value: err,
			}, log2.Field{
				Key:   "book2users",
				Value: book2users,
			})
			return err
		}
		return nil
	})
}

func (b *Book2UsersMySQLRepository) FindAll() ([]domain.BooksLeaned, error) {
	var book2users []domain.BooksLeaned
	if err := b.d.DB.Find(&book2users).Error; err != nil {
		log.Log().Error("Find all book2users failed", log2.Field{
			Key:   "err",
			Value: err,
		})
		return nil, err
	}
	return book2users, nil
}

func (b *Book2UsersMySQLRepository) FindByID(id int) (*domain.BooksLeaned, error) {
	var book2users domain.BooksLeaned
	if err := b.d.DB.First(&book2users, id).Error; err != nil {
		log.Log().Error("Find book2users by id failed", log2.Field{
			Key:   "err",
			Value: err,
		}, log2.Field{
			Key:   "id",
			Value: id,
		})
		return nil, err
	}
	return &book2users, nil
}

func (b *Book2UsersMySQLRepository) Update(book2users *domain.BooksLeaned) error {
	return b.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(book2users).Error; err != nil {
			log.Log().Error("Update book2users failed", log2.Field{
				Key:   "err",
				Value: err,
			}, log2.Field{
				Key:   "book2users",
				Value: book2users,
			})
			return err
		}
		return nil
	})
}

func (b *Book2UsersMySQLRepository) Delete(id int) error {
	return b.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.BooksLeaned{}, id).Error; err != nil {
			log.Log().Error("Delete book2users failed", log2.Field{
				Key:   "err",
				Value: err,
			}, log2.Field{
				Key:   "id",
				Value: id,
			})
			return err
		}
		return nil
	})
}
