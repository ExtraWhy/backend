package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/ExtraWhy/internal-libs/logger"
	pb "github.com/ExtraWhy/internal-libs/proto-models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "CryptoWin"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("Name", defaultName, "Name to greet")
	zl   = logger.ZapperLog{}
	do   sync.Once
)

func log(level int, m string, zpf ...zap.Field) {
	do.Do(func() {
		zl.Init(1)
	})
	zl.Log(level, m, zpf...)
}

type WinRequest struct {
	PlayerRequest  *pb.PlayerRequest
	PlayerResponse *pb.PlayerResponse
}

func (wr *WinRequest) SendWin(id uint64) error {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("No connection to game service ")
	}
	defer conn.Close()
	c := pb.NewServiceGameWonClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pr := &pb.PlayerRequest{Id: &id}
	wr.PlayerResponse, err = c.GetWinForPlayer(ctx, pr)
	if err != nil {
		return errors.New(fmt.Sprint("could not greet: %v", err))
	}
	log(1, "cryptowin proto", zap.Uint64("id", wr.PlayerResponse.GetId()),
		zap.Uint64("won", wr.PlayerResponse.GetMoneyWon()),
		zap.ByteString("lines", wr.PlayerResponse.GetLines()))
	return nil
}
