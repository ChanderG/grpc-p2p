#grpc-p2p

Simple hello world p2p setup using grpc for connections.

The default grpc tutorial base focus on server-client connections. What is presented here, is a peer to peer model, where 2 or more peers/nodes implement a larger service by working together. The grpc service in this context is what the peers expose to each other. By definition, this is identical; that is all peers implement the same grpc service and also are clients to the service running at other nodes.

Uses Consul as a Service Discovery mechanism.

## Running a p2p system

Start the consul agent:
```
consul agent -dev
```
Consul should now be running at default IP:port, "localhost:8500".

Start peer "Node 1":
```
 go run hellonode/main.go "Node 1" :5000 localhost:8500
```

Start peer "Node 2":
```
 go run hellonode/main.go "Node 1" :5001 localhost:8500
```

The peers should discover each other and invoke each other's placeholder service request.
