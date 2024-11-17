package testutil

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/cosmos/interchain-attestation/core/types"
)

type Server struct {
	types.UnimplementedSidecarServer
	grpcServer *grpc.Server

	Response *types.GetIBCDataResponse
}

var _ types.SidecarServer = &Server{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Serve(listenAddr string) error {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.grpcServer = grpc.NewServer()
	types.RegisterSidecarServer(s.grpcServer, s)
	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

func (s *Server) GetIBCData(_ context.Context, _ *types.GetIBCDataRequest) (*types.GetIBCDataResponse, error) {
	return s.Response, nil
}
