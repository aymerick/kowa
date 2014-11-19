package commands

import (
	"github.com/aymerick/kowa/models"
	"github.com/aymerick/kowa/server"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Kowa",
	Long:  `WARNING: Removes all records from database !`,
	Run:   reset,
}

func reset(cmd *cobra.Command, args []string) {
	// @todo Check that we are NOT in production

	// reset models database
	models.DB().DropDatabase()

	// reset oauth database
	server.NewOAuthStorage().DB().DropDatabase()
}
