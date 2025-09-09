package api

import (
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/api/controllers"
	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userController *controllers.UserController
}

func NewRouter(q repo.Querier) http.Handler {
	r := gin.Default()

	//controllers
	userCtrl := controllers.NewUserController()

	//public routes

	//protected routes
	protected := r.Group("/")
	protected.GET("/profile", userCtrl.Profile)

	return r
}
