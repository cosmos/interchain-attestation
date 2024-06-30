package server

import (
	"context"
	"fmt"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	UnimplementedProofServer

	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Serve(listenAddr string) error {
	s.logger.Debug("server.Serve", zap.String("listenAddr", listenAddr))

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", listenAddr)
	}

	server := grpc.NewServer()
	RegisterProofServer(server, s)
	s.logger.Info("server listening", zap.String("addr", lis.Addr().String()))
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
