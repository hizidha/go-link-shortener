package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
)

var (
	hashID *hashids.HashID
	urlMap map[string]string
	port   string = ":8080"
	host   string = "http://localhost" + port + "/"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
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
	// Generate a random number for unique ID
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(100000000) // Generates a random number between 0 and 99999999

	shortID, _ := hashID.Encode([]int{randomID})

	shortURL := fmt.Sprintf("%s%s", host, shortID)
	urlMap[shortID] = longURL

	fmt.Fprintf(w, "Shortened URL: %s", shortURL)
}

func redirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["shortURL"]

	longURL, found := urlMap[shortID]
	if !found {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to GoLang!")
}

func main() {
	urlMap = make(map[string]string)
	hashID, _ = hashids.New()

	r := mux.NewRouter()

	r.HandleFunc("/", welcomeHandler).Methods("GET")
	r.HandleFunc("/shorten", shortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", redirectURL).Methods("GET")

	fmt.Println("Running at", host)
	log.Fatal(http.ListenAndServe(port, r))
}
