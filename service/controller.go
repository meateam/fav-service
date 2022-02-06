package service

import (
	"context"

	pb "github.com/meateam/fav-service/proto"
)

// Controller is an interface for the business logic of the fav.Service which uses a Store.
type Controller interface {
	CreateFavorite(ctx context.Context, fileID string, userID string) (Favorite, error)
	DeleteFavorite(ctx context.Context, fileID string, userID string) (Favorite, error)
	GetManyFavoritesByUserID(ctx context.Context, userID string) (*pb.GetManyFavoritesResponse, error)
	GetByFileAndUser(ctx context.Context, fileID string, userID string) (Favorite ,error)
	DeleteAllfileFav(ctx context.Context, fileID string) (*pb.DeleteAllfileFavResponse, error)
	HealthCheck(ctx context.Context) (bool, error)
	
}
