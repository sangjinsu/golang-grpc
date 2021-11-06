package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
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
	// fmt.Printf("Created client: %f", c)
	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi streaming RPC")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "jinsu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "john",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "heesu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "wanhee",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "hongja",
			},
		},
	}

	waitc := make(chan string)
	go func() {
		for _, request := range requests {
			log.Printf("Sending message: %v\n", request)
			stream.Send(request)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			response, responseErr := stream.Recv()
			if responseErr == io.EOF {
				break
			}
			if responseErr != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			waitc <- response.GetResult()
		}
		close(waitc)
	}()

	for res := range waitc {
		fmt.Printf("Received: %v", res)
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "jinsu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "john",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "heesu",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "wanhee",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "hongja",
			},
		},
	}

	for _, request := range requests {
		fmt.Printf("Sending request %v\n", request)
		stream.Send(request)
		time.Sleep(time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v\n", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "jinsu",
			LastName:  "sang",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes: %v\n")
	}
	var msg *greetpb.GreetManyTimesResponse
	for {
		msg, err = resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "jinsu",
			LastName:  "sang",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
