# gRPC Golang Master Class

section 09 gRPC Bi Directional Streaming

### Bi Directional Streaming API

- New kind API enabled thanks to HTTP/2
- The client will send many message to the server and will receive many responses from the server
- The number of requests and responses does not have to match
- Bi Directional Streaming RPC are will suited for 
  - When the client and server needs to send a lot of data asynchronously
  - Chat protocol
  - Long running connection

### Server

- 받은 숫자들 중에 최대값 클라이언트로 전송 

```go
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
```



### Client

```go
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
```

