package domain

import "gorm.io/gorm"

// Book table books model
type Book struct {
	gorm.Model
	Name  string
	Count int
	ISBN  string
}

// BookRepository interface
type BookRepository interface {
	Create(book *Book) error
	FindAll() ([]Book, error)
	FindByID(id int) (*Book, error)
	Update(book *Book) error
	Delete(id int) error
	StartTransaction() *gorm.DB
}

// BookUseCase  interface
type BookUseCase interface {
	CreateBook(book *Book) error
	GetAllBooks() ([]Book, error)
	GetBookByID(id int) (*Book, error)
	UpdateBook(book *Book) error
	DeleteBook(id int) error
}
