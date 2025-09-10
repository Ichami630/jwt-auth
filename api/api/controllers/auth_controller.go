package controllers

import (
	"net/http"
	"time"

	"github.com/Iknite-Space/sqlc-example-api/db/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	//create access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	accessToken, _ := token.SignedString(jwtSecret)

	//create refresh token
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	refreshToken, _ := refresh.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
