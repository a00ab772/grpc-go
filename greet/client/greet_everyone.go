package main

import (
	"context"
	"io"
	"log"
	"time"

	proto2 "example.com/m/greet/proto"
)

func doGreetEveryone(c proto2.GreetServiceClient) {
	log.Printf("gPRC BI-DIRECTIONAL STREAMING - Client side implementation")
	log.Printf("==========================================================")
	log.Printf("")
	log.Printf("calling doGreetEveryone with stream\n")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("error in GreetEveryone while creating stream : %v", err)
	} else {
		log.Printf("GreetEveryone is streamed\n")
	}

	reqs := []*proto2.GreetRequest{
		{FirstName: "Arturo"},
		{FirstName: "Bob"},
		{FirstName: "John"},
	}

	wait_before_you_close_communication := make(chan struct{})

	go func() {
		for _, req := range reqs {
			log.Printf("sending a request to the server %s", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		log.Printf("waiting for a response from the server")

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while receiving response from the server : %v", err)
				break
			}

			log.Printf("GreetEveryone received %s\n", resp.Result)
		}
		close(wait_before_you_close_communication)
	}()
	<-wait_before_you_close_communication
}
