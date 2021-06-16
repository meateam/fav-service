package service

import (
	"context"
)


// Controller is an interface for the business logic of the fav.Service which uses a Store.
type Controller interface {
	CreateFavorite(ctx context.Context, fileID string, userID string) (Favorite, error)
	DeleteFavorite(ctx context.Context, fileID string, userID string) (Favorite, error)
	GetAll(ctx context.Context, userID string) ([]string, error)
}
