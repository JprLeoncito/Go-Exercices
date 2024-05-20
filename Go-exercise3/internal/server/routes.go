package server

import (
	"net/http"
	"sync"
	"time"

	"Go-exercise3/internal/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// DataStore is a simple in-memory data store using a map
type DataStore struct {
	sync.RWMutex
	data map[string]string
}

func NewDataStore() *DataStore {
	return &DataStore{
		data: make(map[string]string),
	}
}
func savePassword(c *gin.Context) {
	// Validate input
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		URL      string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a strong password hash
	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate password hash"})
		return
	}

	db := database.GetDB() // Get the database connection
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}

	_, err = db.Exec("INSERT INTO passwords (username, password, url, created_at) VALUES ($1, $2, $3, $4)",
		input.Username, hashedPassword, input.URL, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password saved successfully"})
}
func getPasswords(c *gin.Context) {
	db := database.GetDB() // Get the database connection
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}

	// Query passwords from the database
	rows, err := db.Query("SELECT id, username, password, url, created_at FROM passwords")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve passwords"})
		return
	}
	defer rows.Close()

	var passwords []gin.H // Ensure this variable is declared to store the passwords

	// Iterate over rows and build response
	for rows.Next() {
		var id int
		var username, password, url string
		var createdAt time.Time
		// Make sure to include password in the Scan method
		if err := rows.Scan(&id, &username, &password, &url, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan rows"})
			return
		}
		// Include password in the response if necessary
		passwords = append(passwords, gin.H{"id": id, "username": username, "password": password, "url": url, "created_at": createdAt})
	}

	c.JSON(http.StatusOK, passwords)
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	// Define API endpoints
	router.PUT("/passwords", savePassword)
	router.GET("/passwords", getPasswords)
	router.Run(":8080")
	return router
}
