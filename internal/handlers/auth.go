package handlers

import (
	"chess/internal/models"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

var googleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	RedirectURL: os.Getenv("API_URL") + "/auth/google/callback",
	Endpoint:    google.Endpoint,
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
	var username string
	var password_hash string
	err := h.DB.QueryRow(`SELECT username, id, password_hash FROM users WHERE email = ?`, req.Email).Scan(&username, &id, &password_hash)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User does not exist"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong credentials"})
		return
	}

	token := GenerateJWT(id, username)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func GenerateJWT(userId int, username string) string {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"username": username,
	})

	tokenString, _ := token.SignedString(secret)
	return tokenString
}

func generateRandomState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state := generateRandomState() // crypto/rand string
	// store state in a short-lived cookie to verify later
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// verify state
	cookie, _ := c.Cookie("oauth_state")
	if c.Query("state") != cookie {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid state"})
		return
	}

	// exchange code for token
	token, err := googleOauthConfig.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "exchange failed"})
		return
	}

	// fetch user info from Google
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info"})
		return
	}

	var user models.User
	err = h.DB.QueryRow(
		"SELECT id, username, email FROM users WHERE google_id = ?",
		googleUser.ID,
	).Scan(&user.Id, &user.Username, &user.Email)

	if err == sql.ErrNoRows {
		// check email conflict
		var existingID int
		emailErr := h.DB.QueryRow(
			"SELECT id FROM users WHERE email = ?",
			googleUser.Email,
		).Scan(&existingID)

		if emailErr == nil {
			// email exists — link google_id to existing account
			_, err = h.DB.Exec(
				"UPDATE users SET google_id = ? WHERE email = ?",
				googleUser.ID, googleUser.Email,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link google account"})
				return
			}
			user.Id = existingID
			user.Email = googleUser.Email
		} else {
			// brand new user — handle username conflict
			username := googleUser.Name
			var count int
			h.DB.QueryRow(
				"SELECT COUNT(*) FROM users WHERE username = ?",
				username,
			).Scan(&count)
			if count > 0 {
				username = fmt.Sprintf("%s_%s", username, generateRandomState()[:6])
			}

			result, err := h.DB.Exec(
				"INSERT INTO users (username, email, google_id) VALUES (?, ?, ?)",
				username, googleUser.Email, googleUser.ID,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
			newID, _ := result.LastInsertId()
			user.Id = int(newID)
			user.Username = username
			user.Email = googleUser.Email
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// sign JWT and redirect for ALL cases (new + existing)
	jwtToken := GenerateJWT(user.Id, user.Username)
	clientURL := os.Getenv("FRONTEND_URL")
	c.Redirect(http.StatusTemporaryRedirect, clientURL+"/auth/callback?token="+jwtToken)
}
