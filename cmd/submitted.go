// +build ece408ProjectMode

package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rai-project/auth/provider"
	"github.com/rai-project/client"
	"github.com/rai-project/config"
	"github.com/rai-project/database/mongodb"
	"github.com/spf13/cobra"
	upper "upper.io/db.v3"
	//"gopkg.in/yaml.v2"
)

// rankingCmd represents the ranking command
var submittedCmd = &cobra.Command{}

func init() {
	if !ece408ProjectMode {
		return
	}
	submittedCmd = &cobra.Command{
		Use:   "submitted",
		Short: "View history of submissions.",
		Long:  `View history of team submissions associated with user`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the profile (e.g. ~/rai_profile.yml)
			prof, err := provider.New()
			if err != nil {
				return err
			}
			// Verify the profile
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

			// Create a database  using mongodb with the `config.App.Name` name
			db, err := mongodb.NewDatabase(config.App.Name)
			if err != nil {
				return err
			}
			defer db.Close()

			// Create the Fall2017 collection (mongodb's nomenclature for tables)
			col, err := client.NewFa2017Ece408TeamCollection(db)
			if err != nil {
				return err
			}

			var jobs client.Ece408JobResponseBodys

			condInferencesExist := upper.Cond{"inferences.0 $exists": "true"}
			cond := upper.And(
				condInferencesExist,
				upper.Cond{
					"is_submission": true,
					"teamname":      tname,
				},
			)

			// find all jobs which are both submissions and have the
			// team name equal to teamname. This would fill the
			// jobs list with the entries found within the collection
			err = col.Find(cond, 0, 0, &jobs)
			if err != nil {
				return err
			}

			if len(jobs) == 0 {
				print("No jobs associated with userid / teamname. If you have submitted a job, your userid is not associated with a teamname.")
				return nil
			}

			fmt.Println()
			fmt.Println("Last 5 successful submissions for team: " + tname)
			fmt.Println()

			// not sure what the heck this is doing
			// TODO: can use a slice
			x := 0
			for _, i := range jobs {
				//Skip items before last 5
				if x > len(jobs)-6 {
					fmt.Println(i.SubmissionTag + " - " + i.CreatedAt.String() + " (Submitted by: " + i.Username + ")")
				}
				x++
			}

			return nil
		},
	}
	RootCmd.AddCommand(submittedCmd)
}
