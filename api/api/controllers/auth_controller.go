package controllers

import (
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/db/store"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("supersecrets")

type AuthController struct {
	store store.Store
}

func NewAuthController(store store.Store) *AuthController {
	return &AuthController{
		store: store,
	}
}

func (h *AuthController) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	//get user from db
	user, err := h.store.Do().GetUserByEmail(c, body.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email/password"})
		return
	}

	//compare plain password to stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
}
