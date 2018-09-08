package cmd

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize the rai profile",
	Long: `Initialize (rai init) will create a new profile, with a license
and the appropriate structure for usage within rai. The profile is provided
by the rai system administrator usually through an email.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
