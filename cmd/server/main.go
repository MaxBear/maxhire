package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	applicationspb "github.com/MaxBear/maxhire/proto/gen/go/applications/v1"
	"github.com/MaxBear/maxhire/server"
	"github.com/MaxBear/maxhire/service"
)

func main() {
	json := flag.String("json", "", "json file contains job application records")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	port := 9090
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Printf("error starting grpc server on port %d, error: %s", port, err.Error())
		os.Exit(1)
	}

	svc, err := service.NewService(ctx, *json)
	if err != nil {
		log.Printf("error starting grpc service, error: %s", err.Error())
		os.Exit(1)
	}

	srv := server.New(svc)
	grpcServer := grpc.NewServer()
	applicationspb.RegisterApplicationsServer(grpcServer, srv)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("error starting grpc server, error: %s", err.Error())
		os.Exit(1)
	}
}
