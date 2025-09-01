package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tracedock/tracedock/cmd/tracedock/server"
	"github.com/tracedock/tracedock/cmd/tracedock/version"
)

var rootCmd = &cobra.Command{
	Use:   "tracedock",
	Short: "Tracedock is an OpenTelemetry Collector",
	Long:  `Tracedock is an OpenTelemetry Collector that collects and processes telemetry data.`,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(server.ServerCmd)
	rootCmd.AddCommand(version.VersionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
