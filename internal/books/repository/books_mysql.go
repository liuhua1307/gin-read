package repository

import (
	"github.com/liuhua1307/gin-read/internal/configs"
	"github.com/liuhua1307/gin-read/internal/domain"
	"github.com/liuhua1307/gin-read/internal/pkg/log"
	log2 "github.com/liuhua1307/gin-read/pkg/log"
	"gorm.io/gorm"
)

var _ domain.BookRepository = &BookMySQLRepository{}

// BookMySQLRepository is a struct implementing the BookRepository interface
type BookMySQLRepository struct {
	d *configs.DataInstance
}

func (b *BookMySQLRepository) StartTransaction() *gorm.DB {
	return b.d.DB.Begin()
}

func NewBookMySQLRepository(d *configs.DataInstance) *BookMySQLRepository {
	return &BookMySQLRepository{d: d}
}

func (b *BookMySQLRepository) Create(book *domain.Book) error {
	return b.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(book).Error; err != nil {
			log.Log().Error("Create book failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "Book", Value: book.ISBN})
			return err
		}
		return nil
	})
}

func (b *BookMySQLRepository) FindAll() ([]domain.Book, error) {
	var books []domain.Book
	if err := b.d.DB.Find(&books).Error; err != nil {
		log.Log().Error("Find all books failed ", log2.Field{Key: "error", Value: err})
		return nil, err
	}
	return books, nil
}

func (b *BookMySQLRepository) FindByID(id int) (*domain.Book, error) {
	var book domain.Book
	if err := b.d.DB.First(&book, id).Error; err != nil {
		log.Log().Error("Find book by id failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "ID", Value: id})
		return nil, err
	}
	return &book, nil
}

func (b *BookMySQLRepository) Update(book *domain.Book) error {
	return b.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(book).Error; err != nil {
			log.Log().Error("Update book failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "Book", Value: book.ISBN})
			return err
		}
		return nil
	})
}

func (b *BookMySQLRepository) Delete(id int) error {
	// 软删除 填充 deleted_at 字段
	return b.d.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.Book{}, id).Error; err != nil {
			log.Log().Error("Delete book failed ", log2.Field{Key: "error", Value: err}, log2.Field{Key: "ID", Value: id})
			return err
		}
		return nil
	})
}
