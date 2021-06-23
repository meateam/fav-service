package mongodb

import (
	"context"
	"fmt"

	"github.com/meateam/fav-service/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s MongoStore) GetAll(ctx context.Context, filter interface{}) ([]primitive.D, error) {

	collection := s.DB.Collection(FavoriteCollectionName)

	filterCursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var favFiles []bson.D
	if err = filterCursor.All(ctx, &favFiles); err != nil {
		return nil, err
	}

	return favFiles, nil
}

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

	collection.InsertOne(ctx, favObject)

	result := collection.FindOne(ctx, favObject)
	favoriteRes := &BSON{}
	er := result.Decode(favoriteRes)
	if er != nil {
		return nil, er
	}

	return favoriteRes, nil
}

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





