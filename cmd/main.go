package main

import (
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	pb "github.com/ovargas/api-go/item/v1"
	"github.com/ovargas/api-go/storage/v1"
	"github.com/ovargas/item-api/internal/intrastructure"
	"github.com/ovargas/item-api/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	port           = flag.Int("port", 10001, "The service port")
	storageAddress = flag.String("storage_address", "localhost:10000", "The storage service address in the format of host:port")
)

func main() {
	flag.Parse()

	//Configuring the logging
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logEntry := logrus.NewEntry(logger)
	grpc_logrus.ReplaceGrpcLogger(logEntry)

	//Configuring logrus for gRCP
	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	//Creating a gRCP service and registering the middlewares
	grpcServer := grpc.NewServer(grpc_middleware.WithUnaryServerChain(
		grpc_logrus.UnaryServerInterceptor(logEntry, opts...),
	))

	//Creating a grcp connection to create a storage client
	storageConnection, err := grpc.Dial(*storageAddress, grpc.WithInsecure())
	defer func(storageConnection *grpc.ClientConn) {
		_ = storageConnection.Close()
	}(storageConnection)

	//Creating the storage client
	storageClient := storage.NewStorageServiceClient(storageConnection)

	// Registering the storage service
	pb.RegisterItemServiceServer(grpcServer,
		service.New(
			intrastructure.NewItemMemoryRepository(),
			storageClient),
	)

	if err != nil {
		log.Fatalf("unable to create storage: %v", err)
	}

	// The TCP listener where the service will be allocated
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	logger.Info("starting the service in the port ", *port)

	//Start the service
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("unable to start service: %v", err)
	}
}
