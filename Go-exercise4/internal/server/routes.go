package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

// HandleSave handles the POST /save endpoint
func HandleSave(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON payload
	var payload map[string]string
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Store the key-value pair
	store.Lock()
	defer store.Unlock()
	for key, value := range payload {
		store.data[key] = value
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Data saved successfully")
}

// HandleGet handles the GET /get endpoint
func HandleGet(w http.ResponseWriter, r *http.Request) {
	// Extract the key from the query parameter
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key parameter is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the value from the data store
	store.RLock()
	defer store.RUnlock()
	value, found := store.data[key]
	if !found {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	// Write the value to the response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Value for key '%s': %s", key, value)
}

var store = NewDataStore()

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	http.HandleFunc("/save", HandleSave)
	http.HandleFunc("/get", HandleGet)

	return r
}
