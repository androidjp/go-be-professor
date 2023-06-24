package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type HelloService struct{}

func (h *HelloService) Hello(name string, reply *string) error {
	*reply = fmt.Sprintf("Hello, %s!", name)
	return nil
}

func NativeRPCServer() {

	// Create a new HelloService object.
	helloService := &HelloService{}

	// Register the HelloService object with the RPC server.
	_ = rpc.Register(helloService)

	// Listen for connections on port 8080.
	listener, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept connections and handle requests.
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Handle the request.
		rpc.ServeConn(conn)
	}
}
