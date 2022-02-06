# Basic go commands
PROTOC=protoc

# Binary names
BINARY_NAME=fav-service

all: clean deps fmt test build
build: build-proto build-app 
test:
		docker-compose -f "docker-compose.yml" up -d minio && \
		docker-compose down && sudo rm -rf data
clean:
		go clean
		sudo rm -rf $(BINARY_NAME)
run: build
		./$(BINARY_NAME)
deps:
		go get -u github.com/golang/protobuf/protoc-gen-go
build-app:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-extldflags "-static"' -o $(BINARY_NAME) -v
build-proto:
		rm -f proto/*.pb.go
		protoc -I proto/ proto/*.proto --go_out=plugins=grpc:./proto

.PHONY: fmt
fmt:
	./gofmt.sh
