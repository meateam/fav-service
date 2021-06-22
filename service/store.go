package service

import (
	"context"
)

// Store is an interface for handling the storing of favorites.
// ----> ./mongodb/store.go
type Store interface {
	GetAllFavorites(ctx context.Context, filter interface{}) ([]Favorite, error)
	CreateFavorites(ctx context.Context, favorite Favorite) (Favorite, error)
	DeleteFavorites(ctx context.Context, filter interface{}) (Favorite, error)
}