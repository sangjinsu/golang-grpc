# gRPC Golang Master Class

section 05 gRPC Unary

## What's an Unary API?

- Unary RPC calls are the basic Request / Response that everyone is familiar with
- **The client will send one message to the server and will receive one response from the server**
- Unary RPC calls will be the most common for APIs

- In gRPC Unary Calls are defined using Protocol Buffers
- For each RPC call we have to define "Request" and "Response" message

## Greet API

- client

  ```go
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
  ```

- server

  ```go
  func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
  	fmt.Printf("Greet function was invoked with %v\n", req)
  	firstName := req.GetGreeting().GetFirstName()
  	lastName := req.GetGreeting().GetLastName()
  	result := fmt.Sprintf("Hello %s %s", firstName, lastName)
  	res := &greetpb.GreetResponse{Result: result}
  	return res, nil
  }
  ```

  

