package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/aymerick/kowa/core"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display kowa version",
	Long:  "Outputs kowa version",
	Run:   displayVersion,
}

func displayVersion(cmd *cobra.Command, args []string) {
	fmt.Printf(core.FormatVersion())
}
