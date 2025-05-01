package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	gametest "proto/player/server/slot-game"
	"sync"

	"github.com/ExtraWhy/internal-libs/logger"
	pb "github.com/ExtraWhy/internal-libs/proto-models/player"
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

type server struct {
	pb.UnimplementedServiceGameWonServer
}

func (s *server) GetWinForCleopatra(context.Context, *pb.PlayerRequest) (*pb.CleopatraWins, error) {

	retwins := pb.CleopatraWins{}
	retwins.Wins = make([]*pb.CleopatraWin, 1)

	wins, cl := gametest.CleopatraSpin(100)

	for j := 0; j < 5; j++ {
		for i := 0; i < 3; i++ {
			retwins.Syms = append(retwins.Syms, int32(cl.Scr[j][i]))
		}
	}

	for _, j := range *wins {
		bid := uint32(j.BID)
		free := uint32(j.Free)
		jid := uint32(j.JID)
		jack := float32(j.Jack)
		line := uint32(j.Line)
		mult := float32(j.Mult)
		pay := float32(j.Pay)
		num := uint32(j.Num)
		sym := uint32(j.Sym)
		retwins.Wins = append(retwins.Wins, &pb.CleopatraWin{})

		retwins.Wins[len(retwins.Wins)-1].BID = &bid
		retwins.Wins[len(retwins.Wins)-1].Free = &free
		retwins.Wins[len(retwins.Wins)-1].JID = &jid
		retwins.Wins[len(retwins.Wins)-1].Jack = &jack
		retwins.Wins[len(retwins.Wins)-1].Line = &line
		retwins.Wins[len(retwins.Wins)-1].Mult = &mult
		retwins.Wins[len(retwins.Wins)-1].Pay = &pay
		retwins.Wins[len(retwins.Wins)-1].Num = &num
		retwins.Wins[len(retwins.Wins)-1].Sym = &sym

		for h := 0; h < len(j.XY); h++ {
			retwins.Wins[len(retwins.Wins)-1].Linex = append(retwins.Wins[len(retwins.Wins)-1].Linex, uint32(j.XY[h]))
		}

	}
	return &retwins, nil

}

func (s *server) GetWinForPlayer(ctx context.Context, in *pb.PlayerRequest) (*pb.PlayerResponse, error) {
	var m0 uint64 = 0
	v := fmt.Sprintf("%s ", in.GetName())
	mult, lines, reels := gametest.RollLines()

	id := in.GetId()
	autor := &pb.PlayerResponse{}
	autor.Name = &v
	autor.MoneyWon = &m0
	autor.Id = &id
	autor.Lines = lines
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			autor.Reels = append(autor.Reels, byte(reels[j][i]))
		}
	}
	if lines != nil {
		autor.MoneyWon = &mult
	}
	return autor, nil
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "t" {
		gametest.SetupGame(true)
		log(logger.DEV, "setting up game in test mode (very high win ratio)")
	} else {
		log(logger.DEV, "setting up game in normal mode (normal win ratio)")
		gametest.SetupGame(false)
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log(logger.CRITICAL, "failed to listen", zap.String("address", lis.Addr().String()))
		os.Exit(-1)
		//Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	srv := &server{}
	pb.RegisterServiceGameWonServer(s, srv)
	if err := s.Serve(lis); err != nil {
		log(logger.CRITICAL, "failed to listen", zap.String("address", lis.Addr().String()))
		os.Exit(-1)
	}
	log(logger.DEV, "server listening", zap.String("addres", lis.Addr().String()))
}
