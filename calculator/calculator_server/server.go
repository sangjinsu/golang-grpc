package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.Request) (*calculatorpb.Response, error) {
	log.Printf("Request from RPC: %v %v", req.FirstNumber, req.SecondNumber)

	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calculatorpb.Response{
		Result: sum,
	}
	return res, nil
}

func (s *server) Minus(ctx context.Context, req *calculatorpb.Request) (*calculatorpb.Response, error) {
	log.Printf("Request from RPC: %v %v", req.FirstNumber, req.SecondNumber)

	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	minus := firstNumber - secondNumber
	res := &calculatorpb.Response{
		Result: minus,
	}
	return res, nil
}

func (s *server) Multiply(ctx context.Context, req *calculatorpb.Request) (*calculatorpb.Response, error) {
	log.Printf("Request from RPC: %v %v", req.FirstNumber, req.SecondNumber)

	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	mul := firstNumber * secondNumber
	res := &calculatorpb.Response{
		Result: mul,
	}
	return res, nil
}

func (s *server) Divide(ctx context.Context, req *calculatorpb.Request) (*calculatorpb.Response, error) {
	log.Printf("Request from RPC: %v %v", req.FirstNumber, req.SecondNumber)

	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	if secondNumber == 0 {
		return nil, fmt.Errorf("second number is Zero: %v", secondNumber)
	}
	divide := firstNumber / secondNumber
	res := &calculatorpb.Response{
		Result: divide,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
