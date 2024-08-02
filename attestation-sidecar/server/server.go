package server

import (
	"context"
	"github.com/gjermundgaraba/pessimistic-validation/attestationsidecar/attestors"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	UnimplementedClaimServer
	listener net.Listener

	logger      *zap.Logger
	coordinator attestors.Coordinator
	grpcServer  *grpc.Server
}

func NewServer(logger *zap.Logger, coordinator attestors.Coordinator) *Server {
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
	RegisterClaimServer(s.grpcServer, s)
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

func (s *Server) GetClaim(ctx context.Context, request *ClaimRequest) (*ClaimResponse, error) {
	s.logger.Debug("server.GetLatestSignedClaim", zap.String("chainId", request.ChainId))

	chainProver := s.coordinator.GetChainProver(request.ChainId)
	claim := chainProver.GetLatestSignedClaim()
	return &ClaimResponse{
		Claim: claim,
	}, nil
}

var _ ClaimServer = &Server{}
