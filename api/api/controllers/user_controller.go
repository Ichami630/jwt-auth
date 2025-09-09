package controllers

import (
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/Iknite-Space/sqlc-example-api/db/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	store store.Store
}

func NewUserController(store store.Store) *UserController {
	return &UserController{
		store: store,
	}
}

func (u *UserController) Profile(c *gin.Context) {
	userId := c.GetInt("userId")
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to your profile 🎉",
		"userId":  userId,
	})
}

func (u *UserController) Register(c *gin.Context) {
	var req repo.CreateUserParams

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//hash the password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	//create new user
	user, err := u.store.Do().CreateUser(c, repo.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPwd),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    user,
	})

}
