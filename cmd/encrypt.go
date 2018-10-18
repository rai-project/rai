package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rai-project/utils"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:          "encrypt",
	Short:        "encrypts an input string using the system's secret",
	Hidden:       true,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.Errorf("expected one argument, got %v", len(args))
		}
		val, err := utils.EncryptStringBase64(appSecret, args[0])
		if err != nil {
			return err
		}
		fmt.Println(val)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(encryptCmd)
}
