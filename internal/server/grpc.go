package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/tracedock/tracedock/internal/logger"
	"google.golang.org/grpc"

	tracecollectorv1 "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

// GRPCServer implements Server interface for the gRPC protocol
//
// Notice: It isn't implementing 100% of the OpenTelemetry HTTP specification
// regarding to the responses content.
//
// This behaviour wasn't tested yet.
//
// For more details: https://opentelemetry.io/docs/specs/otlp/#otlpgrpc-response
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
	var err error

	if s.traceIngestor == nil {
		return nil, ErrNoIngestorRegistered
	}

	for _, resource := range req.GetResourceSpans() {
		if thisErr := s.traceIngestor(resource); thisErr != nil {
			err = errors.Join(err, thisErr)
		}
	}

	return &tracecollectorv1.ExportTraceServiceResponse{}, err
}

// Start the gRPC server
func (s *GRPCServer) Start(addr string) error {
	logger.Info(fmt.Sprintf("starting gRPC server at %s", addr))

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
