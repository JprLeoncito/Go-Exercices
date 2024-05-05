package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
}

type service struct {
	db *sql.DB
}

var db *sql.DB
var (
	host     = "localhost"
	port     = "5432"
	database = "posgres"
	username = "jasper31"
	password = "posgres"

	dbInstance *service
)
var (
//database = os.Getenv("DB_DATABASE")
//	password = os.Getenv("DB_PASSWORD")
//username = os.Getenv("DB_USERNAME")
//port     = os.Getenv("DB_PORT")
//host     = os.Getenv("DB_HOST")

// dbInstance *service
)

type Password struct {
	ID           int       `json:"id"`
	Password     string    `json:"password"`
	CreationDate time.Time `json:"creation_date"`
	UserID       int       `json:"user_id"`
}

const (
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers          = "0123456789"
	symbols          = "!@#$%^&*()-_=+,.?/:;{}[]~"
)

type PasswordRequest struct {
	Length           int    `json:"length"`
	IncludeNumbers   bool   `json:"includeNumbers"`
	IncludeSymbols   bool   `json:"includeSymbols"`
	IncludeUppercase bool   `json:"includeUppercase"`
	Type             string `json:"type"`
}

type PasswordResponse struct {
	Password string `json:"password"`
}

func generateRandomPassword(length int, includeNumbers, includeSymbols, includeUppercase bool) string {
	var chars string
	if includeNumbers {
		chars += numbers
	}
	if includeSymbols {
		chars += symbols
	}
	if includeUppercase {
		chars += uppercaseLetters
	}
	chars += lowercaseLetters

	var password strings.Builder
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		password.WriteByte(chars[randomIndex.Int64()])
	}
	return password.String()
}
func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	//connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	connStr := "host=" + host + " port=" + (port) + " user=" + username + " password=" + password + " dbname=" + database + " sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func GeneratePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Generate password...
	var req PasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var passwordType string
	switch req.Type {
	case "random":
		passwordType = "random"
	case "alphanumeric":
		passwordType = "alphanumeric"
	case "pin":
		passwordType = "pin"
	default:
		http.Error(w, "Invalid password type. Please choose 'random', 'alphanumeric', or 'pin'.", http.StatusBadRequest)
		return
	}

	var password string
	switch passwordType {
	case "random":
		password = generateRandomPassword(req.Length, req.IncludeNumbers, req.IncludeSymbols, req.IncludeUppercase)
	case "alphanumeric":
		password = generateRandomPassword(req.Length, true, false, req.IncludeUppercase)
	case "pin":
		password = generateRandomPassword(6, true, false, false)
	}

	res := PasswordResponse{Password: password}
	json.NewEncoder(w).Encode(res)

	// Save password to database
	userID := 1 // Example user ID
	err = SavePassword(password, userID)
	if err != nil {
		http.Error(w, "Error saving password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func SavePassword(password string, userID int) error {
	_, err := db.Exec("INSERT INTO passwords (password_string, user_id) VALUES ($1, $2)", password, userID)
	if err != nil {
		return err
	}
	return nil
}

func RetrievePasswords(userID int) ([]Password, error) {
	rows, err := db.Query("SELECT id, password_string, creation_date FROM passwords WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []Password
	for rows.Next() {
		var p Password
		err := rows.Scan(&p.ID, &p.Password, &p.CreationDate)
		if err != nil {
			return nil, err
		}
		p.UserID = userID
		passwords = append(passwords, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return passwords, nil
}
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
