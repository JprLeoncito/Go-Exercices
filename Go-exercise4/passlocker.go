package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "admin"
	dbName     = "postgres"
)

var db *sql.DB

func main() {
	// Initialize database connection
	connectionString := "host=" + dbHost + " port=" + strconv.Itoa(dbPort) + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create Gin router
	router := gin.Default()

	// Define API endpoints
	router.POST("/passwords", savePassword)
	router.GET("/passwords", getPasswords)
	router.POST("/save", saveHandler) // Add /save endpoint
	router.GET("/get", getHandler)    // Add /get endpoint
	router.POST("/login", login)

	// Start server
	router.Run(":8080")

}

type saveRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type getResponse struct {
	Value string `json:"value"`
}

// In-memory store
var store = make(map[string]string)

func login(c *gin.Context) {
	// Validate input
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user from the database
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM passwords WHERE username = $1", input.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		}
		return
	}
	// Compare hashed password with provided password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// If valid, return success message (you can also generate and return a token here)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
func saveHandler(c *gin.Context) {
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	store[req.Key] = req.Value
	c.JSON(http.StatusOK, gin.H{"message": "Data saved successfully"})
}

func getHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key is required"})
		return
	}

	value, exists := store[key]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	c.JSON(http.StatusOK, getResponse{Value: value})
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

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	db := db // Get the database connection
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
		return
	}
	// Save hashed password to the database
	_, err = db.Exec("INSERT INTO passwords (username, password, url, created_at) VALUES ($1, $2, $3, $4)",
		input.Username, hashedPassword, input.URL, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password saved successfully"})
}

func getPasswords(c *gin.Context) {
	// Query passwords from the database
	rows, err := db.Query("SELECT id, username, password, url, created_at FROM passwords")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve passwords"})
		return
	}
	defer rows.Close()

	var passwords []gin.H

	// Iterate over rows and build response
	for rows.Next() {
		var id int
		var username, password, url string
		var createdAt time.Time
		if err := rows.Scan(&id, &username, &password, &url, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan rows"})
			return
		}
		passwords = append(passwords, gin.H{"id": id, "username": username, "password": password, "url": url, "created_at": createdAt})
	}

	c.JSON(http.StatusOK, passwords)
}
