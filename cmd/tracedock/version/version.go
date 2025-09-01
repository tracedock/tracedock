package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var BuildVersion string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "return the version for tracedock cli",
	Long:  `return the version for tracedock cli`,
	Args:  cobra.MinimumNArgs(0),
	Run:   execVersionCmd,
}

func execVersionCmd(cmd *cobra.Command, args []string) {
	fmt.Println(BuildVersion)
}
