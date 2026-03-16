package main

import (
	"log"
	"os"

	"bizkit-backend/config"
	"bizkit-backend/internal/router"

	"github.com/joho/godotenv"
)

// @title           Bizkit Backend API
// @version         1.0
// @description     API documentation for the Bizkit backend application.
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	godotenv.Load()

	config.ConnectDB()

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	log.Println("Server berjalan di port", port)
	r.Run(":" + port)
}