package configs

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
	err  error
)

func GetDB() *gorm.DB {
	once.Do(func() {
		db, err = gorm.Open(mysql.Open(os.Getenv("DB_DSN")), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err != nil {
			panic(err)
		}
	})
	return db
}
