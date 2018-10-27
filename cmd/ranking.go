// +build ece408ProjectMode

package cmd

import (
	//"fmt"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rai-project/auth/provider"
	"github.com/rai-project/client"
	"github.com/rai-project/config"
	"github.com/rai-project/database/mongodb"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	upper "upper.io/db.v3"
)

var numResults int

const (
	maxResults = 100
)

// rankingCmd represents the ranking command
var rankingCmd = &cobra.Command{}

func init() {
	if !ece408ProjectMode {
		return
	}
	// rankingCmd represents the ranking command
	rankingCmd = &cobra.Command{
		Use:   "ranking",
		Short: "View anonymous rankings.",
		Long:  `View anonymized convolution performance. Only the fastest result for each team is reported.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			min := func(a, b int) int {
				if a < b {
					return a
				}
				return b
			}

			db, err := mongodb.NewDatabase(config.App.Name)
			if err != nil {
				return err
			}
			defer db.Close()

			col, err := client.NewEce408JobResponseBodyCollection(db)
			if err != nil {
				return err
			}
			defer col.Close()

			// Get submissions

			condInferencesExist := upper.Cond{"inferences.0 $exists": "true"}
			cond := upper.And(
				condInferencesExist,
				upper.Cond{
					"is_submission":          true,
					"inferences.correctness": 0.8451},
			)

			var jobs client.Ece408JobResponseBodys
			err = col.Find(cond, 0, 0, &jobs)
			if err != nil {
				return err
			}

			// keep only jobs with non-zero runtimes
			jobs = client.FilterNonZeroTimes(jobs)

			// Sort by fastest
			sort.Sort(client.ByMinOpRuntime(jobs))

			// Keep first instance of every team
			jobs = client.KeepFirstTeam(jobs) // Keep fastest entry for each team

			// show only numResults
			if numResults < 0 {
				numResults = maxResults
			}

			numResults = min(numResults, len(jobs))

			// Get current user details
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

			tname, err := client.FindTeamName(prof.Info().Username)
			if err != nil {
				return err
			}

			if tname == "" {
				return errors.Errorf("No team name for %v", prof.Info().Username)
			}

			// Create table of ranking
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"You", "Rank", "Anonymized Team", "Fastest (ms)"})

			currentRank := 1
			currentMinOpRunTime := time.Duration(0)

			for ii, job := range jobs {
				if currentMinOpRunTime != job.MinOpRuntime() {
					currentMinOpRunTime = job.MinOpRuntime()
					currentRank = ii
				}

				srank := cast.ToString(currentRank)
				sMinOpTime := fmt.Sprintf("%v", currentMinOpRunTime)
				anonymizedTeamName := job.Anonymize().Teamname

				row := []string{tname + " -->", srank, anonymizedTeamName, sMinOpTime}

				if tname == job.Teamname {
					row[0] = ""
				}
				table.Append(row)
			}
			table.Render()
			return nil
		},
	}
	rankingCmd.Flags().IntVarP(&numResults, "num-results", "n", 10, "Number of results to show (<"+strconv.Itoa(maxResults)+")")
	RootCmd.AddCommand(rankingCmd)
}
