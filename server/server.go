package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
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
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	configPort                  		= "port"
	envPrefix 							= "FVS"
	configHealthCheckInterval          	= "health_check_interval"
	configMongoConnectionString 		= "mongo_host"
	configMongoClientConnectionTimeout 	= "mongo_client_connection_timeout"
	configMongoClientPingTimeout       	= "mongo_client_ping_timeout"
	configElasticAPMIgnoreURLS         	= "elastic_apm_ignore_urls"


)

func init() {
	viper.SetDefault(configPort, "8080")
	viper.SetDefault(configHealthCheckInterval, 3)
	viper.SetDefault(configElasticAPMIgnoreURLS, "/grpc.health.v1.Health/Check")
	viper.SetDefault(configMongoConnectionString, "mongodb://localhost:27017/favorite")
	viper.SetDefault(configMongoClientConnectionTimeout, 10)
	viper.SetDefault(configMongoClientPingTimeout, 10)
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
}

// FavoriteServer is a structure that holds the favorite grpc server and its services and configuration. 
type FavoriteServer struct {
	*grpc.Server
	logger 		*logrus.Logger
	port 		string
	healthCheckInterval int
	favoriteService service.Service
}


// Serve accepts incoming connections on the listener `lis`, creating a new
// ServerTransport and service goroutine for each. The service goroutines
// read gRPC requests and then call the registered handlers to reply to them.
// Serve returns when `lis.Accept` fails with fatal errors. `lis` will be closed when
// this method returns.
// If `lis` is nil then Serve creates a `net.Listener` with "tcp" network listening
// on the configured `TCP_PORT`, which defaults to "8080".
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (s FavoriteServer) Serve(lis net.Listener) {
	listener := lis
	if lis ==nil {
		l, err := net.Listen("tcp", ":"+s.port)
		if err != nil {
			s.logger.Fatalf("failed to listen: %v", err)
		}

		listener = l
	}

	log.Printf("listening and serving grpc server on port %s", s.port)
	if err := s.Server.Serve(listener); err != nil {
		s.logger.Fatalf(err.Error())
	}

}


// NewServer configures and creates a grpc.Server instance
// health check service.
// Configure using environment variables.
// `HEALTH_CHECK_INTERVAL`: Interval to update serving state of the health check server.
// `PORT`: TCP port on which the grpc server would serve on.
func NewServer(logger *logrus.Logger) *FavoriteServer {
	if logger == nil {
		logger = ilogger.NewLogger()
	}

	// Set up grpc server opts with logger interceptor.
	serverOpts := append(
		serverLoggerInterceptor(logger),
		grpc.MaxRecvMsgSize(16<<20),
	)

	grpcServer := grpc.NewServer(
		serverOpts...,
	)

	controller, err := initMongoDBController()
	if err != nil {
		logger.Fatalf("%v", err)
	}

	favoriteService := service.NewService(controller, logger)
	pb.RegisterFavoriteServer(grpcServer, favoriteService)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)


	favoriteServer := &FavoriteServer{
		Server: grpcServer,
		logger: logger,
		port: viper.GetString(configPort),
		healthCheckInterval: viper.GetInt(configHealthCheckInterval),
		favoriteService: favoriteService,
	}

	// Health check validation goroutine worker.
	go favoriteServer.healthCheckWorker(healthServer)

	return favoriteServer

}

func initMongoDBController() (service.Controller, error) {
	mongoClient, err := connectToMongoDB(viper.GetString(configMongoConnectionString))
	if err != nil {
		return nil, err
	}

	db, err := getMongoDatabaseName(mongoClient, viper.GetString(configMongoConnectionString))
	if err != nil {
		return nil, err
	}

	controller, err := mongodb.NewMongoController(db)
	if err != nil {
		return nil, fmt.Errorf("failed creating mongo store: %v", err)
	}

	return controller, nil

}

// mongoOptions:
// SetMonitor specifies a CommandMonitor to receive command events.
// CommandMonitor represents a monitor that is triggered for different events.

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


func serverLoggerInterceptor(logger *logrus.Logger) []grpc.ServerOption {
	// Create new logrus entry for logger interceptor.
	logrusEntry := logrus.NewEntry(logger)

	ignorePayload := ilogger.IgnoreServerMethodsDecider(
		strings.Split(viper.GetString(configElasticAPMIgnoreURLS), ",")...,
	)

	ignoreInitialRequest := ilogger.IgnoreServerMethodsDecider(
		strings.Split(viper.GetString(configElasticAPMIgnoreURLS), ",")...,
	)

	loggerOpts := []grpc_logrus.Option {
		grpc_logrus.WithDecider(func(fullMethodName string, err error) bool {
			return ignorePayload(fullMethodName)
		}),
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}

	return ilogger.ElasticsearchLoggerServerInterceptor(
		logrusEntry,
		ignorePayload,
		ignoreInitialRequest,
		loggerOpts...,
	)

}



// healthCheckWorker is running an infinite loop that sets the serving status once
// in s.healthCheckInterval seconds.
func (s FavoriteServer) healthCheckWorker(healthServer *health.Server) {
	//GetDuration returns the value associated with the key as a duration.
	mongoClientPingTimeout := viper.GetDuration(configMongoClientPingTimeout)

	for {
		// HealthCheck checks the health of the service, returns true if healthy, or false otherwise
		if s.favoriteService.HealthCheck(mongoClientPingTimeout * time.Second) {
			healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
		} else {
			healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		}

		time.Sleep(time.Second * time.Duration(s.healthCheckInterval))
	}

}



