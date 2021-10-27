
# gRPC Golang Master Class

section 06 gRPC Server Streaming

## What's a Server Streaming API?

- Server Steaming RPC API are a New kind API enabled thanks to HTTP/2
- The client will send one message to the server 
- The client will receive many responses from the server, possibly an infinite number

- Streaming Server are well suited for when the server needs to send a lot of data
- or When the server needs to PUSH data to the client without having the client request for more 



- In gRPC Server streaming Calls are defined using the keyword stream

  

### Server

```go
func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManytimes function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	// waitGroup 사용 
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(firstName string, i int) {
			defer wg.Done()
			result := "Hello " + firstName + " number " + strconv.Itoa(i)
			res := &greetpb.GreetManyTimesResponse{
				Result: result,
			}
			stream.Send(res)
		}(firstName, i)
	}
	wg.Wait()
	return nil
}
```



### Client

```go
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
```

