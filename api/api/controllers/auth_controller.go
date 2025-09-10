package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Iknite-Space/sqlc-example-api/db/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

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
		"exp": time.Now().Add(1 * time.Minute).Unix(),
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

// refresh token function
func (h *AuthController) Refresh(c *gin.Context) {
	//get the refresh token from request
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}

	//parse and validate the refresh token
	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		// Ensure HMAC signing
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		}
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	userId := claims["sub"]

	//issue new access token
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	})
	accessTokenString, _ := newAccessToken.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessTokenString,
	})
}
