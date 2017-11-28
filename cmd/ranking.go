package cmd

import (
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

		col, err := model.NewFa2017Ece408JobCollection(db)
		if err != nil {
			return err
		}
		defer col.Close()

		// Get submissions
		var jobs model.Fa2017Ece408Jobs
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

		err = col.Find(condInferencesExist, 0, 0, &jobs)
		if err != nil {
			return err
		}

		// keep only jobs with correct inferences
		jobs = model.FilterCorrectInferences(jobs)

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
		jobs = jobs[0:numResults]

		// Anonymize
		for i, j := range jobs {
			jobs[i] = j.Anonymize()
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Anonymized Team", "Team's Fastest Conv (ms)"})
		for _, j := range jobs {
			table.Append([]string{j.Teamname, strconv.FormatFloat(float64(j.MinOpRuntime())/float64(time.Millisecond), 'f', -1, 64)})
		}

		table.Render()
		return nil
	},
}

func init() {
	rankingCmd.Flags().IntVarP(&numResults, "num-results", "n", 10, "Number of results to show (<"+strconv.Itoa(maxResults)+")")
	RootCmd.AddCommand(rankingCmd)
}
