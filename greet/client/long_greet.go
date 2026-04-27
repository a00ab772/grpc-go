package main

import (
	"context"
	"log"
	"time"

	pb "example.com/m/greet/proto"
)

func doLongGreet(c pb.GreetServiceClient) {
	log.Printf("gPRC CLIENT STREAMING - Client side implementation")
	log.Printf("==================================================")
	log.Printf("")
	log.Printf("calling doLongGreet with stream\n")

	reqs := []*pb.GreetRequest{
		{FirstName: "Arturo"},
		{FirstName: "Bob"},
		{FirstName: "John"},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
		return
	}
	for _, req := range reqs {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet: %v", err)
		return
	} else {
		log.Printf("Response from LongGreet: %v\n", res.Result)
	}
}
