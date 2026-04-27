package main

import (
	"io"
	"log"

	pb "example.com/m/greet/proto"
)

func (s *Server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	log.Printf("gPRC BI-DIRECTIONAL STREAMING - server side implementation")
	log.Printf("==========================================================")
	log.Printf("")
	log.Println("GreetEveryone start")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while reading stream in GreetEveryone: %v", err)
		}

		res := "Hello " + req.FirstName + "!"
		err = stream.Send(&pb.GreetResponse{
			Result: res,
		})

		if err != nil {
			log.Fatalf("Error in GreetEveryone while sending data to client: %v", err)
		}
	}
}
