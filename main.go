package main

import (

	"github.com/meateam/fav-service/server"
)

func main() {
	server.NewServer(nil).Serve(nil)
}