package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/liuhua1307/gin-read/internal/domain"
	"github.com/liuhua1307/gin-read/internal/pkg/response"
	"strconv"
)

type BooksHandler struct {
	uc domain.BookUseCase
}

func NewBooksHandler(uc domain.BookUseCase, router *gin.RouterGroup, mw ...gin.HandlerFunc) {
	h := &BooksHandler{
		uc: uc,
	}
	router.Use(mw...)
	router.POST("", h.CreateBook)
	router.GET("", h.GetAllBooks)
	router.GET("/:id", h.GetBookByID)
	router.PUT("/:id", h.UpdateBook)
	router.DELETE("/:id", h.DeleteBook)
}

func (h *BooksHandler) CreateBook(c *gin.Context) {
	var book domain.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		response.Error(c, response.FormError)
	}
	if err := h.uc.CreateBook(&book); err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithoutData(c)
}

func (h *BooksHandler) GetAllBooks(c *gin.Context) {
	books, err := h.uc.GetAllBooks()
	if err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithData(c, books)
}

func (h *BooksHandler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	book, err := h.uc.GetBookByID(id)
	if err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithData(c, book)
}

func (h *BooksHandler) UpdateBook(c *gin.Context) {
	var book domain.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		response.Error(c, response.FormError)
	}
	if err := h.uc.UpdateBook(&book); err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithoutData(c)
}

func (h *BooksHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err = h.uc.DeleteBook(id); err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithoutData(c)
}
