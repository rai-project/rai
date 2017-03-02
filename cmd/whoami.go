package cmd

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/user"
	"github.com/spf13/cobra"
)

var WhoamiCmd = &cobra.Command{
	Use:    "whoami",
	Short:  "Prints the user information.",
	Hidden: false,
	Run: func(cmd *cobra.Command, args []string) {
		profile, err := user.NewProfile("")
		if err != nil {
			return
		}
		pp.Println(*profile)
	},
}

func init() {
	RootCmd.AddCommand(WhoamiCmd)
}
