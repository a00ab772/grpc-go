package main

import (
	"log"

	pb "example.com/m/blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "0.0.0.0:5002"

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln("Failed to connect:", err)
	}
	defer conn.Close()

	c := pb.NewBlogServiceClient(conn)

	createBlog(c)
}
