package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/ExtraWhy/internal-libs/proto-models"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.PlayerRequest) (*pb.PlayerResponse, error) {
	v := "helllllo"
	log.Printf("Received: %v", in.GetName())
	return &pb.PlayerResponse{Name: &v}, nil
}

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &server{}
	pb.RegisterGreeterServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
