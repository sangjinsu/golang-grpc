package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/calculator/calculatorpb"
	"google.golang.org/grpc"
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


	for i := 1; i < 10; i++ {
		wg.Add(1)
		go doUnary(c, i, i+1)
	}
	wg.Wait()
}

func doUnary(c calculatorpb.CalculatorServiceClient, n1, n2 int) {
	defer wg.Done()
	fmt.Println("Starting to do a Unary RPC")
	req := &calculatorpb.SumRequest{
		FirstNumber: int64(n1),
		SecondNumber: int64(n2),
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}
