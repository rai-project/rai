package cmd

import (
	"fmt"
	"os"

	ccmd "github.com/rai-project/cmd"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// self-updateCmd represents the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:     "self-update",
	Aliases: []string{"update"},
	Short:   "Update RAI if a new version exists",
	Long:    `This will allow RAI to update it's self.`,
	Run: func(cmd *cobra.Command, args []string) {
		updater := &selfupdate.Updater{
			CurrentVersion: ccmd.Version.Version,
			ApiURL:         "https://files.rai-project.com.s3.amazonaws.com/dist/",
			BinURL:         "https://files.rai-project.com.s3.amazonaws.com/dist/",
			DiffURL:        "https://files.rai-project.com.s3.amazonaws.com/dist/",
			Dir:            "rai/stable/",
			CmdName:        "rai",
			ForceCheck:     true,
		}

		fmt.Println("Running Self Update")
		err := updater.BackgroundRun()
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		fmt.Println("Self Update Finished")
	},
}

func init() {
	RootCmd.AddCommand(selfUpdateCmd)
}
