package internal

import (
	"context"
	"github.com/guythatdrinkscoffee/ShortyURL/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	Collection *mongo.Collection
}

func NewDB(c *mongo.Collection) *Database {
	return &Database{c}
}

func (d *Database) Create(ctx context.Context, shortUrl models.ShortURL) (interface{}, error) {
	id, err := d.Collection.InsertOne(ctx, bson.D{
		{"hash", shortUrl.Hash},
		{"original_url", shortUrl.OriginalURL},
		{"created_at", shortUrl.DateCreated},
		{"expires_at", shortUrl.ExpirationDate},
	})

	if err != nil {
		return nil, err
	}

	return id.InsertedID, nil
}

func (d *Database) Find(ctx context.Context, hash string) (*models.ShortURL, error) {
	var doc *models.ShortURL

	err := d.Collection.FindOne(ctx, bson.M{"hash": hash}).Decode(&doc)

	if err != nil {
		return nil, err
	}

	return doc, err
}
