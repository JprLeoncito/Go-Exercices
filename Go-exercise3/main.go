package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "jasper31"
	dbName     = "postgres"
)

var db *sql.DB

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

func initDB() {
	connectionString := "host=" + dbHost + " port=" + strconv.Itoa(dbPort) + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Ping the database to check if the connection is successful
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}
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

func RetrievePasswordsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve passwords from database
	userID := 1 // Example user ID
	passwords, err := RetrievePasswords(userID)
	if err != nil {
		http.Error(w, "Error retrieving passwords", http.StatusInternalServerError)
		return
	}

	// Return passwords as JSON
	json.NewEncoder(w).Encode(passwords)
}

func main() {
	// Initialize the database connection
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/generate-password", GeneratePasswordHandler).Methods("POST")
	r.HandleFunc("/retrieve-passwords", RetrievePasswordsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
