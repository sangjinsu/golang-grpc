package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
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

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("Received PrimeNumberDecomposition from RPC: %v", req)
	number := req.GetNumber()
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Printf("Received Average RPC\n")
	var sum int64
	var cnt int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result := float64(sum) / float64(cnt)
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		number := req.GetNumber()
		sum += number
		cnt++
	}
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
