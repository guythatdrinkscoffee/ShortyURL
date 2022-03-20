package models

import "time"

type ShortURL struct {
	ID             string    `bson:"_id,omitempty" json:"id,omitempty"`
	Hash           string    `bson:"hash,omitempty" json:"hash,omitempty"`
	OriginalURL    string    `bson:"original_url,omitempty" json:"original_url,omitempty"`
	DateCreated    time.Time `bson:"date_created,omitempty" json:"date_created,omitempty"`
	ExpirationDate time.Time `bson:"expiration_date,omitempty" json:"expiration_date,omitempty"`
}