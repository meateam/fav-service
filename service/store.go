package service

import (
	"context"
)

// Store is an interface for handling the storing of favorites.
type Store interface {
	GetAll(ctx context.Context, filter interface{}) ([]Favorite, error)
	Create(ctx context.Context, favorite Favorite) (Favorite, error)
	Delete(ctx context.Context, filter interface{}) (Favorite, error)
	Get(ctx context.Context, filter interface{}) (Favorite, error)
	HealthCheck(ctx context.Context) (bool, error)

}