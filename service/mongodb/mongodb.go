package mongodb

import (
	"fmt"

	pb "github.com/meateam/fav-service/proto"
)

// BSON is the structure that represents a favorite as it's stored.
type BSON struct {
	FileID    string             `bson:"fileID,omitempty"`
	UserID    string             `bson:"userID,omitempty"`
}

// FileID is the structure that represents a fileID as it's stored.
type FileID struct {
	FileID  []string			`bson:"fileID,omitempty"`
}

// GetFileID returns b.FileID.
func (b BSON) GetFileID() string {
	return b.FileID

}

// SetFileID sets b.FileID to fileID.
func (b *BSON) SetFileID(fileID string) error {
	if b == nil {
		panic("b == nil")
	}

	if fileID == "" {
		return fmt.Errorf("FileID is required")
	}

	b.FileID = fileID
	return nil

}

// GetUserID returns b.UserID.
func (b BSON) GetUserID() string {
	return b.UserID

}

// SetUserID sets b.UserID to userID.
func (b *BSON) SetUserID(userID string) error {
	if b == nil {
		panic("b == nil")
	}

	if userID == "" {
		return fmt.Errorf("UserID is required")
	}

	b.UserID = userID
	return nil

}

// MarshalProto marshals b into a favorite.
func (b BSON) MarshalProto(favorite *pb.FavoriteObject) error {
	favorite.FileID = b.GetFileID()
	favorite.UserID = b.GetUserID()

	return nil

}