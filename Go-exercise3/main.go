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

	// Start server
	router.Run(":8081")
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

	// Save password to the database
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
