package server

import (
	"net"

	pb "github.com/meateam/fav-service/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

const (
	configMongoConnectionString        = "mongo_host"

)

func init() {
	viper.SetDefault(configMongoConnectionString, "mongodb://localhost:27017/favorite")
}

type FavServer struct {
	*grpc.Server
	logger            *logrus.Logger
	port              string
}

func (s FavServer) Serve(lis net.Listener) {
	listener := lis
	if lis == nil {
		l, err := net.Listen("tcp", ":"+s.port)
		if err != nil {
			s.logger.Fatalf("failed to listen: %v", err)
		}

		listener = l
	}

	s.logger.Infof("listening and serving grpc server on port %s", s.port)
	if err := s.Server.Serve(listener); err != nil {
		s.logger.Fatalf(err.Error())
	}
}

func connectToMongoDB(connectionString string) (*mongo.Client, error) {
	
}