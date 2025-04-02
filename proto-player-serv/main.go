package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"

	pb "github.com/ExtraWhy/internal-libs/proto-models"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedServiceGameWonServer
}

func (s *server) GetWinForPlayer(ctx context.Context, in *pb.PlayerRequest) (*pb.PlayerResponse, error) {
	v := fmt.Sprintf("Hello %s ", in.GetName())
	won := rand.Uint64N(2)
	log.Printf("Received: %v", in.GetName())
	autor := &pb.PlayerResponse{}
	autor.Name = &v
	autor.MoneyWon = &won
	return autor, nil
}

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &server{}
	pb.RegisterServiceGameWonServer(s, srv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
