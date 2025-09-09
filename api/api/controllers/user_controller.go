package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (u *UserController) Profile(c *gin.Context) {
	userId := c.GetInt("userId")
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your profile 🎉",
		"userId":  userId,
	})
}
