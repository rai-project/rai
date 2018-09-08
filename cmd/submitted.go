package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rai-project/auth/provider"
	"github.com/rai-project/client"
	"github.com/rai-project/config"
	"github.com/rai-project/database/mongodb"
	"github.com/rai-project/model"
	"github.com/spf13/cobra"
	upper "upper.io/db.v3"
	//"gopkg.in/yaml.v2"
)

func init() {
	if !projectMode {
		return
	}
	submittedCmd = &cobra.Command{
		Use:   "submitted",
		Short: "View history of submissions.",
		Long:  `View history of team submissions associated with user`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := checkWorkingDir(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

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

			var jobs model.Sp2018Ece408Jobs

			condInferencesExist := upper.Cond{"inferences.0 $exists": "true"}
			cond := upper.And(
				condInferencesExist,
				upper.Cond{
					"is_submission": true,
					"teamname":      tname,
				},
			)

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
