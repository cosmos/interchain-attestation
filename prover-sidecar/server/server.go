package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"proversidecar/proof"
)

type Server struct {
	proof.UnimplementedProofServiceServer
}

func (s *Server) Serve(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	proof.RegisterProofServiceServer(server, s)
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) GetProof(ctx context.Context, request *proof.ProofRequest) (*proof.ProofResponse, error) {
	return &proof.ProofResponse{
		Proof: fmt.Sprintf("hello proof! (clientID: %s)", request.ClientID),
	}, nil
}

var _ proof.ProofServiceServer = &Server{}
