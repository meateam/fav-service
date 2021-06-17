package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

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
	// ilogger "github.com/meateam/elasticsearch-logger"
)

const (
	configPort                         = "port"
	configMongoConnectionString        = "mongo_host"
	configMongoClientConnectionTimeout = "mongo_client_connection_timeout"
	configMongoClientPingTimeout       = "mongo_client_ping_timeout"

)

func init() {
	viper.SetDefault(configPort, "8080")
	viper.SetDefault(configMongoConnectionString, "mongodb://localhost:27017/permission")
	viper.SetDefault(configMongoClientConnectionTimeout, 10)
	viper.SetDefault(configMongoClientPingTimeout, 10)
	
}

type FavServer struct {
	*grpc.Server
	logger            *logrus.Logger
	port              string
	favoriteService   service.Service
}

func (s FavServer) Serve(lis net.Listener) {

}

func NewServer(logger *logrus.Logger) *FavServer {
	grpcServer := grpc.NewServer()

	// controller, err := 

	// favService := service.NewService()
}

func connectToMongoDB(connectionString string) (*mongo.Client, error) {
	// Create mongodb client.
	mongoOptions := options.Client().ApplyURI(connectionString).SetMonitor(apmmongo.CommandMonitor())
	mongoClient, err := mongo.NewClient(mongoOptions)
	if err != nil {
		return nil, fmt.Errorf("failed creating mongodb client with connection string %s: %v", connectionString, err)
	}

	// Connect client to mongodb.
	mongoClientConnectionTimout := viper.GetDuration(configMongoClientConnectionTimeout)
	connectionTimeoutCtx, cancelConn := context.WithTimeout(context.TODO(), mongoClientConnectionTimout*time.Second)
	defer cancelConn()
	err = mongoClient.Connect(connectionTimeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to mongodb with connection string %s: %v", connectionString, err)
	}

	// Check the connection.
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
