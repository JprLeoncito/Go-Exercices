package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

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

func GeneratePasswordHandler(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/generate-password", GeneratePasswordHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
