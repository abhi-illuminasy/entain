package main

import (
	"flag"
	"net"

	"git.neds.sh/matty/entain/racing/proto/sport"
	"git.neds.sh/matty/entain/racing/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "localhost:9001", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("failed running grpc server: %s", err)
	}
}

func run() error {
	conn, err := net.Listen("tcp", ":9001")
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	sport.RegisterSportServer(
		grpcServer,
		service.NewSportService(),
	)

	log.Infof("gRPC server listening on: %s", *grpcEndpoint)

	if err := grpcServer.Serve(conn); err != nil {
		return err
	}

	return nil
}
