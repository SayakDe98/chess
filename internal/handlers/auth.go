package handlers

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *sql.DB
}

type RegisterRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	res, err := h.DB.Exec(
		`INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`, req.UserName, req.Email, hash)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User exists"})
		return
	}
	id, _ := res.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": gin.H{
		"id":       id,
		"username": req.UserName,
		"email":    req.Email,
	}})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	var id int
	var password_hash string
	err := h.DB.QueryRow(`SELECT id, password_hash FROM users WHERE email = ?`, req.Email).Scan(&id, &password_hash)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong credentials"})
		return
	}

	token := GenerateJWT(id)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func GenerateJWT(userId int) string {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})

	tokenString, _ := token.SignedString(secret)
	return tokenString
}
