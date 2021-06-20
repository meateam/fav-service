package service

import (
	"context"
)

// Store is an interface for handling the storing of favorites.
// ----> ./mongodb/store.go
type Store interface {
	GetAll(ctx context.Context, filter interface{}) ([]Favorite, error)
	CreateFavorite(ctx context.Context, favorite Favorite) (Favorite, error)
	DeleteFavorite(ctx context.Context, filter interface{}) (Favorite, error)
}