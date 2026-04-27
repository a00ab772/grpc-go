package main

import (
	"fmt"
	"log"
	"os"

	pb "example.com/m/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr string = "localhost:5003"

func main() {
	wd, _ := os.Getwd()

	opts := []grpc.DialOption{}
	tls := true
	if tls {
		log.Println("TLS is enabled")
		creds, err := credentials.NewClientTLSFromFile(wd+"/ssl/public/ca.crt", "")
		if err != nil {
			log.Fatalf("Error while loading CA trust certificate %v", err)
		} else {
			log.Println("Successfully read ssl client credentials")
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	//conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(addr, opts...)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to " + addr)
	}

	// Close at the end of the function
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	doGreet(client)
	//doGreetManyTimes(client)
	//doLongGreet(client)
	//doGreetEveryone(client)
}
