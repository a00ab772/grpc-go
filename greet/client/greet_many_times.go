package main

import (
	"context"
	"io"
	"log"

	proto2 "example.com/m/greet/proto"
)

func doGreetManyTimes(c proto2.GreetServiceClient) {
	log.Printf("gPRC SERVER STREAMING - Client side implementation")
	log.Printf("==================================================")
	log.Printf("")
	log.Println("start doGreetManyTimesClient")

	req := &proto2.GreetRequest{
		FirstName: "Arturo",
	}

	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes: %v", err)
		log.Fatalf("%v.GreetManyTimes(_) = _, %v", c, err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error while reading stream: %v", err)
			log.Fatalf("%v.GreetManyTimes(_) = _, %v", c, err)
		}

		log.Printf("GreetManyTimes (_) = _, %v\n", msg.Result)

	}

}
