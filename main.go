package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"shorted-link/handlers"
	"shorted-link/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	port   string // Default port
	domain string // Domain from .env
)

func main() {
	// Load environment variables from .env file
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

	// Connect to MongoDB
	client, err := utils.ConnectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Set MongoDB collection for handlers
	urlCollection := client.Database("go").Collection("shortenedLink")
	handlers.SetURLCollection(urlCollection)

	// Initialize router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", handlers.WelcomeHandler).Methods("GET")
	r.HandleFunc("/shorten", handlers.ShortenURL).Methods("POST")
	r.HandleFunc("/{shortURL}", handlers.RedirectURL).Methods("GET")

	// Start server
	serverAddr := fmt.Sprintf("%s:%s", domain, port)
	serverPrint := fmt.Sprintf("http://%s", serverAddr)
	fmt.Println("Running at", serverPrint)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
