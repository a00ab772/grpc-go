package main

import (
	"io"
	"log"

	proto2 "example.com/m/greet/proto"
)

func (s *Server) GreetEveryone(stream proto2.GreetService_GreetEveryoneServer) error {
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
		err = stream.Send(&proto2.GreetResponse{
			Result: res,
		})

		if err != nil {
			log.Fatalf("Error in GreetEveryone while sending data to client: %v", err)
		}
	}
}
