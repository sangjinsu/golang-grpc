package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"sync"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("Calculator Client")

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

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//
	//for i := 1; i < 10; i++ {
	//	wg.Add(1)
	//	go doUnary(c, float64(i), float64(i+1))
	//}
	//wg.Wait()
	//
	//doServerStreaming(c)

	// doClientStreaming(c)

	doBiDiStreaming(c)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Maximum Server Client Streaming RPC...")

	stream, err := c.Maximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling Maximum RPC: %v", err)
	}

	nums := []int64{
		4, 5, 8, 1, 2, 6, -7, 89, 5, 2, 100, -7, 5,
	}

	waitc := make(chan int64)
	go func() {
		for _, num := range nums {
			request := &calculatorpb.MaximumRequest{Number: num}
			log.Printf("Sending message: %v\n", request)
			sendErr := stream.Send(request)
			if sendErr != nil {
				log.Fatalf("Error while sending: %v", sendErr)
				break
			}
		}
		closeSendErr := stream.CloseSend()
		if closeSendErr != nil {
			log.Fatalf("Error while closing sending: %v", closeSendErr)
		}
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
			waitc <- response.Result
		}
		close(waitc)
	}()

	for v := range waitc {
		log.Printf("Maxnum is %v\n", v)
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Average Server Client Streaming RPC...")

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error while calling Average RPC: %v", err)
	}

	requests := []*calculatorpb.AverageRequest{
		{Number: 5},
		{Number: 6},
		{Number: 7},
		{Number: 8},
	}
	for _, request := range requests {
		stream.Send(request)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is %v", res.GetResult())
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 17,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Prime Number Decomposition RPC: %v", err)
	}

	var res *calculatorpb.PrimeNumberDecompositionResponse
	for {
		res, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happend: %v\n", err)
		}
		log.Printf("Response from Prime Number Decomposition: %v", res.GetPrimeFactor())
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient, n1, n2 float64) {
	defer wg.Done()
	fmt.Println("Starting to do a Unary RPC")
	req := &calculatorpb.Request{
		FirstNumber:  n1,
		SecondNumber: n2,
	}
	sum, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", sum.Result)

	minus, err := c.Minus(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Minus RPC: %v", err)
	}
	log.Printf("Response from Minus: %v", minus.Result)

	mul, err := c.Multiply(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Multiply RPC: %v", err)
	}
	log.Printf("Response from Multiply: %v", mul.Result)

	divide, err := c.Divide(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Divide RPC: %v", err)
	}
	log.Printf("Response from Divide: %v", divide.Result)
}
