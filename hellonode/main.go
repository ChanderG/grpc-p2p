package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	hs "github.com/chanderg/grpc-p2p/helloservice"
	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type node struct {
	// Self information
	Name string
	Addr string

	// Consul related variables
	SDAddress string
	SDKV      api.KV

	// used to make requests
	Clients map[string]hs.HelloServiceClient
}

// SayHello implements helloworld.GreeterServer
func (n *node) SayHello(ctx context.Context, in *hs.HelloRequest) (*hs.HelloReply, error) {
	return &hs.HelloReply{Message: "Hello from" + n.Name}, nil
}

// Start listening/service.
func (n *node) StartListening() {

	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	_n := grpc.NewServer() // n is for serving purpose

	hs.RegisterHelloServiceServer(_n, n)
	// Register reflection service on gRPC server.
	reflection.Register(_n)

	// start listening
	if err := _n.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Register self with the service discovery module.
func (n *node) registerService() {
	config := api.DefaultConfig()
	config.Address = n.SDAddress
	consul, err := api.NewClient(config)
	if err != nil {
		log.Panicln("Unable to contact Service Discovery.")
	}

	kv := consul.KV()
	p := &api.KVPair{Key: n.Name, Value: []byte(n.Addr)}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Panicln("Unable to register with Service Discovery.")
	}

	// store the kv for future use
	n.SDKV = *kv

	log.Println("Successfully registered with Consul.")
}

// Start the node.
// This starts listening at the configured address. It also sets up clients for it's peers.
func (n *node) Start() {

	// start service / listening
	go n.StartListening()

	// register with the service discovery unit
	n.registerService()

	// start the main loop here
	// in our case, simply time out for 1 minute and greet all

	// wait for other nodes to come up
	for {
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	// pass the port as an argument and also the port of the other node
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("Arguments required: <name> <listening address> <consul address>")
		os.Exit(1)
	}

	// args in order
	name := args[0]
	listenaddr := args[1]
	sdaddress := args[2]

	noden := node{Name: name, Addr: listenaddr, SDAddress: sdaddress, Clients: nil} // noden is for opeartional purposes

	// start the node
	noden.Start()

	/*
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
	*/
}
