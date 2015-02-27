package commands

import (
	"github.com/aymerick/kowa/server"
	"github.com/aymerick/kowa/utils"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Kowa",
	Long:  `Setup database and stuff`,
	Run:   setup,
}

func setup(cmd *cobra.Command, args []string) {
	utils.AppEnsureUploadDir()

	// Insert oauth client
	oauthStorage := server.NewOAuthStorage()
	oauthStorage.SetupDefaultClient()
}
