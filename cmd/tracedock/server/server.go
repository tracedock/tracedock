package server

import (
	"fmt"

	"github.com/spf13/cobra"
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
	Short: "manages the tracedock server",
	Long:  `manages the tracedock server`,
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
	orchestrator := server.NewOrchestrator()

	grpcServer := server.NewGRPCServer()
	httpServer := server.NewHTTPServer()

	orchestrator.Add(paramGRPCPort, grpcServer)
	orchestrator.Add(paramHTTPPort, httpServer)

	ingestor := func([]*trace.ResourceSpans) error {
		return nil
	}

	grpcServer.RegisterTraceIngestor(ingestor)
	httpServer.RegisterTraceIngestor(ingestor)

	if err := orchestrator.Run(); err != nil {
		fmt.Printf("Error starting orchestrator: %v\n", err)
		return
	}

	if err := orchestrator.Wait(); err != nil {
		fmt.Printf("Error waiting for orchestrator: %v\n", err)
		return
	}
}
