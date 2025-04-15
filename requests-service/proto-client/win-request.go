package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/ExtraWhy/internal-libs/proto-models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "CryptoWin"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("Name", defaultName, "Name to greet")
)

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
	log.Printf("CryptoWin Proto : %d %d %v", wr.PlayerResponse.GetId(),
		wr.PlayerResponse.GetMoneyWon(), wr.PlayerResponse.GetLines())
	return nil
}
