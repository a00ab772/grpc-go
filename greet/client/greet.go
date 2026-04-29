package main

import (
	"context"
	"log"

	proto2 "example.com/m/greet/proto"
)

func doGreet(c proto2.GreetServiceClient) {
	log.Printf("gPRC UNARY - Client side implementation")
	log.Printf("=======================================")
	log.Printf("")
	log.Println("call doGreet")
	greet, err := c.Greet(context.Background(), &proto2.GreetRequest{
		FirstName: "Arturo",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		panic(err)
	}
	log.Printf("Greeting: %v", greet.Result)
}
