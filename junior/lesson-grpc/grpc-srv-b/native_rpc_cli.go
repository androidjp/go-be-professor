package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type HelloServiceClient struct {
	client *rpc.Client
}

func (c *HelloServiceClient) Hello(name string) (string, error) {
	//args := &struct {
	//	Name string `json:"name"`
	//}{
	//	Name: name,
	//}
	var reply string
	err := c.client.Call("HelloService.Hello", name, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}

func NativeRPCClient(input string) {
	cli, err := rpc.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		log.Fatalf("dialing", err)
	}
	// Create a new HelloServiceClient object.
	helloServiceClient := &HelloServiceClient{client: cli}

	// Call the Hello method.
	greeting, err := helloServiceClient.Hello(input)
	if err != nil {
		log.Fatal("calling Hello failed", err)
	}

	// Print the greeting.
	fmt.Println(greeting)

	// 结束
	cli.Close()
}
