package server

import (
	"fmt"
	"log"
	"net"
	qpb "perScoreCal/perScoreProto/question"
	upb "perScoreCal/perScoreProto/user"

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
	upb.RegisterUserServer(s, &UserServer{})
	qpb.RegisterQuestionServer(s, &QuestionServer{})
	fmt.Println("perScoreCal server started on port 6060 ...")
	s.Serve(lis)
}
