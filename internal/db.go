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

func NewDB(c *mongo.Collection) Database {
	return Database{
		Collection: c,
	}
}

func (d *Database) insert(ctx context.Context, shortUrl models.ShortURL) (interface{}, error) {
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
