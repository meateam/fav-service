package service

import (

	pb "github.com/meateam/fav-service/proto"
)


// Favorite is an interface of a favorite object (UserID, FileID)
type Favorite interface {
	GetFileID() string

	SetFileID(fileID string) error

	GetUserID() string

	SetUserID(userID string) error

	MarshalProto(permission *pb.FavoriteObject) error
}