package domain

import "gorm.io/gorm"

type BooksLeaned struct {
	gorm.Model
	BookID int
	UserID int
}

type Book2UsersRepository interface {
	Create(book2users *BooksLeaned) error
	FindAll() ([]BooksLeaned, error)
	FindByID(id int) (*BooksLeaned, error)
	Update(book2users *BooksLeaned) error
	Delete(id int) error
	FindByUserID(userId int) (*BooksLeaned, error)
	StartTransaction() *gorm.DB
}

type Book2UsersUseCase interface {
	LendBook(book2users *BooksLeaned) error
	ReturnBook(book2users *BooksLeaned) error
	GetAllBooksLeaned() ([]BooksLeaned, error)
	GetBookLeanedByID(id int) (*BooksLeaned, error)
	GetBookLeanedByUserID(userId int) (*BooksLeaned, error)
}
