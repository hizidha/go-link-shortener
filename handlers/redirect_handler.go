package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"shorted-link/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	urlCollection *mongo.Collection
)

func SetURLCollection(collection *mongo.Collection) {
	urlCollection = collection
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["shortURL"]

	var shortenedURL models.ShortenedURL
	err := urlCollection.FindOne(context.Background(), bson.M{"short_id": shortID}).Decode(&shortenedURL)
	if err != nil {
		// Prepare JSON response for not found
		response := map[string]interface{}{
			"error": "Short URL not found",
		}

		// Convert response to JSON
		jsonResponse, _ := json.Marshal(response)

		// Set Content-Type header and write JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
		return
	}

	http.Redirect(w, r, shortenedURL.LongURL, http.StatusFound)
}
