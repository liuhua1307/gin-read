package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/liuhua1307/gin-read/internal/domain"
	"github.com/liuhua1307/gin-read/internal/pkg/response"
	"strconv"
)

type Book2UsersHandler struct {
	book2usersUsecase domain.Book2UsersUseCase
}

func (h Book2UsersHandler) LeanedBook2User(context *gin.Context) {
	var book2users domain.BooksLeaned
	if err := context.ShouldBindJSON(&book2users); err != nil {
		response.Error(context, response.FormError)
		return
	}
	if err := h.book2usersUsecase.LendBook(&book2users); err != nil {
		response.Error(context, response.ServerError)
		return
	}
	response.SuccessWithoutData(context)
}

func (h Book2UsersHandler) GetAllBook2Users(context *gin.Context) {
	leaned, err := h.book2usersUsecase.GetAllBooksLeaned()
	if err != nil {
		response.Error(context, response.ServerError)
		return
	}
	response.SuccessWithData(context, leaned)
}

func (h Book2UsersHandler) GetBook2UsersByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		response.Error(context, response.FormError)
		return
	}
	book2users, err := h.book2usersUsecase.GetBookLeanedByID(id)
	if err != nil {
		response.Error(context, response.ServerError)
		return
	}
	response.SuccessWithData(context, book2users)
}

func (h Book2UsersHandler) ReturnBook(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		response.Error(context, response.FormError)
		return
	}
	if err := h.book2usersUsecase.ReturnBook(&domain.BooksLeaned{BookID: id}); err != nil {
		response.Error(context, response.ServerError)
		return
	}
	response.SuccessWithoutData(context)
}

func NewBook2UsersHandler(uc domain.Book2UsersUseCase, router *gin.RouterGroup, mw ...gin.HandlerFunc) {
	h := &Book2UsersHandler{
		book2usersUsecase: uc,
	}
	router.Use(mw...)
	router.POST("", h.LeanedBook2User)
	router.GET("", h.GetAllBook2Users)
	router.GET("/:id", h.GetBook2UsersByID)
	router.DELETE("/:id", h.ReturnBook)
}
