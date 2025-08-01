package server

import (
	"fmt"

	"github.com/spf13/cobra"
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

	ServerStartCmd.PersistentFlags().StringVarP(&paramGRPCPort, "grpc-port", "", "4317", "tcp port for gRPC server")
	ServerStartCmd.PersistentFlags().StringVarP(&paramHTTPPort, "http-port", "", "4318", "tcp port for HTTP server")
	ServerStartCmd.PersistentFlags().StringVarP(&paramConfigFile, "config", "c", "/etc/tracedock.yaml", "path to the configuration file")
}

func execServerStartCmd(cmd *cobra.Command, args []string) {
	fmt.Println("starting application with params")
	fmt.Println()
	fmt.Println("gRPC port:\t", paramGRPCPort)
	fmt.Println("HTTP port:\t", paramHTTPPort)
	fmt.Println("Config file:\t", paramConfigFile)
}
