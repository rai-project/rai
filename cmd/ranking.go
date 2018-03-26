package cmd

import (
	//"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/rai-project/config"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/model"
	"github.com/spf13/cobra"
	upper "upper.io/db.v3"
  "github.com/rai-project/auth/provider"
  "github.com/pkg/errors"
  "github.com/rai-project/client"
)

var numResults int

const (
	maxResults = 100
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// rankingCmd represents the ranking command
var rankingCmd = &cobra.Command{
	Use:   "ranking",
	Short: "View anonymous rankings.",
	Long:  `View anonymized convolution performance. Only the fastest result for each team is reported.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		db, err := mongodb.NewDatabase(config.App.Name)
		if err != nil {
			return err
		}
		defer db.Close()

		col, err := model.NewSp2018Ece408JobCollection(db)
		if err != nil {
			return err
		}
		defer col.Close()

		// Get submissions
		var jobs model.Sp2018Ece408Jobs
		// cond := upper.Or(
		// 	upper.Cond{
		// 		"model":       "ece408-high",
		// 		"correctness": 0.8562,
		// 	},
		// 	upper.Cond{
		// 		"model":       "ece408-low",
		// 		"correctness": 0.629,
		// 	},
		// )

		condInferencesExist := upper.Cond{"inferences.0 $exists": "true"}
		cond := upper.And(
			condInferencesExist,
			upper.Cond{
			  "is_submission": true,
			  "inferences.correctness": 0.8451,},
		)

		err = col.Find(cond, 0, 0, &jobs)
		if err != nil {
			return err
		}

		// keep only jobs with non-zero runtimes
		jobs = model.FilterNonZeroTimes(jobs)

		// Sort by fastest
		sort.Sort(model.ByMinOpRuntime(jobs))

		// Keep first instance of every team
		jobs = model.KeepFirstTeam(jobs) // Keep fastest entry for each team

		// show only numResults
		if numResults < 0 {
			numResults = maxResults
		}
		numResults = min(numResults, maxResults)
		numResults = min(numResults, len(jobs))
		//jobs = jobs[0:numResults]

		//for _, j := range jobs {
		//	fmt.Println(j)
		//}

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

    tname, err := client.ReturnTeamName(prof.Info().Username)

    // Create table of ranking
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"You","Rank", "Anonymized Team", "Fastest Conv (ms)"})

		var x int64
		var currentRank int64
		var currentMinOpRunTime float64
		x=1
		currentRank = 1
		currentMinOpRunTime = 0

		for _, j := range jobs {
		  if currentMinOpRunTime != float64(j.MinOpRuntime())/float64(time.Millisecond) {
        currentMinOpRunTime = float64(j.MinOpRuntime()) / float64(time.Millisecond)
        currentRank = x
      }

		  if tname == j.Teamname {
        table.Append([]string{tname + " -->", strconv.FormatInt(currentRank, 10), j.Anonymize().Teamname, strconv.FormatFloat(currentMinOpRunTime, 'f', 3, 64)})
      } else {
        table.Append([]string{"", strconv.FormatInt(currentRank, 10), j.Anonymize().Teamname, strconv.FormatFloat(currentMinOpRunTime, 'f', 3, 64)})
      }
      x++
		}

		table.Render()
		return nil
	},
}

func init() {
	rankingCmd.Flags().IntVarP(&numResults, "num-results", "n", 10, "Number of results to show (<"+strconv.Itoa(maxResults)+")")
	RootCmd.AddCommand(rankingCmd)
}
