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

	rdb := config.ConnectRedis()
	idempotencyRepo := repository.NewIdempotencyRepository(rdb)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, userRepo) // Pass userRepo di sini
	taskHandler := handler.NewTaskHandler(taskService)

	r := gin.New()
	r.Use(gin.Logger())

	routeConfig := routes.RouteConfig{
		App:             r,
		AuthHandler:     authHandler,
		TaskHandler:     taskHandler,
		IdempotencyRepo: idempotencyRepo,
	}
	routeConfig.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
