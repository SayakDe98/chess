package main

import (
	"chess/internal/db"
	"chess/internal/migration"
	"chess/internal/routes"
	"chess/internal/server"
	"fmt"

	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	manager := server.NewGameManager()
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load("../../.env")
		if err != nil {
			fmt.Println("No .env file present")
		}
	}

	db := db.InitDB()

	r := gin.Default()

	origin := os.Getenv("FRONTEND_URL")
	if origin == "" {
		log.Fatal("FRONTEND_URL is required")
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/ws", func(c *gin.Context) {
		manager.HandleWebSocket(c.Writer, c.Request)
	})
	routes.RegisterRoutes(r, db)
	migration.RunMigrations(db)
	log.Println("Server running on :8080")
	log.Fatal(r.Run(":8080"))
}
