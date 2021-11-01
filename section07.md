
# gRPC Golang Master Class

section 07 gRPC Client Streaming

## What's a Client Streaming API?

- Client Streaming RPC API are a NEW kind API enabled thanks to HTTP/2

- The client will send many message to the server and will receive one response from the server 

- Streaming Client are well suited for When the client needs to send a lot of data

  - When the Server processing is expensive and should happen as the client sends data

  - When the client needs to PUSH data to the server without really expecting a response

    

- In gRPC Client Streaming Calls are defined using the keyword "stream"



## Server

```go
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}
```

## Client

```go
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
```

## proto

```protobuf
message LongGreetRequest {
  Greeting greeting = 1;
}

message LongGreetResponse {
  string result = 1;
}

service GreetService{
  // Unary
  rpc Greet(GreetRequest) returns (GreetResponse) {};

  // Server Streaming
  rpc GreetManyTimes(GreetManyTimesRequest) returns (stream GreetManyTimesResponse) {}

  // Client Streaming
  rpc LongGreet(stream LongGreetRequest) returns  (LongGreetResponse) {}
}
```

