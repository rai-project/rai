package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rai-project/auth/provider"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
		buf, err := yaml.Marshal(prof.Info())
		if err != nil {
			return err
		}
		fmt.Print(string(buf))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(WhoamiCmd)
}
