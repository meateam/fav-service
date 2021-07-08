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

// NewMongoController returns a new controller.
func NewMongoController(db *mongo.Database) (Controller, error) {
	store, err := newMongoStore(db)
	if err != nil {
		return Controller{}, err
	}

	return Controller{store: store}, nil

}


// GetAllFavoritesByUserID gets all user favorite files by userID
func (c Controller) GetAllFavoritesByUserID(ctx context.Context, userID string) (*pb.GetManyFavoritesResponse, error) {
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

	var returnedFiles []*pb.FileIDObject

	for _, file := range favoriteFiles {
		returnedFiles = append(returnedFiles, &pb.FileIDObject{FileID: file.FileID})
	}

	return &pb.GetManyFavoritesResponse{FavFileIDList: returnedFiles}, nil

}


// CreateFavorite creates a Favorite in store and returns the created favorite.
func (c Controller) CreateFavorite(ctx context.Context, fileID string, userID string,) (service.Favorite, error) {
	FavoriteObject := &BSON{FileID: fileID, UserID: userID}
	createdFavorite, err := c.store.Create(ctx, FavoriteObject)
	if err != nil {
		return nil, fmt.Errorf("failed creating favorite: %v", err)
	}

	return createdFavorite, nil

}


// DeleteFavorite deletes the favorite in store that matches userID and fileID. 
// returns the deleted favorite / error 
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

	favorite, err := c.store.Delete(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "favorite not found")
	}

	return favorite, nil

}

// HealthCheck runs store's healthcheck and returns true if healthy, otherwise returns false
// and any error if occurred.
func (c Controller) HealthCheck(ctx context.Context) (bool, error) {
	return c.store.HealthCheck(ctx)
	
}


// GetByFileAndUser retrieves the favorite that matches fileID and userID, and any error if occurred.
func (c Controller) GetByFileAndUser(ctx context.Context, fileID string, userID string) (service.Favorite, error) {
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

	favorite, err := c.store.Get(ctx, filter)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	
	if err == mongo.ErrNoDocuments {
		return nil, status.Error(codes.NotFound, "favorite not found")
	}

	return favorite, nil

}

// DeleteAllfileFav deletes all favorites of a file and returns true and the number of favorites that were deleted.
func (c Controller) DeleteAllfileFav(ctx context.Context, fileID string) (*pb.DeleteAllfileFavResponse, error) {
	filter := bson.D{
		bson.E{
			Key:   FavoriteBSONFileIDField,
			Value: fileID,
		},
	}

	result, err := c.store.DeleteAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteAllfileFavResponse{Acknownledged: result.acknownledged, DeletedCount: result.deletedCount}, nil

}