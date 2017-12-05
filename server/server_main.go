package server

import (
	"fmt"
	"log"
	"net"
<<<<<<< HEAD
	qpb "perScoreCal/perScoreProto/question"
	upb "perScoreCal/perScoreProto/user"
=======
	pb "perScoreCal/perScoreProto/question"
>>>>>>> 5ed5fd7002ae0f8df7dddd9ef69b1ddc7987f3a2

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
<<<<<<< HEAD
	upb.RegisterUserServer(s, &UserServer{})
	qpb.RegisterQuestionServer(s, &QuestionServer{})
=======
	pb.RegisterQuestionServer(s, &Server{})
>>>>>>> 5ed5fd7002ae0f8df7dddd9ef69b1ddc7987f3a2
	fmt.Println("perScoreCal server started on port 6060 ...")
	s.Serve(lis)
}
