package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/meateam/fav-service/proto"
	"github.com/sirupsen/logrus"
)

// Service is a structure used for handling favorite Service grpc requests.
type Service struct {
	controller Controller
	logger     *logrus.Logger
	pb.UnimplementedFavoriteServer
}

// NewService creates a Service and returns it.
func NewService(controller Controller, logger *logrus.Logger) Service {
	return Service{controller: controller, logger: logger}

}


// CreateFavorite is the request handler for creating a favorite. 
func (s Service) CreateFavorite(ctx context.Context, req *pb.CreateFavoriteRequest,) (*pb.FavoriteObject, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	favorite, err := s.controller.CreateFavorite(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	var response pb.FavoriteObject
	if err := favorite.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil

}

// DeleteFavorite is the request handler for deleting favorite.
func (s Service) DeleteFavorite(ctx context.Context, req *pb.DeleteFavoriteRequest,) (*pb.FavoriteObject, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("fileID is required")
	}

	favorite, err := s.controller.DeleteFavorite(ctx, fileID, userID)
	if err != nil {
		return nil, err
	}

	var response pb.FavoriteObject
	if err := favorite.MarshalProto(&response); err != nil {
		return nil, err
	}

	return &response, nil

}

// GetAllFavoritesByUserID is the request handler for getting all user favorite files.  
func (s Service) GetAllFavoritesByUserID(ctx context.Context, req *pb.GetAllFavoritesRequest,) (*pb.GetAllFavoritesResponse, error) {
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	favorite, err := s.controller.GetAllFavoritesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &pb.GetAllFavoritesResponse{FavFileIDList: favorite}, nil

}


// HealthCheck checks the health of the service, returns true if healthy, or false otherwise.
func (s Service) HealthCheck(mongoClientPingTimeout time.Duration) bool {
	timeoutCtx, cancel := context.WithTimeout(context.TODO(), mongoClientPingTimeout)
	defer cancel()
	healthy, err := s.controller.HealthCheck(timeoutCtx)
	if err != nil {
		s.logger.Errorf("%v", err)
		return false
	}

	return healthy
	
}

// IsFavorite is the request handler for checking user favorite by userID and fileID.
func (s Service) IsFavorite(ctx context.Context, req *pb.IsFavoriteRequest) (*pb.IsFavoriteResponse, error) {
	fileID := req.GetFileID()
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("UserID is required")
	}

	if fileID == "" {
		return nil, fmt.Errorf("FileID is required")
	}

	_, err := s.controller.GetByFileAndUser(ctx, fileID, userID)
	if err != nil {
		return &pb.IsFavoriteResponse{IsFavorite: false}, err
	}

	return &pb.IsFavoriteResponse{IsFavorite: true}, err

}


// DeleteAllfileFav is the request handler for deleteing all favorite files by fileID.  
func (s Service) DeleteAllfileFav(ctx context.Context, req *pb.DeleteAllfileFavRequest) (*pb.DeleteAllfileFavResponse, error) {
	fileID := req.GetFileID()
	if fileID == "" {
		return nil, fmt.Errorf("FileID is required")
	}

	result, _ := s.controller.DeleteAllfileFav(ctx, fileID)
	return result, nil
}