package mongodb

import (
	"context"
	"fmt"

	"github.com/meateam/fav-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Controller is the favorite service business logic implementation using MongoStore.
type Controller struct {
	store MongoStore
}

func NewMongoController(db *mongo.Database) (Controller, error) {
	store, err := newMongoStore(db)
	if err != nil {
		return Controller{}, err
	}

	return Controller{store: store}, nil
}

func (c Controller) DeleteFavorite(ctx context.Context, fileID string, userID string,) (service.Favorite, error) {
	filter := bson.D{
		bson.E{
			Key:   FavoriteBSONFileIDField,
			Value: fileID,
		},
		bson.E{
			Key:   FavoriteBSONUserIDField,
			Value: userID,
		},
	}

	favorite, err := c.store.DeleteFavorite(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "favorite not found")
	}

	return favorite, nil
}

func (c Controller) CreateFavorite(ctx context.Context, fileID string, userID string,) (service.Favorite, error) {
	FavoriteObject := &BSON{FileID: fileID, UserID: userID}
	createdFavorite, err := c.store.CreateFavorite(ctx, FavoriteObject)
	if err != nil {
		return nil, fmt.Errorf("failed creating favorite: %v", err)
	}

	return createdFavorite, nil
}

func (c Controller) GetAll(ctx context.Context, userID string) ([]service.Favorite, error) {
	filter := bson.D{
		bson.E{
			Key:   FavoriteBSONUserIDField,
			Value: userID,
		},
	}

	favorite, err := c.store.GetAll(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "favorite not found")
	}

	return favorite, nil
}



