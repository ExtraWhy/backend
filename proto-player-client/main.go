package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/ExtraWhy/internal-libs/proto-models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "Pishki"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("Name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pr := pb.PlayerRequest{}
	pr.Name = name
	r, err := c.SayHello(ctx, &pr)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Pish mi yajkata : %s", r.GetName())
}
