package main

import (
	"fmt"
	"log"

	pb "example.com/m/greet/proto"
	"google.golang.org/grpc"
)

func (s *Server) GreetManyTimes(in *pb.GreetRequest, stream grpc.ServerStreamingServer[pb.GreetResponse]) error {
	log.Printf("gPRC SERVER STREAMING - server side implementation")
	log.Printf("==================================================")
	log.Printf("")
	log.Printf("GreetManyTimes function was invoked with %v\n", in)

	for i := 0; i < 10; i++ {
		res := fmt.Sprintf("Hello %s, this is the invoke number %d", in.FirstName, i)
		err := stream.Send(&pb.GreetResponse{
			Result: res,
		})

		if err != nil {
			return err
		}
	}
	return nil
}
