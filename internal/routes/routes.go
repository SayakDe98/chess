package routes

import (
	"chess/internal/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/register", func(c *gin.Context) {
		handlers.Register(c, db)
	})
	r.POST("/login", func(c *gin.Context) {
		handlers.Login(c, db)
	})
}
