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

// GetAllFavorites is the request handler for getting all user favorite files.  
func (s Service) GetAllFavorites(ctx context.Context, req *pb.GetAllFavoritesRequest,) (*pb.GetAllFavoritesResponse, error) {
	userID := req.GetUserID()

	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	favorite, err := s.controller.GetAllFavorites(ctx, userID)
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

