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

var _ domain.BookUseCase = &BookUseCase{}

type BookUseCase struct {
	bookRepo domain.BookRepository
	Cache    configs.Cache
}

func NewBookUseCase(bookRepo domain.BookRepository) *BookUseCase {
	return &BookUseCase{bookRepo: bookRepo}
}

func (b BookUseCase) CreateBook(book *domain.Book) error {
	tx := b.bookRepo.StartTransaction()
	err := tx.Model(&domain.Book{}).Create(book).Error
	if err != nil {
		log.Log().Error("err create book:", log2.Field{
			Key:   "err",
			Value: err,
		})
		tx.Rollback()
		return err
	}
	//初始化库存
	//其实这里还要锁
	info, s := getLockInfo(int(book.ID))
	cnt := 0
	for lock, err := b.Cache.TryLock(info, s, time.Minute); !lock; {
		if err != nil {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: err,
			}, log2.Field{
				Key: "bookId", Value: book.ID,
			})
		}
		cnt++
		if cnt > 5 {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: "Get lock failed",
			}, log2.Field{
				Key: "bookId", Value: book.ID,
			})
			tx.Rollback()
			return err
		}
	}
	err = b.Cache.Set(fmt.Sprintf("stock:%d", book.ID), strconv.Itoa(book.Count))
	if err != nil {
		log.Log().Error("Set stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book.ID,
		})
		tx.Rollback()
		return err
	}
	tx.Commit()
	err = b.Cache.Unlock(info)
	if err != nil {
		log.Log().Error("Unlock failed", log2.Field{
			Key: "err", Value: err,
		})
	}
	return nil
}

func (b BookUseCase) GetAllBooks() ([]domain.Book, error) {
	return b.bookRepo.FindAll()
}

func (b BookUseCase) GetBookByID(id int) (*domain.Book, error) {
	return b.bookRepo.FindByID(id)
}

func (b BookUseCase) UpdateBook(book *domain.Book) error {
	tx := b.bookRepo.StartTransaction()
	err := tx.Model(&domain.Book{}).Where("id = ?", book.ID).Update("count", book.Count).Error
	if err != nil {
		log.Log().Error("err update book:", log2.Field{
			Key:   "err",
			Value: err,
		})
		tx.Rollback()
		return err
	}
	//获取锁 更新一下库存 使用CAS
	//获取锁
	info, s := getLockInfo(int(book.ID))
	cnt := 0
	for lock, err := b.Cache.TryLock(info, s, time.Minute); !lock; {
		if err != nil {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: err,
			}, log2.Field{
				Key: "bookId", Value: book.ID,
			})
		}
		cnt++
		if cnt > 5 {
			log.Log().Error("Get lock failed", log2.Field{
				Key: "err", Value: "Get lock failed",
			}, log2.Field{
				Key: "bookId", Value: book.ID,
			})
			tx.Rollback()
			return err
		}
	}

	//查询缓存查看库存是不是需要更新
	stock, err := b.Cache.Get(fmt.Sprintf("stock:%d", book.ID))
	if err != nil {
		log.Log().Error("Get stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book.ID,
		})
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
	if stockInt == book.Count {
		return nil
	}
	//更新库存
	err = b.Cache.Set(fmt.Sprintf("stock:%d", book.ID), strconv.Itoa(book.Count))
	if err != nil {
		log.Log().Error("Set stock failed", log2.Field{
			Key: "err", Value: err,
		}, log2.Field{
			Key: "bookId", Value: book,
		})
		tx.Rollback()
		return err
	}
	tx.Commit()
	err = b.Cache.Unlock(info)
	if err != nil {
		log.Log().Error("Unlock failed", log2.Field{
			Key: "err", Value: err,
		})
	}
	return nil
}

func (b BookUseCase) DeleteBook(id int) error {
	return b.bookRepo.Delete(id)
}

func getLockInfo(bookId int) (string, string) {
	key := fmt.Sprintf("lock:book:%d", bookId)
	value := fmt.Sprintf("lock:%d", time.Now().Unix())
	return key, value
}
