package service

import (
	"context"
	"fmt"

	pb "github.com/meateam/fav-service/proto"
	"github.com/sirupsen/logrus"
)

// Service is a structure used for handling favorite Service grpc requests.
type Service struct {
	pb.UnimplementedFavoriteServer
	controller Controller
	logger     *logrus.Logger
}

func NewService(controller Controller, logger *logrus.Logger) Service {
	return Service{controller: controller, logger: logger}
}

func (s Service) CreateFavorite(ctx context.Context, req *pb.CreateFavoriteRequest,) (*pb.FavoriteObject, error) {
	// fileID := req.FileID - what are the difference
	fileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	favorite, err := s.controller.CreateFavorites(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	var response pb.FavoriteObject
	if err := favorite.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s Service) DeleteFavorite(ctx context.Context, req *pb.DeleteFavoriteRequest,) (*pb.FavoriteObject, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	favorite, err := s.controller.DeleteFavorites(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	var response pb.FavoriteObject
	if err := favorite.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil

}

func (s Service) GetAllFavorite(ctx context.Context, req *pb.GetAllFavoriteRequest,) (*pb.GetAllFavoriteResponse, error) {
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	favorite, err := s.controller.GetAllFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.GetAllFavoriteResponse{FavFileList: favorite}, nil
}



