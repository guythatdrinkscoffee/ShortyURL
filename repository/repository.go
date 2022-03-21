package repository

import (
	"context"
	"github.com/guythatdrinkscoffee/ShortyURL/internal"
	"github.com/guythatdrinkscoffee/ShortyURL/models"
)

type Repository struct {
	DB *internal.Database
}

func NewRepository(db *internal.Database) Repository {
	return Repository{db}
}

func (r *Repository) CreateURL(ctx context.Context, shortUrl *models.ShortURL) (interface{}, error) {
	return r.DB.Create(ctx, *shortUrl)
}

func (r *Repository) Find(ctx context.Context, hash string) (*models.ShortURL, error) {
	return r.DB.Find(ctx, hash)
}
