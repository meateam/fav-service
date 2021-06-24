package server

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/meateam/fav-service/proto"
	"github.com/meateam/fav-service/service"
	"github.com/meateam/fav-service/service/mongodb"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmmongo"
	ilogger "github.com/meateam/elasticsearch-logger"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"google.golang.org/grpc"
)

const (
	configPort                  = "port"
	configMongoConnectionString = "mongo_host"
	configMongoClientConnectionTimeout = "mongo_client_connection_timeout"

)

func init() {
	viper.SetDefault(configPort, "8080")
	viper.SetDefault(configMongoConnectionString, "mongodb://127.0.0.1:27017/favorite")
	viper.SetDefault(configMongoClientConnectionTimeout, 10)
}

type FavoriteServer struct {
	pb.UnimplementedFavoriteServer
	grpc.Server
	logger 		*logrus.Logger
	port 		string
	favoriteService service.Service
}


func (s FavoriteServer) Serve(lis net.Listener) {
	listener := lis
	if lis ==nil {
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


func NewServer(logger *logrus.Logger) *FavoriteServer {
	if logger == nil {
		logger = ilogger.NewLogger()
	}

	grpcServer := grpc.NewServer()

	controller, err := initMongoDBController(viper.GetString(configMongoConnectionString))


	if err != nil {
		logger.Fatalf("%v", err)
	}

	favoriteService := service.NewService(controller, logger)
	pb.RegisterFavoriteServer(grpcServer, favoriteService)



	favoriteServer := &FavoriteServer{
		Server: *grpcServer,
		logger: logger,
		port: viper.GetString(configPort),
		favoriteService: favoriteService,
	}

	return favoriteServer

}


func initMongoDBController(connectionString string) (service.Controller, error) {
	mongoClient, err := connectToMongoDB(connectionString)
	if err != nil {
		return nil, err
	}

	db, err := getMongoDatabaseName(mongoClient, connectionString)
	if err != nil {
		return nil, err
	}

	controller, err := mongodb.NewMongoController(db)
	if err != nil {
		return nil, fmt.Errorf("failed creating mongo store: %v", err)
	}
	return controller, nil

}

func connectToMongoDB(connectionString string) (*mongo.Client, error) {
	mongoOptions := options.Client().ApplyURI(connectionString).SetMonitor(apmmongo.CommandMonitor())
	mongoClient, err := mongo.NewClient(mongoOptions)

	if err != nil {
		return nil, fmt.Errorf("failed creating mongodb client with connection string %s: %v", connectionString, err)
	}

	mongoClientConnectionTimout := viper.GetDuration(configMongoClientConnectionTimeout)
	connectionTimeoutCtx, cancelConn := context.WithTimeout(context.TODO(), mongoClientConnectionTimout*time.Second)
	defer cancelConn()
	err = mongoClient.Connect(connectionTimeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to mongodb with connection string %s: %v", connectionString, err)
	}

	// check the connection

	return mongoClient, nil
}


func getMongoDatabaseName(mongoClient *mongo.Client, connectionString string) (*mongo.Database, error) {
	connString, err := connstring.Parse(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed parsing connection string %s: %v", connectionString, err)
	}

	return mongoClient.Database(connString.Database), nil
}

