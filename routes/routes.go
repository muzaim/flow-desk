package routes

import (
	"flow-desk/handler"
	"flow-desk/middleware"
	"flow-desk/repository"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App             *gin.Engine
	AuthHandler     *handler.AuthHandler
	TaskHandler     *handler.TaskHandler
	IdempotencyRepo repository.IdempotencyRepository
}

func (c *RouteConfig) SetupRoutes() {
	c.App.Use(middleware.GlobalErrorHandler())

	c.App.GET("/debug-panic", func(ctx *gin.Context) {
		panic("Simulasi crash fatal pada server!")
	})

	c.App.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
			"status":  "server is running",
		})
	})

	api := c.App.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", c.AuthHandler.Register)
			auth.POST("/login", c.AuthHandler.Login)
		}

		tasks := api.Group("/tasks")
		tasks.Use(middleware.AuthMiddleware())
		{
			tasks.POST("", middleware.IdempotencyMiddleware(c.IdempotencyRepo), c.TaskHandler.Create)
			tasks.GET("", c.TaskHandler.GetAll)
			tasks.GET("/:id", c.TaskHandler.GetByID)
			tasks.PUT("/:id", c.TaskHandler.Update)
			tasks.DELETE("/:id", c.TaskHandler.Delete)
		}
	}
}
