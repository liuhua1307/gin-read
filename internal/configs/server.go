package configs

import (
	"github.com/gin-gonic/gin"
	book2usersdelivery "github.com/liuhua1307/gin-read/internal/book2users/delivery"
	book2usersrepo "github.com/liuhua1307/gin-read/internal/book2users/repository"
	book2usersuc "github.com/liuhua1307/gin-read/internal/book2users/usecase"
	bookdelivery "github.com/liuhua1307/gin-read/internal/books/delivery"
	bookrepo "github.com/liuhua1307/gin-read/internal/books/repository"
	bookuc "github.com/liuhua1307/gin-read/internal/books/usecase"
	userdelivery "github.com/liuhua1307/gin-read/internal/users/delivery"
	userepo "github.com/liuhua1307/gin-read/internal/users/repository"
	useruc "github.com/liuhua1307/gin-read/internal/users/usecase"
	"net/http"
)

func ServerInit() *http.Server {

	e := gin.Default()
	db := GetDB()
	dataInstance := NewDataInstance(db)

	//Init book service
	bookRepo := bookrepo.NewBookMySQLRepository(dataInstance)
	bookUseCase := bookuc.NewBookUseCase(bookRepo)
	bookdelivery.NewBooksHandler(bookUseCase, e.Group("/api/v1/books"))

	//Init user service
	userRepo := userepo.NewUserMySQLRepository(dataInstance)
	userUseCase := useruc.NewUserUseCase(userRepo)
	userdelivery.NewUsersHandler(userUseCase, e.Group("/api/v1/users"))

	//Init Leaned service
	book2usersRepo := book2usersrepo.NewBook2UsersMySQLRepository(dataInstance)
	book2usersUseCase := book2usersuc.NewBook2UsersUseCase(book2usersRepo, bookRepo, userRepo, NewCache(rs))
	book2usersdelivery.NewBook2UsersHandler(book2usersUseCase, e.Group("/api/v1/book2users"))

	return &http.Server{
		Addr:    ":8080",
		Handler: e,
	}
}
