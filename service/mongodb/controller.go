package mongodb

import (
	"context"
	"fmt"

	"github.com/meateam/fav-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/meateam/fav-service/proto"

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

func (c Controller) GetAllFavorites(ctx context.Context, userID string) ([]*pb.FavoriteObject, error) {
	filter := bson.D{
		bson.E{
			Key:   FavoriteBSONUserIDField,
			Value: userID,
		},
	}

	favoriteFiles, err := c.store.GetAll(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "favorite not found")
	}

	returnedFavFiles := make([]*pb.FavoriteObject, 0, len(favoriteFiles))
	for _, favFiles := range favoriteFiles {
		fileID := fmt.Sprintf("%v",favFiles[1].Value)
		userID := fmt.Sprintf("%v",favFiles[2].Value)
		returnedFavFiles = append(returnedFavFiles, &pb.FavoriteObject{
			UserID: userID,
			FileID: fileID,
		})
	}

	return returnedFavFiles, nil
}

func (c Controller) CreateFavorites(ctx context.Context, fileID string, userID string,) (service.Favorite, error) {
	FavoriteObject := &BSON{FileID: fileID, UserID: userID}
	createdFavorite, err := c.store.Create(ctx, FavoriteObject)
	if err != nil {
		return nil, fmt.Errorf("failed creating favorite: %v", err)
	}

	return createdFavorite, nil
}

func (c Controller) DeleteFavorites(ctx context.Context, fileID string, userID string,) (service.Favorite, error) {
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

	favorite, err := c.store.Delete(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "favorite not found")
	}

	return favorite, nil
}








