package server

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tracedock/tracedock/internal/logger"
	"github.com/tracedock/tracedock/internal/server"

	trace "go.opentelemetry.io/proto/otlp/trace/v1"
)

var (
	paramGRPCPort   string
	paramHTTPPort   string
	paramConfigFile string
)

var ServerCmd = &cobra.Command{
	Use:   "server [flags]",
	Short: "Manages the tracedock server",
	Long:  `Manages the tracedock server`,
	Args:  cobra.MinimumNArgs(1),
}

var ServerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the Tracedock server",
	Long:  `Starts the Tracedock server to collect and process telemetry data`,
	Args:  cobra.NoArgs,
	Run:   execServerStartCmd,
}

func init() {
	ServerCmd.AddCommand(ServerStartCmd)

	ServerStartCmd.PersistentFlags().StringVarP(&paramGRPCPort, "grpc-port", "", "0.0.0.0:4317", "tcp port for gRPC server")
	ServerStartCmd.PersistentFlags().StringVarP(&paramHTTPPort, "http-port", "", "0.0.0.0:4318", "tcp port for HTTP server")
	ServerStartCmd.PersistentFlags().StringVarP(&paramConfigFile, "config", "c", "/etc/tracedock.yaml", "path to the configuration file")
}

func execServerStartCmd(cmd *cobra.Command, args []string) {
	supervisor := server.NewSupervisor()

	grpcServer := server.NewGRPCServer()
	httpServer := server.NewHTTPServer()

	supervisor.Add(paramGRPCPort, grpcServer)
	supervisor.Add(paramHTTPPort, httpServer)

	ingestor := func(resource *trace.ResourceSpans) error {
		logger.Error(fmt.Sprintf("received a resource with %d attributes", len(resource.Resource.Attributes)))

		return nil
	}

	grpcServer.RegisterTraceIngestor(ingestor)
	httpServer.RegisterTraceIngestor(ingestor)

	if err := supervisor.Run(); err != nil {
		logger.Error(fmt.Sprintf("error starting supervisor: %v", err))
		return
	}

	time.Sleep(10 * time.Millisecond)

	if err := supervisor.Wait(); err != nil {
		logger.Error(fmt.Sprintf("error waiting for supervisor: %v", err))
		return
	}
}
