package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSaveHandler(t *testing.T) {
	// Create a Gin router for testing
	router := gin.Default()
	router.POST("/save", saveHandler)

	// Create a new server using the router
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Prepare the request payload
	payload := []byte(`{"key": "example", "value": "data"}`)

	// Send the POST request to /save
	resp, err := http.Post(ts.URL+"/save", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestGetHandler(t *testing.T) {
	// Create a Gin router for testing
	router := gin.Default()
	router.GET("/get", getHandler)

	// Create a new server using the router
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Initialize the in-memory store
	store["example"] = "data"

	// Send the GET request to /get
	resp, err := http.Get(ts.URL + "/get?key=example")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check the response body
	var res getResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if res.Value != "data" {
		t.Errorf("Expected value 'data', got '%s'", res.Value)
	}
}

func TestGetHandlerKeyNotFound(t *testing.T) {
	// Create a Gin router for testing
	router := gin.Default()
	router.GET("/get", getHandler)

	// Create a new server using the router
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Send the GET request to /get with a non-existent key
	resp, err := http.Get(ts.URL + "/get?key=nonexistent")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}
