package server

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tracedock/tracedock/internal/config"
	"github.com/tracedock/tracedock/internal/logger"
	"github.com/tracedock/tracedock/internal/orchestrator"
	"github.com/tracedock/tracedock/internal/server"
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
	cfg := config.NewConfig()
	if err := cfg.Load(paramConfigFile); err != nil {
		logger.Error(fmt.Sprintf("error loading config file: %v", err))
		return
	}

	orchestrator := orchestrator.NewIngestor(cfg)

	supervisor := server.NewSupervisor()
	grpcServer := server.NewGRPCServer()
	httpServer := server.NewHTTPServer()

	supervisor.Add(paramGRPCPort, grpcServer)
	supervisor.Add(paramHTTPPort, httpServer)

	grpcServer.RegisterTraceIngestor(orchestrator.IngestTrace)
	httpServer.RegisterTraceIngestor(orchestrator.IngestTrace)

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
