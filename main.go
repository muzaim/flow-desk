package main

import (
	"log"
	"os"

	"flow-desk/config"
	"flow-desk/entity"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable bawaan.")
	}

	// Connect Database
	db := config.ConnectDatabase()

	// Auto Migration Model ke Database
	err = db.AutoMigrate(&entity.User{}, &entity.Task{})
	if err != nil {
		log.Fatalf("Gagal melakukan auto migration: %v", err)
	}

	// Inisialisasi Router Gin
	r := gin.Default()

	// Health check route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "server is running",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
