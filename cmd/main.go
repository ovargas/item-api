package main

import (
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	pb "github.com/ovargas/api-go/item/v1"
	"github.com/ovargas/api-go/storage/v1"
	"github.com/ovargas/item-api/internal/intrastructure"
	"github.com/ovargas/item-api/internal/server"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	port           = flag.Int("port", 10001, "The server port")
	storageAddress = flag.String("storage_address", "localhost:10000", "The storage server address in the format of host:port")
)

func main() {

	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logEntry := logrus.NewEntry(logger)
	grpc_logrus.ReplaceGrpcLogger(logEntry)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	grpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		grpc_logrus.UnaryServerInterceptor(logEntry, opts...),
	))

	storageConnection, err := grpc.Dial(*storageAddress, grpc.WithInsecure())
	defer func(storageConnection *grpc.ClientConn) {
		_ = storageConnection.Close()
	}(storageConnection)

	pb.RegisterItemServiceServer(grpcServer,
		server.New(
			intrastructure.NewItemMemoryRepository(),
			storage.NewStorageServiceClient(storageConnection)),
	)

	if err != nil {
		log.Fatalf("unable to create storage: %v", err)
	}

	logger.Info("server start")
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
