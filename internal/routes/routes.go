package routes

import (
	"chess/internal/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	authHandler := &handlers.AuthHandler{DB: db}
	r.POST("/users/register", authHandler.Register)
	r.POST("/users/login", authHandler.Login)

	r.GET("/auth/google", authHandler.GoogleLogin)
	r.GET("/auth/google/callback", authHandler.GoogleCallback)
}
