package server

import (
	"context"
	"fmt"

	// "log"
	"net"
	"time"

	ilogger "github.com/meateam/elasticsearch-logger"
	pb "github.com/meateam/fav-service/proto"
	"github.com/meateam/fav-service/service"
	"github.com/meateam/fav-service/service/mongodb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	configPort                         = "port"
	configMongoConnectionString        = "mongo_host"
	// configMongoClientConnectionTimeout = "mongo_client_connection_timeout"
	configMongoClientPingTimeout       = "mongo_client_ping_timeout"

)

func init() {
	viper.SetDefault(configPort, "8080")
	viper.SetDefault(configMongoConnectionString, "mongodb://localhost:27017/favorite")
	// viper.SetDefault(configMongoClientConnectionTimeout, 10)
	viper.SetDefault(configMongoClientPingTimeout, 10)
	
}

type FavoriteServer struct {
	*grpc.Server
	logger            *logrus.Logger
	port              string
	favoriteService   service.Service
}

func (s FavoriteServer) Serve(lis net.Listener) {

}


//NewServer configures and creates a grpc.Server instance 
//returns FavoriteServer 
func NewServer(logger *logrus.Logger) *FavoriteServer {

	// If no logger is given, create a new default logger for the server.
	if logger == nil {
		logger = ilogger.NewLogger()
	}


	// Create a new grpc server.
	grpcServer := grpc.NewServer()

	mongoDBcontroller, err := initMongoDBcontroller(viper.GetString(configMongoConnectionString))
	if err != nil {
		logger.Fatalf("%v", err)
	}

	// fav := &pb.FavoriteServer
	favoriteService := service.NewService(mongoDBcontroller, logger)
	pb.RegisterFavoriteServer(grpcServer, favoriteService)

	// reflection.Register(grpcServer)


}











func initMongoDBcontroller(connectionString string) (service.Controller, error)  {
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
//Create mongodb client:
	
	//SetMonitor specifies a CommandMonitor (a monitor that is triggered for different events) to receive command events 
	mongoOptions := options.Client().ApplyURI(connectionString).SetMonitor(apmmongo.CommandMonitor())
	mongoClient, err := mongo.NewClient(mongoOptions)
	if err != nil {
		return nil, fmt.Errorf("failed creating mongodb client with connection string %s: %v", connectionString, err)
	}

// Connect client to mongodb:

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to mongodb with connection string %s: %v", connectionString, err)
	}

	defer mongoClient.Disconnect(ctx)

// check the connection: (copied)
	mongoClientPingTimeout := viper.GetDuration(configMongoClientPingTimeout)
	pingTimeoutCtx, cancelPing := context.WithTimeout(context.TODO(), mongoClientPingTimeout*time.Second)
	defer cancelPing()
	err = mongoClient.Ping(pingTimeoutCtx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed pinging to mongodb with connection string %s: %v", connectionString, err)
	}

	return mongoClient, nil

}








func getMongoDatabaseName(mongoClient *mongo.Client, connectionString string) (*mongo.Database, error) {
	connString, err := connstring.Parse(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed parsing connection string %s: %v", connectionString, err)
	}

	return mongoClient.Database(connString.Database), nil
}

















