package main

import (
	"fmt"
	"github.com/grpc-project/greet/greetpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Hello This Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer func(cc *grpc.ClientConn) {
		err = cc.Close()
		if err != nil {
			log.Fatalf("Could not disconnect: %v", err)
		}
	}(cc)

	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created client: %f", c)
}
