package main

import (
	"context"
	"log"

	pb "example.com/m/greet/proto"
)

func doGreet(c pb.GreetServiceClient) {
	log.Printf("gPRC UNARY - Client side implementation")
	log.Printf("=======================================")
	log.Printf("")
	log.Println("call doGreet")
	greet, err := c.Greet(context.Background(), &pb.GreetRequest{
		FirstName: "Arturo",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		panic(err)
	}
	log.Printf("Greeting: %v", greet.Result)
}
