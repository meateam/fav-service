package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"time"

	pb "github.com/meateam/fav-service/proto"
	"github.com/meateam/fav-service/service"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	configPort                  = "port"
	configMongoConnectionString = "mongo_host"
)

func init() {
	viper.SetDefault(configPort, "8080")
	viper.SetDefault(configMongoConnectionString, "mongodb://127.0.0.1:27017/favorite")
}

type FavoriteServer struct {
	pb.UnimplementedFavoriteServer
	grpc.Server
	port 		string
	favoriteService service.Service
}


func NewServer() {

//connection to mongo:


	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString(configMongoConnectionString)))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)


//connection to grpc:

	lis, err := net.Listen("tcp", ":"+viper.GetString(configPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println(":"+viper.GetString(configPort))
	server := grpc.NewServer()

	pb.RegisterFavoriteServer(server, &FavoriteServer{})

	reflection.Register(server)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}



}



