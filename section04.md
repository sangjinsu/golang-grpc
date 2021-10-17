# gRPC Golang Master Class

section04

## Code Generation

- greet.proto

  ```protobuf
  syntax = "proto3";
  
  package greet;
  option go_package = "./greet/greetpb";
  
  service GreetService{}
  ```

- command

  ```sh
  protoc ./greet/greetpb/greet.proto --go_out=plugins=grpc:.
  ```


## Server Boilerplate

```go
package main

import (
	"fmt"
	"github.com/grpc-project/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
}

func main() {
	fmt.Println("hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
```

## Client Boilerplate

```go
package main

import (
   "fmt"
   "github.com/grpc-project/greet/greetpb"
   "google.golang.org/grpc"
   "log"
)

func main() {
   fmt.Println("Hello This Client")
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

   c := greetpb.NewGreetServiceClient(cc)
   fmt.Printf("Created client: %f", c)
}
```

