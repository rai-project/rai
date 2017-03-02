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
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, err := user.NewProfile("")
		if err != nil {
			return err
		}
		pp.Println(*profile)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(WhoamiCmd)
}
