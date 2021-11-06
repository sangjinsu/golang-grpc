package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
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

func (*server) Maximum(stream calculatorpb.CalculatorService_MaximumServer) error {
	log.Printf("Received Maximum RPC\n")

	maxNum := int64(math.MinInt64)
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading clinet stream: %v", err)
			return err
		}

		number := request.GetNumber()
		if maxNum < number {
			maxNum = number
			sendErr := stream.Send(&calculatorpb.MaximumResponse{
				Result: maxNum,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("Received a negative number %v", number))
	}
	return &calculatorpb.SquareRootResponse{
		Result: math.Sqrt(float64(number)),
	}, nil
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
