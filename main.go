package main

import (
	"log"
	"os"

	"flow-desk/config"
	"flow-desk/handler"
	"flow-desk/repository"
	"flow-desk/routes"
	"flow-desk/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env tidak ditemukan")
	}

	db := config.ConnectDatabase()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	r := gin.Default()

	routeConfig := routes.RouteConfig{
		App:         r,
		AuthHandler: authHandler,
		TaskHandler: taskHandler,
	}
	routeConfig.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
