package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/liuhua1307/gin-read/internal/domain"
	"github.com/liuhua1307/gin-read/internal/pkg/response"
	"strconv"
)

type UsersHandler struct {
	uc domain.UserUseCase
}

func NewUsersHandler(uc domain.UserUseCase, router *gin.RouterGroup, mw ...gin.HandlerFunc) {
	h := &UsersHandler{
		uc: uc,
	}
	router.Use(mw...)
	router.POST("/users", h.CreateUser)
	router.GET("/users", h.GetAllUsers)
	router.GET("/users/:id", h.GetUserByID)
	router.PUT("/users/:id", h.UpdateUser)
	router.DELETE("/users/:id", h.DeleteUser)
}

func (h *UsersHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, response.FormError)
	}
	if err := h.uc.CreateUser(&user); err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithoutData(c)
}

func (h *UsersHandler) GetAllUsers(c *gin.Context) {
	users, err := h.uc.GetAllUsers()
	if err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithData(c, users)
}

func (h *UsersHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	user, err := h.uc.GetUserByID(id)
	if err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithData(c, user)
}

func (h *UsersHandler) UpdateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, response.FormError)
	}
	if err := h.uc.UpdateUser(&user); err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithoutData(c)
}

func (h *UsersHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err = h.uc.DeleteUser(id); err != nil {
		response.Error(c, response.ServerError)
		return
	}
	response.SuccessWithoutData(c)
}
