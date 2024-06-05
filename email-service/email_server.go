package main

import (
	"email/handlers"
	"email/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("[Email Service] failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEmailServiceServer(s, &handlers.Server{})
	log.Printf("[Email Service] Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[Email Service] failed to serve: %v", err)
	}
}
