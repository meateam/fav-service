package service

import (
	"context"
	pb "github.com/meateam/fav-service/proto"

)


// Controller is an interface for the business logic of the fav.Service which uses a Store.
type Controller interface {
	CreateFavoriteByUserAndFile(ctx context.Context, fileID string, userID string) (Favorite, error)
	DeleteFavoriteByUserAndFile(ctx context.Context, fileID string, userID string) (Favorite, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*pb.FavoriteObject, error)
}
