package main

import (
	"fmt"
	"io"
	"log"

	proto2 "example.com/m/greet/proto"
)

func (s *Server) LongGreet(stream proto2.GreetService_LongGreetServer) error {
	log.Printf("gPRC CLIENT STREAMING - server side implementation")
	log.Printf("==================================================")
	log.Printf("")
	log.Println("LongGreet start")

	res := ""

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("LongGreet End")
			return stream.SendAndClose(&proto2.GreetResponse{
				Result: res,
			})
		}
		if err != nil {
			log.Fatalf("LongGreet error while reading client stream with stream.Recv(): %v", err)
		}

		res += fmt.Sprintf("Hello %s!\n", req.FirstName)
	}
}
