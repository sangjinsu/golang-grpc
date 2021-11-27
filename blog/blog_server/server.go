package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/blog/blogpb"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
)

var collection *mongo.Collection

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

type server struct {
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Create blog request")
	blog := req.GetBlog()
	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID: %v", err))
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func main() {
	// if we crash the go code, we get the file name and line number
	// 버그나 에러 발생시 파일 이름과 줄 번호를 알 수 있다
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Connecting to MongoDB")

	client, dbClientErr := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_DB_URI")))
	if dbClientErr != nil {
		log.Fatalln(dbClientErr)
	}

	dbConnectErr := client.Connect(context.TODO())
	if dbConnectErr != nil {
		log.Fatalln(dbConnectErr)
	}

	fmt.Println("Blog Service Started")
	collection = client.Database("go-grpc").Collection("Blog")

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

	fmt.Println("Closing MongoDB Connection")
	if dbDisconnectErr := client.Disconnect(context.TODO()); dbDisconnectErr != nil {
		log.Fatalf("Error on disconnection with MongoDB : %v", dbDisconnectErr)
	}

	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("End of Program")
}
