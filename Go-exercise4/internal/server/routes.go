package server

import (
	"Go-exercise4/internal/database"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RetrievePasswordsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve passwords from database
	userID := 1 // Example user ID
	passwords, err := database.RetrievePasswords(userID)
	if err != nil {
		http.Error(w, "Error retrieving passwords", http.StatusInternalServerError)
		return
	}

	// Return passwords as JSON
	json.NewEncoder(w).Encode(passwords)
}
func (s *Server) RegisterRoutes() http.Handler {
	// Initialize the database connection
	r := mux.NewRouter()
	r.HandleFunc("/generate-password", database.GeneratePasswordHandler).Methods("POST")
	r.HandleFunc("/retrieve-passwords", RetrievePasswordsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
	return r
}
