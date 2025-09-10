package api

import (
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/api/controllers"
	"github.com/Iknite-Space/sqlc-example-api/api/middlewares"
	"github.com/Iknite-Space/sqlc-example-api/db/store"
	"github.com/gin-gonic/gin"
)

type Router struct {
	store store.Store
}

func NewRouter(store store.Store) *Router {
	return &Router{
		store: store,
	}
}

func (h *Router) WireHttpHandler() http.Handler {
	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	//controllers
	userCtrl := controllers.NewUserController(h.store)
	authCtrl := controllers.NewAuthController(h.store)

	//public routes
	r.POST("/login", authCtrl.Login)
	r.POST("/register", userCtrl.Register)
	r.POST("/refresh", authCtrl.Refresh)

	//protected routes
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/profile", userCtrl.Profile)

	return r
}
