package main

import (
	"context"
	"log"

	pb "example.com/m/greet/proto"
)

func (s *Server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("gPRC UNARY - server side implementation")
	log.Printf("=======================================")
	log.Printf("")
	log.Printf("Greet called with req: %v", req)
	return &pb.GreetResponse{
		Result: "Hello " + req.FirstName,
	}, nil
}
