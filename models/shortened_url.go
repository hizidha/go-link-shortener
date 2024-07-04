package models

import (
	"time"
)

type ShortenedURL struct {
	ShortID        string    `json:"short_id" bson:"short_id"`
	LongURL        string    `json:"long_url" bson:"long_url"`
	TimestampAdded time.Time `json:"timestamp_added" bson:"timestamp_added"`
	ExpirationDate time.Time `json:"expiration_date" bson:"expiration_date"`
}
