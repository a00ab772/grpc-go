package main

import (
	"context"
	"fmt"
	"log"

	pb "example.com/m/blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateBlog(ctx context.Context, in *pb.Blog) (*pb.BlogId, error) {
	log.Printf("CreateBlog was invoked with %v", in)

	data := BlogItem{
		AuthorID: in.AuthorId,
		Title:    in.Title,
		Content:  in.Content,
	}

	res, err := collection.InsertOne(ctx, data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("CreateBlog failed with internal error %v", err),
		)
	} else {
		log.Printf("CreateBlog was invoked with %v", in)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return nil, status.Errorf(codes.Internal, "CreateBlog cannot convert to OID")
	}
	return &pb.BlogId{
		Id: oid.Hex(),
	}, nil
}
