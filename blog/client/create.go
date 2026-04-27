package main

import (
	"context"
	"log"

	pb "example.com/m/blog/proto"
)

func createBlog(client pb.BlogServiceClient) string {
	log.Println("call client.CreateBlog")

	blog := &pb.Blog{
		Id:      "Arturo",
		Title:   "How are you?",
		Content: "This is my blog",
	}

	result, err := client.CreateBlog(context.Background(), blog)
	if err != nil {
		log.Fatal("Unable to create blog", err)
	}

	log.Printf("Created blog with id: %s", result.Id)
	return result.Id

}
