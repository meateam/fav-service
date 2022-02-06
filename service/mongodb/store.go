package mongodb

import (
	"context"
	"fmt"

	"github.com/meateam/fav-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

const (
	// MongoObjectIDField is the default mongodb unique key.
	MongoObjectIDField = "_id"

	// FavoriteCollectionName is the name of the favorites collection.
	FavoriteCollectionName = "favorite"	

	// FavoriteBSONFileIDField is the name of the fileID field in BSON.
	FavoriteBSONFileIDField = "fileID"

	// FavoriteBSONUserIDField is the name of the userID field in BSON.
	FavoriteBSONUserIDField = "userID"
)

// MongoStore holds the mongodb database and implements Store interface.
type MongoStore struct {
	DB *mongo.Database
}


// DeletedFavRes is the struct of the result from the delete all requests
type DeletedFavRes struct {
	acknownledged 	bool
	deletedCount 	int32
}


// newMongoStore returns a new store.
func newMongoStore(db *mongo.Database) (MongoStore, error) {
	collection := db.Collection(FavoriteCollectionName)
	indexes := collection.Indexes()
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			bson.E{
				Key:   FavoriteBSONFileIDField,
				Value: 1,
			},
			bson.E{
				Key:   FavoriteBSONUserIDField,
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := indexes.CreateOne(context.Background(), indexModel)
	if err != nil {
		return MongoStore{}, err
	}

	return MongoStore{DB: db}, nil
}



// GetAll gets all user favorite files by userID. 
// If the user doesnt have favorite files at all, it will return empty array. 
// If successful returns an array of favorite objects. 
func (s MongoStore) GetAll(ctx context.Context, filter interface{}) ([]service.Favorite, error) {
	collection := s.DB.Collection(FavoriteCollectionName)

	filterCursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer filterCursor.Close(ctx)

	favFiles := []service.Favorite{}
	for filterCursor.Next(ctx) {
		fav := &BSON{}
		err := filterCursor.Decode(fav)
		if err != nil {
			return nil, err
		}
		favFiles = append(favFiles, fav)
	}

	if err := filterCursor.Err(); err != nil {
		return nil, err
	}

	return favFiles, nil

}


// Create creates a favorite object of userID and fileID.
// If favorite already exists then it will return nil and error.
// If successful returns the favorite obejct and a nil error. 
func (s MongoStore) Create(ctx context.Context, favorite service.Favorite,) (service.Favorite, error) {
	collection := s.DB.Collection(FavoriteCollectionName)

	fileID := favorite.GetFileID()
	userID := favorite.GetUserID()


	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	favObject := bson.D{
		{Key: "userID", Value: userID},
		{Key: "fileID", Value: fileID},
	}

	_, err := collection.InsertOne(ctx, favObject)
	if err != nil {
		return nil, err
	}

	result := collection.FindOne(ctx, favObject)
	favoriteRes := &BSON{}
	err = result.Decode(favoriteRes)
	if err != nil {
		return nil, err
	}

	return favoriteRes, nil
}


// Delete deletes a favorite by userID and fileID. 
// If favorite does not exists it will return nil and error. 
// If successful returns the deleted favorite object. 
func (s MongoStore) Delete(ctx context.Context, filter interface{}) (service.Favorite, error){
	collection := s.DB.Collection(FavoriteCollectionName)

	result := collection.FindOneAndDelete(ctx, filter)

	deletedFav := &BSON{}
	err := result.Decode(deletedFav)

	if err != nil {
		return nil, err
	}

	return deletedFav, nil

}

// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s MongoStore) HealthCheck(ctx context.Context) (bool, error) {
	if err := s.DB.Client().Ping(ctx, readpref.Primary()); err != nil {
		return false, err
	}

	return true, nil
}


// Get finds one favorite that matches filter,
// if successful returns the favorite, and a nil error,
// if the favorite is not found it would return nil and NotFound error,
// otherwise returns nil and non-nil error if any occurred.
func (s MongoStore) Get(ctx context.Context, filter interface{}) (service.Favorite, error) {
	collection := s.DB.Collection(FavoriteCollectionName)
	favorite := &BSON{}
	foundFav := collection.FindOne(ctx, filter)
	err := foundFav.Decode(favorite)
	if err != nil {
		return nil, err
	}
	return favorite, nil
}


// DeleteAll deletes all favorites that matches filter
// if successful returns the deletedFavRes type with acknownledged true and deletedCount of favorites and nil
// otherwise returns the deletedFavRes type with acknownledged false and 0 deletedCount and nil
func (s MongoStore) DeleteAll(ctx context.Context, filter interface{}) (*DeletedFavRes, error) {
	collection := s.DB.Collection(FavoriteCollectionName)
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return &DeletedFavRes{acknownledged: false, deletedCount: int32(result.DeletedCount)}, nil
	}

	return &DeletedFavRes{acknownledged: true, deletedCount: int32(result.DeletedCount)}, nil

}