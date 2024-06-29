package server

import (
	"context"
	"fmt"
	"gitlab.com/tozd/go/errors"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	UnimplementedProofServer
}

func (s *Server) Serve(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrapf(err, "failed to listen on port %d", port)
	}

	server := grpc.NewServer()
	RegisterProofServer(server, s)
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) GetProof(ctx context.Context, request *ProofRequest) (*ProofResponse, error) {
	return &ProofResponse{
		Proof: fmt.Sprintf("hello proof! (clientID: %s)", request.ClientID),
	}, nil
}

var _ ProofServer = &Server{}
