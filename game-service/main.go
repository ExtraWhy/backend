package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	gametest "proto/player/server/game-test"
	"sync"

	"github.com/ExtraWhy/internal-libs/logger"
	pb "github.com/ExtraWhy/internal-libs/proto-models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
	zl   = logger.ZapperLog{}
	do   sync.Once
)

func log(level int, m string, zpf ...zap.Field) {
	do.Do(func() {
		zl.Init(1)
	})
	zl.Log(level, m, zpf...)
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedServiceGameWonServer
}

func (s *server) GetWinForPlayer(ctx context.Context, in *pb.PlayerRequest) (*pb.PlayerResponse, error) {
	var m0 uint64 = 0
	v := fmt.Sprintf("%s ", in.GetName())
	mult, lines := gametest.RollLines()
	id := in.GetId()
	autor := &pb.PlayerResponse{}
	autor.Name = &v
	autor.MoneyWon = &m0
	autor.Id = &id
	if lines != nil {
		autor.MoneyWon = &mult
		autor.Lines = lines
	}
	return autor, nil
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "t" {
		gametest.SetupGame(true)
		log(1, "setting up game in test mode (very high win ratio)")
	} else {
		log(1, "setting up game in normal mode (normal win ratio)")
		gametest.SetupGame(false)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log(3, "failed to listen", zap.String("address", lis.Addr().String()))
		//Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &server{}
	pb.RegisterServiceGameWonServer(s, srv)
	log(1, "server listening", zap.String("addres", lis.Addr().String()))
	if err := s.Serve(lis); err != nil {
		log(3, "failed to listen", zap.String("address", lis.Addr().String()))

	}
}
