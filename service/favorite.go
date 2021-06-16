package service

import (

	pb "github.com/meateam/fav-service/proto"
)

type Favorite interface {
	GetFileID() string

	SetFileID(fileID string) error

	GetUserID() string

	SetUserID(userID string) error

	MarshalProto(permission *pb.FavoriteObject) error
}