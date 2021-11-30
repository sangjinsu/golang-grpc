package main

import (
	"context"
	"fmt"
	"github.com/grpc-project/blog/blogpb"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Create blog request")
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	var data *blogItem
	decodeErr := collection.FindOne(context.Background(), bson.D{{"_id", oid}}).Decode(&data)
	if decodeErr != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}
	var data *blogItem
	filter := bson.D{{"_id", oid}}
	decodeErr := collection.FindOne(context.Background(), filter).Decode(&data)
	if decodeErr != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB %v", updateErr))
	}

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Delete blog request")
	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	filter := bson.D{{"_id", oid}}

	res, deleteErr := collection.DeleteOne(context.Background(), filter)
	if deleteErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot Delete object in MongoDB %v", deleteErr))
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find document in MongoDB %v", deleteErr))
	}

	return &blogpb.DeleteBlogResponse{
		BlogId: req.GetBlogId(),
	}, nil
}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		closeErr := cursor.Close(ctx)
		if closeErr != nil {
			fmt.Printf("Cursor close error: %v\n", closeErr)
		}
	}(cursor, context.Background())

	for cursor.Next(context.Background()) {
		data := &blogItem{}
		decodeErr := cursor.Decode(data)
		if decodeErr != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", decodeErr),
			)
		}
		stream.Send(&blogpb.ListBlogResponse{
			Blog: dataToBlogPb(data),
		})
	}

	if cursorErr := cursor.Err(); cursorErr != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}

	return nil
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
