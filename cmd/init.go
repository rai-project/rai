package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/AlecAivazis/survey"
	"github.com/Unknwon/com"
	"github.com/pkg/errors"
	"github.com/rai-project/auth"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize the rai profile",
	Long: `Initialize (rai init) will create a new profile, with a license
and the appropriate structure for usage within rai. The profile is provided
by the rai system administrator usually through an email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if DefaultProfilePath == "" {
			return errors.New("Unable to infer the profile path. " +
				"You may need to initialize the profile manually.")
		}

		if com.IsFile(DefaultProfilePath) {
			confirm := false
			prompt := &survey.Confirm{
				Message: "Do you want to override you existing rai profile?",
			}
			survey.AskOne(prompt, &confirm, nil)
			if !confirm {
				return errors.New("The rai profile already exist.")
			}
		}

		fmt.Println("Paste the profile content that was included in the email bellow " +
			"(insert and empty line [i.e. press enter] to signify the end of the profile): ")

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

		if err := scn.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}

		profileContent := strings.Join(lines, "\n")

		prof := auth.ProfileBase{}
		if err := yaml.Unmarshal([]byte(profileContent), &prof); err != nil {
			return errors.Wrapf(err, "Invalid profile input. Make sure you have "+
				"correctly copied the profile content from the email")
		}

		err := ioutil.WriteFile(DefaultProfilePath, []byte(profileContent), 0644)
		if err != nil {
			return err
		}

		fmt.Println("Profile written. Checking if everything is ok.")

		return whoami()
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
