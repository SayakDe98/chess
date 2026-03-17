package main

import (
	"chess/internal/db"
	"chess/internal/migration"
	"chess/internal/routes"
	"chess/internal/server"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	manager := server.NewGameManager()

	db := db.InitDB()

	r := gin.Default()
	r.GET("/ws", func(c *gin.Context) {
		manager.HandleWebSocket(c.Writer, c.Request)
	})
	routes.RegisterRoutes(r, db)
	migration.RunMigrations(db)
	log.Println("Server running on :8080")
	log.Fatal(r.Run(":8080"))
}
