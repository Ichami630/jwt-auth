package api

import (
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/api/controllers"
	"github.com/Iknite-Space/sqlc-example-api/api/middlewares"
	"github.com/Iknite-Space/sqlc-example-api/db/store"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authController *controllers.AuthController
	userController *controllers.UserController
}

func NewRouter(s store.Store) http.Handler {
	r := gin.Default()

	//controllers
	userCtrl := controllers.NewUserController(s)
	authCtrl := controllers.NewAuthController(s)

	//public routes
	r.POST("/login", authCtrl.Login)
	r.POST("/register", userCtrl.Register)

	//protected routes
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/profile", userCtrl.Profile)

	return r
}
