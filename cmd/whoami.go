package cmd

import (
	"github.com/spf13/cobra"
)

var WhoamiCmd = &cobra.Command{
	Use:          "whoami",
	Short:        "Prints the user information.",
	Hidden:       false,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return whoami()
	},
}

func init() {
	RootCmd.AddCommand(WhoamiCmd)
}
