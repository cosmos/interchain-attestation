package server

import (
	"context"
	"fmt"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"proversidecar/provers"
)

type Server struct {
	UnimplementedProofServer
	listener net.Listener

	logger *zap.Logger
	coordinator *provers.Coordinator
	grpcServer *grpc.Server
}

func NewServer(logger *zap.Logger, coordinator *provers.Coordinator) *Server {
	return &Server{
		logger: logger,
		coordinator: coordinator,
	}
}

func (s *Server) Serve(listenAddr string) error {
	s.logger.Debug("server.Serve", zap.String("listenAddr", listenAddr))

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", listenAddr)
	}

	s.grpcServer = grpc.NewServer()
	RegisterProofServer(s.grpcServer, s)
	s.logger.Info("server listening", zap.String("addr", lis.Addr().String()))
	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	s.logger.Debug("server.Stop")

	s.grpcServer.GracefulStop()
}

func (s *Server) GetProof(ctx context.Context, request *ProofRequest) (*ProofResponse, error) {
	s.logger.Debug("server.GetProof", zap.String("chainId", request.ChainId))

	chainProver := s.coordinator.GetChainProver(request.ChainId)
	proof := chainProver.GetProof()
	return &ProofResponse{
		Proof: fmt.Sprintf("proof: %s (clientID: %s)", string(proof), request.ChainId),
	}, nil
}

var _ ProofServer = &Server{}
