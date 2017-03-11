package cmd

import (
	"strings"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/auth"
	"github.com/rai-project/auth/auth0"
	"github.com/rai-project/auth/secret"
	"github.com/spf13/cobra"
)

var WhoamiCmd = &cobra.Command{
	Use:          "whoami",
	Short:        "Prints the user information.",
	Hidden:       false,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error
		var prof auth.Profile

		provider := auth.Provider(strings.ToLower(auth.Config.Provider))
		switch provider {
		case auth.Auth0Provider:
			prof, err = auth0.NewProfile()
		case auth.SecretProvider:
			prof, err = secret.NewProfile()
		default:
			err = errors.Errorf("the auth provider %v specified is not supported", provider)
		}
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
