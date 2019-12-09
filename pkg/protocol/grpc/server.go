package grpc

import (
	"context"
	v1 "github.com/donggangcj/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"github.com/donggangcj/go-grpc-http-rest-microservice-tutorial/pkg/logger"
	"github.com/donggangcj/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc/middleware"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

//RunServer run gRPC service to public ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server statup option
	var opts []grpc.ServerOption

	// add middleware
	opts = middleware.AddLogging(logger.Log, opts)

	// register service
	server := grpc.NewServer(opts...)
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Log.Info("shutting down gRPC server ....")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Info("starting gRPC server ...")
	return server.Serve(listen)
}
