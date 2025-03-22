package main

import (
	"casino/rest-backend/server"
	"fmt"
)

func main() {
	fmt.Printf("--- server up ---\r\n")
	s := server.Server{}
	s.SetHostPort("localhost", 8080)
	s.DoRun()

}
