package server

import (
	"context"
	"net"

	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/cosmos/interchain-attestation/core/types"
	"github.com/cosmos/interchain-attestation/sidecar/attestators"
)

type Server struct {
	types.UnimplementedSidecarServer

	logger      *zap.Logger
	coordinator attestators.Coordinator
	grpcServer  *grpc.Server
}

var _ types.SidecarServer = &Server{}

func NewServer(logger *zap.Logger, coordinator attestators.Coordinator) *Server {
	return &Server{
		logger:      logger,
		coordinator: coordinator,
	}
}

func (s *Server) Serve(listenAddr string) error {
	s.logger.Debug("server.Serve", zap.String("listenAddr", listenAddr))

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return errors.Errorf("failed to listen on %s: %w", listenAddr, err)
	}

	s.grpcServer = grpc.NewServer()
	types.RegisterSidecarServer(s.grpcServer, s)
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

func (s *Server) GetAttestations(_ context.Context, _ *types.GetAttestationsRequest) (*types.GetAttestationsResponse, error) {
	s.logger.Debug("server.GetLatestAttestation")

	attestations, err := s.coordinator.GetLatestAttestations()
	if err != nil {
		return nil, err
	}

	return &types.GetAttestationsResponse{
		Attestations: attestations,
	}, nil
}
