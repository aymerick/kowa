package commands

import (
	"github.com/aymerick/kowa/core"
	"github.com/aymerick/kowa/server"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Kowa",
	Long:  `Setup database and stuff`,
	Run:   setup,
}

func setup(cmd *cobra.Command, args []string) {
	core.EnsureUploadDir()

	// Insert oauth client
	oauthStorage := server.NewOAuthStorage()
	oauthStorage.SetupDefaultClient()
}
