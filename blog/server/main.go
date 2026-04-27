package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "example.com/m/blog/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var collection *mongo.Collection
var addr = "0.0.0.0:5002"

type Server struct {
	pb.UnimplementedBlogServiceServer
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection = client.Database("blogdb").Collection("blog")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Blog failed to listen: %v", err)
	} else {
		log.Printf("Blog server listening at %v", lis.Addr())
	}

	srv := grpc.NewServer()
	pb.RegisterBlogServiceServer(srv, &Server{})

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Blog server failed to serve: %v", err)
	} else {
		log.Printf("Blog server started")
	}
}
