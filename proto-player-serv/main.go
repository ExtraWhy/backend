package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	gametest "proto/player/server/game-test"

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
	var m0 uint64 = 0
	v := fmt.Sprintf("Hello %s ", in.GetName())
	res := gametest.RollLines()
	autor := &pb.PlayerResponse{}
	autor.Name = &v
	autor.MoneyWon = &m0
	if res != nil {
		log.Printf("[%d][%d][%d][%d][%d][%d][%d][%d][%d][%d]\r\n",
			res.Top, res.Mid,
			res.Bottom, res.DHigh,
			res.DLow, res.ZigRight,
			res.ZizLeft, res.ZigDoubleLeft, res.ZigDoubleRight,
			res.ZigLongLeft)
		won := uint64(res.Bottom*25 + res.DHigh*10 + res.DLow*10 + res.Mid*100 + res.Top*25 + res.ZigRight*5 + res.ZizLeft*5 + res.ZigLongLeft*7)
		autor.MoneyWon = &won
		return autor, nil
	}
	return autor, nil
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "t" {
		gametest.SetupGame(true)
		fmt.Println("Setting up game in test mode (very high win ration)")
	} else {
		fmt.Println("Setting up game in normal game (win ratio as designed)")
		gametest.SetupGame(false)
	}

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
