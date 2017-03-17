package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	ccmd "github.com/rai-project/cmd"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

type HTTPRequester struct {
}

// Fetch will return an HTTP request to the specified url and return
// the body of the result. An error will occur for a non 200 status code.
func (httpRequester *HTTPRequester) Fetch(url string) (io.ReadCloser, error) {
	if isDebug || isVerbose {
		fmt.Println("GET " + url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad http status from %s: %v", url, resp.Status)
	}

	return resp.Body, nil
}

// self-updateCmd represents the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:     "self-update",
	Aliases: []string{"update", "selfupdate"},
	Short:   "Update RAI if a new version exists",
	Long:    `This will allow RAI to update it's self.`,
	Run: func(cmd *cobra.Command, args []string) {
		updateDir, _ := homedir.Expand("~/.rai_update")
		updater := &selfupdate.Updater{
			CurrentVersion: ccmd.Version.Version,
			ApiURL:         "http://files.rai-project.com/dist/rai/stable/",
			BinURL:         "http://files.rai-project.com/dist/rai/stable/",
			CmdName:        "updates",
			Dir:            updateDir,
			Requester:      &HTTPRequester{},
			ForceCheck:     true,
		}

		fmt.Println("Running Self Update")
		err := updater.BackgroundRun()
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	},
}

func init() {
	RootCmd.AddCommand(selfUpdateCmd)
}
