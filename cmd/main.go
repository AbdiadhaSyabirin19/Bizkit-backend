package main

import (
	"log"
	"os"

	"bizkit-backend/config"
	"bizkit-backend/internal/router"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.ConnectDB()

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	log.Println("Server berjalan di port", port)
	r.Run(":" + port)
}