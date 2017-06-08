package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	hs "github.com/chanderg/grpc-p2p/helloservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type node struct {
	I int
	C hs.HelloServiceClient
}

// SayHello implements helloworld.GreeterServer
func (n *node) SayHello(ctx context.Context, in *hs.HelloRequest) (*hs.HelloReply, error) {
	return &hs.HelloReply{Message: "Hello " + strconv.Itoa(n.I)}, nil
}

func main() {
	// pass the port as an argument and also the port of the other node
	args := os.Args[1:]

	fmt.Println("Args: ", args[0])

	// example arg[0] -> :5000
	port := args[0]
	address := args[1]

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	n := grpc.NewServer()        // n is for serving purpose
	noden := node{I: 42, C: nil} // noden is for opeartional purposes

	hs.RegisterHelloServiceServer(n, &noden)
	// Register reflection service on gRPC server.
	reflection.Register(n)

	// start listening
	go func() {
		if err := n.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// wait for other nodes to come up
	time.Sleep(5 * time.Second)

	// setup connection with other node
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	noden.C = hs.NewHelloServiceClient(conn)

	r, err := noden.C.SayHello(context.Background(), &hs.HelloRequest{Name: "John"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting from the other node: %s", r.Message)

	for {
		time.Sleep(10 * time.Second)
	}
}
