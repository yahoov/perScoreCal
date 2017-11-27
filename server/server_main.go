package server

import (
	"fmt"
	"log"
	"net"
	pb "perScoreCal/perScoreProto/question"

	"google.golang.org/grpc"
)

const address = "0.0.0.0:6060"

// StartServer ...
func StartServer() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterQuestionServer(s, &Server{})
	fmt.Println("perScoreCal server started on port 6060 ...")
	s.Serve(lis)
}
