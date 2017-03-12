package cmd

import (
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/auth/provider"
	"github.com/spf13/cobra"
)

var WhoamiCmd = &cobra.Command{
	Use:          "whoami",
	Short:        "Prints the user information.",
	Hidden:       false,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		prof, err := provider.New()
		if err != nil {
			return err
		}

		ok, err := prof.Verify()
		if err != nil {
			return err
		}
		if !ok {
			return errors.Errorf("cannot authenticate using the credentials in %v", prof.Options().ProfilePath)
		}
		pp.Println(prof.Info())
		return nil
	},
}

func init() {
	RootCmd.AddCommand(WhoamiCmd)
}
