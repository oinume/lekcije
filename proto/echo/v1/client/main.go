package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oinume/lekcije/proto-gen/go/proto/echo/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address        = "localhost:50051"
	defaultMessage = "Hello world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := echo.NewEchoClient(conn)

	// Contact the server and print out its response.
	message := defaultMessage
	if len(os.Args) > 1 {
		message = os.Args[1]
	}
	r, err := c.Echo(context.Background(), &echo.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("could not echo: %v", err)
	}
	fmt.Printf("%s\n", r.Message)
}
