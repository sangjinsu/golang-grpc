package main

import (
	"fmt"
	"github.com/grpc-project/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct {
}

func main() {
	// if we crash the go code, we get the file name and line number
	// 버그나 에러 발생시 파일 이름과 줄 번호를 알 수 있다
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service Started")

	lis, listenErr := net.Listen("tcp", "0.0.0.0:50051")
	if listenErr != nil {
		log.Fatalf("Failed to listen: %v", listenErr)
	}

	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if serveErr := s.Serve(lis); serveErr != nil {
			log.Fatalf("Failed to serve: %v", serveErr)
		}
	}()

	// Wait for Control C to exit
	// ctrl + c 대기
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	closeErr := lis.Close()
	if closeErr != nil {
		log.Fatalf("Failed to close: %v", closeErr)
	}
	fmt.Println("End of Program")
}
