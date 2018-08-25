package cmd

import (
	"github.com/almostmoore/gbquestion/vars"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// RootCmd is an entrypoint to application
var RootCmd = &cobra.Command{
	Use:     "gbquestion",
	Short:   "Questions server and client tool",
	Version: vars.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return godotenv.Load()
	},
}

func init() {
	RootCmd.AddCommand(server)
	RootCmd.AddCommand(upsertCmd)
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(viewCmd)
}
