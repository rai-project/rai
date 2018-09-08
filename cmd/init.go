package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rai-project/auth/provider"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize the rai profile",
	Long: `Initialize (rai init) will create a new profile, with a license
and the appropriate structure for usage within rai. The profile is provided
by the rai system administrator usually through an email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if prof, err := provider.New(); err == nil {
			ok, err := prof.Verify()
			if err == nil && ok {
				return err
			}
		}

		scn := bufio.NewScanner(os.Stdin)
		var lines []string
		for scn.Scan() {
			line := scn.Text()
			if len(lines) == 0 && strings.TrimSpace(line) == "" {
				continue
			}
			if len(lines) > 1 && strings.TrimSpace(line) == "" {
				break
			}
			lines = append(lines, line)
		}
		profileContent := strings.Join(lines, "\n")
		if len(profileContent) > 0 {
			fmt.Println()
			fmt.Println("Profile:")
			fmt.Println(profileContent)
			fmt.Println()
		}

		if err := scn.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
