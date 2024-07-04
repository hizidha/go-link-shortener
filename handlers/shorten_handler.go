package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"shorted-link/models"

	"github.com/joho/godotenv"
	"github.com/speps/go-hashids"
	"go.mongodb.org/mongo-driver/bson"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

var (
	domain string
	port   string
	hashID *hashids.HashID
)

func init() {
	hashID, _ = hashids.New()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Accessing environment variables
	domain = os.Getenv("DOMAIN")
	port = os.Getenv("PORT")
	if domain == "" || port == "" {
		log.Fatal("DOMAIN or PORT not set in .env file")
	}
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest

	// Decode JSON body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	longURL := req.URL
	if longURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	// Generate a unique short code based on the longURL
	randomID := rand.Intn(100000000)
	shortID, _ := hashID.Encode([]int{randomID})

	// Check if shortID already exists
	var existingURL models.ShortenedURL
	err = urlCollection.FindOne(context.Background(), bson.M{"short_id": shortID}).Decode(&existingURL)
	if err == nil {
		// ShortID exists, check expiration
		if existingURL.ExpirationDate.After(time.Now()) {
			http.Error(w, "ShortID already exists and is still active", http.StatusConflict)
			return
		}
	}

	shortenedURL := models.ShortenedURL{
		ShortID:        shortID,
		LongURL:        longURL,
		TimestampAdded: time.Now(),
		ExpirationDate: time.Now().AddDate(0, 0, 30), // Example: 30 days expiration
	}

	// Save to MongoDB
	_, err = urlCollection.InsertOne(context.Background(), shortenedURL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	serverAddr := fmt.Sprintf("http://%s:%s", domain, port)
	shortURL := fmt.Sprintf("%s/%s", serverAddr, shortID)

	// Prepare JSON response
	response := map[string]interface{}{
		"shortened_url": shortURL,
		"status":        "success",
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
