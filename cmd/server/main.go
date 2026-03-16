package main

import (
	"chess/internal/db"
	"chess/internal/migration"
	"chess/internal/routes"
	"chess/internal/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	manager := server.NewGameManager()

	http.HandleFunc("/ws", manager.HandleWebSocket)
	db := db.InitDB()

	r := gin.Default()
	routes.RegisterRoutes(r, db)
	migration.RunMigrations(db)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
