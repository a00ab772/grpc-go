package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "example.com/m/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr string = "0.0.0.0:5003"

type Server struct {
	pb.GreetServiceServer
}

func main() {
	wd, _ := os.Getwd()

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Listening on " + addr)
	}

	opts := []grpc.ServerOption{}
	tls := true
	if tls {
		log.Println("TLS is enabled")
		creds, err := credentials.NewServerTLSFromFile(wd+"/ssl/private/server.crt", wd+"/ssl/private/server.pem")
		if err != nil {
			log.Fatalf("Failed to read credentials %v", err)
		} else {
			log.Println("Successfully read server ssl credentials")
		}
		opts = append(opts, grpc.Creds(creds))
	}

	//s := grpc.NewServer()
	s := grpc.NewServer(opts...)

	pb.RegisterGreetServiceServer(s, &Server{})
	//s.Serve(lis)

	if err = s.Serve(lis); err != nil {
		panic(err)
	} else {
		fmt.Println("Listening on " + addr)
	}

}
