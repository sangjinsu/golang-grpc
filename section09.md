# gRPC Golang Master Class

section 09 gRPC Advanced Features Deep Dive

## Errors in gRPC

#### Error Codes

- It is common for API to sometimes return error code

- Http has many error codes

- While HTTP codes are standardized they're not usually clear

- gRPC has a few error codes

  - [grpc 공식 문서 error](https://www.grpc.io/docs/guides/error/)
  - [golang grpc error](https://avi.im/grpc-errors/#go)

  ![image-20211106154837705](section09.assets/image-20211106154837705.png)

  ![image-20211106154858627](section09.assets/image-20211106154858627.png)

  ![image-20211106154916061](section09.assets/image-20211106154916061.png)

### Client

```go
func doErrorCall(c calculatorpb.CalculatorServiceClient, n int64) {
	response, err := c.SquareRoot(
		context.Background(),
		&calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			log.Println(respErr.Message())
			log.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				log.Println("Negative Number Error")
				return
			}
		} else {
			// severe error
			log.Fatalf("Severe Error calling SquareRoot: %v", err)
			return
		}
	}
	log.Printf("Result of square root of %v : %v\n", n, response.GetResult())
}
```



### Server 

```go
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument, fmt.Sprintf("Received a negative number %v", number))
	}
```



## gRPC Deadlines

- Deadlines allow gRPC clients to specify how long they are willing to wait for an RPC to complete before the RPC is terminated with the error DEADLINE_EXCEEDED
- The gRPC documentation recommends set a deadline for all client RPC calls 

- The server should check if the deadline has exceeded and cancel the work it is doing

- [gRPC Deadline](https://grpc.io/blog/deadlines/)

### Server

```go
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("The client canceled the request")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
```

### Client

```go
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("Timeout was hit. Deadline was exceeded")
			} else {
				log.Fatalf("unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("error while calling Greet RPC: %v", err)
		}
		return
	}
	log.Printf("Response from Greet: %v", res.Result)
```

