package main

import (
	"log"
	"net/http"

	"chess/internal/server"
)

func main() {

	manager := server.NewGameManager()

	http.HandleFunc("/ws", manager.HandleWebSocket)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
