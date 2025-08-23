package server

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tracecollectorv1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

// GRPCServer implements Server interface for the gRPC protocol
type GRPCServer struct {
	server        *grpc.Server
	traceIngestor TraceIngestor
	tracecollectorv1.UnimplementedTraceServiceServer
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer() *GRPCServer {
	grpcServer := &GRPCServer{
		server: grpc.NewServer(),
	}

	tracecollectorv1.RegisterTraceServiceServer(grpcServer.server, grpcServer)

	return grpcServer
}

// Export implements the interface UnimplementedTraceServiceServer that allows it
// to process incoming trace data
func (s *GRPCServer) Export(ctx context.Context, req *tracecollectorv1.ExportTraceServiceRequest) (*tracecollectorv1.ExportTraceServiceResponse, error) {
	if s.traceIngestor == nil {
		return nil, status.Error(codes.Internal, "no trace ingestor registered")
	}

	for _, resourceSpan := range req.GetResourceSpans() {
		if err := s.traceIngestor(resourceSpan); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to ingest trace: %v", err)
		}
	}

	return &tracecollectorv1.ExportTraceServiceResponse{}, nil
}

// Start the gRPC server
func (s *GRPCServer) Start(addr string) error {
	if s.traceIngestor == nil {
		return ErrNoIngestorRegistered
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.server.Serve(listener)
}

// Stop the gRPC server
func (s *GRPCServer) Stop() error {
	s.server.GracefulStop()
	return nil
}

// RegisterTraceIngestor registers a TraceIngestor function that will process all the
// incoming trace data
func (s *GRPCServer) RegisterTraceIngestor(ingestor TraceIngestor) {
	s.traceIngestor = ingestor
}
