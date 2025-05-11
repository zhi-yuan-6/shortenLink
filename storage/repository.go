package storage

import (
	"context"
	"shortenLink/models"
)

type Repository interface {
	CreateShortenUrl(ctx context.Context, su *models.ShortUrl) error
	GetShortenUrl(ctx context.Context, code string) (string, error)
	IncrementVisit(ctx context.Context, code string) error
	DeleteExpired(ctx context.Context) (int64, error)
}
