package usecase

import (
	"fmt"
	"github.com/liuhua1307/gin-read/internal/configs"
	"github.com/liuhua1307/gin-read/internal/domain"
	"github.com/liuhua1307/gin-read/internal/pkg/log"
	log2 "github.com/liuhua1307/gin-read/pkg/log"
	"strconv"
	"time"
)

var _ domain.Book2UsersUseCase = &Book2UsersUseCase{}

type Book2UsersUseCase struct {
	book2usersRepo domain.Book2UsersRepository
	bookRepo       domain.BookRepository
	userRepo       domain.UserRepository
	Cache          configs.Cache
}

func NewBook2UsersUseCase(book2usersRepo domain.Book2UsersRepository, bookRepo domain.BookRepository, userRepo domain.UserRepository, cache configs.Cache) *Book2UsersUseCase {
	return &Book2UsersUseCase{book2usersRepo: book2usersRepo, bookRepo: bookRepo, userRepo: userRepo, Cache: cache}
}

func (b Book2UsersUseCase) LendBook(book2users *domain.BooksLeaned) error {
	// get book by id
	_, err := b.bookRepo.FindByID(book2users.BookID)
	if err != nil {
		log.Log().Error("Get book by id failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book2users.BookID,
		})
		return err
	}
	// 获取分布式锁
	info, s := getLockInfo(book2users.BookID)
	cnt := 0
	for lock, err := b.Cache.TryLock(info, s, time.Minute); !lock; {
		if err != nil {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: err,
			}, log2.Field{
				Key: "bookId", Value: book2users.BookID,
			})
		}
		cnt++
		if cnt > 5 {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: "Get lock failed",
			}, log2.Field{
				Key: "bookId", Value: book2users.BookID,
			})
			return err
		}
	}
	// 获取到锁后，查询缓存判断库存是否足够
	stock, err := b.Cache.Get(fmt.Sprintf("stock:%d", book2users.BookID))
	if err != nil {
		log.Log().Error("Get stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book2users.BookID,
		})
		return err
	}
	stockInt, err := strconv.Atoi(stock)
	if err != nil {
		log.Log().Error("Convert stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "stock", Value: stock,
		})
		return err
	}
	if stockInt <= 0 {
		log.Log().Error("Stock is not enough", log2.Field{
			Key: "err", Value: "Stock is not enough",
		}, log2.Field{
			Key: "stock", Value: stock,
		})
		return fmt.Errorf("stock is not enough")
	}
	// 库存足够，减少库存
	stockInt--
	err = b.Cache.Set(fmt.Sprintf("stock:%d", book2users.BookID), strconv.Itoa(stockInt))
	if err != nil {
		log.Log().Error("Set stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "stock", Value: stockInt,
		})
		return err
	}
	// 创建借书记录
	err = b.book2usersRepo.Create(book2users)
	if err != nil {
		log.Log().Error("Create book2users failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "book2users", Value: book2users,
		})
		return err
	}
	// 释放锁
	err = b.Cache.Unlock(info)
	if err != nil {
		log.Log().Error("Unlock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book2users.BookID,
		})
	}
	return nil
}

func (b Book2UsersUseCase) ReturnBook(book2users *domain.BooksLeaned) error {
	// get book by id
	_, err := b.bookRepo.FindByID(book2users.BookID)
	if err != nil {
		log.Log().Error("Get book by id failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book2users.BookID,
		})
		return err
	}
	tx := b.book2usersRepo.StartTransaction()
	// 查询借书记录
	book2users, err = b.book2usersRepo.FindByID(int(book2users.ID))
	if err != nil {
		log.Log().Error("Find book2users by id failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "id", Value: book2users.ID,
		})
		return err
	}
	// 删除借书记录
	err = tx.Model(&domain.BooksLeaned{}).Where("id=?", book2users.ID).Delete(book2users).Error
	if err != nil {
		log.Log().Error("Delete book2users failed", log2.Field{
			Key: "err", Value: err,
		})
		tx.Rollback()
		return err
	}
	// 获取分布式锁
	info, s := getLockInfo(book2users.BookID)
	cnt := 0
	for lock, err := b.Cache.TryLock(info, s, time.Minute); !lock; {
		if err != nil {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: err,
			}, log2.Field{
				Key: "bookId", Value: book2users.BookID,
			})
		}
		cnt++
		if cnt > 5 {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: "Get lock failed",
			}, log2.Field{
				Key: "bookId", Value: book2users.BookID,
			})
			tx.Rollback()
			return err
		}
	}
	// 获取到锁后，查询缓存增加库存
	stock, err := b.Cache.Get(fmt.Sprintf("stock:%d", book2users.BookID))
	if err != nil {
		log.Log().Error("Get stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book2users.BookID,
		})
		tx.Rollback()
		return err
	}
	stockInt, err := strconv.Atoi(stock)
	if err != nil {
		log.Log().Error("Convert stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "stock", Value: stock,
		})
		tx.Rollback()
		return err
	}
	stockInt++
	err = b.Cache.Set(fmt.Sprintf("stock:%d", book2users.BookID), strconv.Itoa(stockInt))
	if err != nil {
		log.Log().Error("Set stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "stock", Value: stockInt,
		})
		tx.Rollback()
		return err
	}
	tx.Commit()
	// 释放锁
	err = b.Cache.Unlock(info)
	if err != nil {
		log.Log().Error("Unlock failed", log2.Field{
			Key: "err", Value: err,
		})
	}

	return nil
}

func (b Book2UsersUseCase) GetAllBooksLeaned() ([]domain.BooksLeaned, error) {
	var book2users []domain.BooksLeaned
	book2users, err := b.book2usersRepo.FindAll()
	if err != nil {
		log.Log().Error("Find all book2users failed", log2.Field{
			Key: "err", Value: err,
		})
		return nil, err
	}
	return book2users, nil
}

func (b Book2UsersUseCase) GetBookLeanedByID(id int) (*domain.BooksLeaned, error) {
	book2users, err := b.book2usersRepo.FindByID(id)
	if err != nil {
		log.Log().Error("Find book2users by id failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "id", Value: id,
		})
		return nil, err
	}
	return book2users, nil
}

func (b Book2UsersUseCase) GetBookLeanedByUserID(userId int) (*domain.BooksLeaned, error) {
	book2users, err := b.book2usersRepo.FindByUserID(userId)
	if err != nil {
		log.Log().Error("Find book2users by user id failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "userId", Value: userId,
		})
		return nil, err
	}
	return book2users, nil
}

func getLockInfo(bookId int) (string, string) {
	key := fmt.Sprintf("lock:book:%d", bookId)
	value := fmt.Sprintf("lock:%d", time.Now().Unix())
	return key, value
}
