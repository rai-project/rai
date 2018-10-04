package cmd

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/Jeffail/tunny"
	"github.com/rai-project/client"
	log "github.com/rai-project/logger"
	_ "github.com/rai-project/logger/hooks"
	"github.com/spf13/cobra"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	iterationCount   int
	concurrencyCount int
)

func init() {
	RootCmd.AddCommand(benchCmd)
	benchCmd.PersistentFlags().IntVar(&iterationCount, "iteration_count", 100, "Number of iterations.")
	benchCmd.PersistentFlags().IntVar(&concurrencyCount, "concrrency_count", 10, "Number of concurrent runs")
}

// benchmark the server
var benchCmd = &cobra.Command{
	Use:          "bench",
	Short:        "Bench",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := []client.Option{
			client.Directory(workingDir),
			client.Stdout(os.Stdout),
			client.Stderr(os.Stderr),
			client.JobQueueName(jobQueueName),
			client.DisableRatelimit(),
			client.SubmissionCustom("bench"),
		}

		if buildFilePath != "" {
			absPath, err := filepath.Abs(buildFilePath)
			if err != nil {
				buildFilePath = absPath
			}
			opts = append(opts, client.BuildFilePath(absPath))
		}

		client, err := client.New(opts...)
		if err != nil {
			return err
		}

		progress := pb.StartNew(iterationCount)
		defer progress.FinishPrint("finished benchmarking")

		var wg sync.WaitGroup
		runClient := func() error {
			defer wg.Done()
			for ii := 0; ii < iterationCount; ii++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if err := client.Validate(); err != nil {
						return err
					}
					if err := client.Subscribe(); err != nil {
						return err
					}
					if err := client.Upload(); err != nil {
						return err
					}
					if err := client.Publish(); err != nil {
						return err
					}
					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Disconnect()
					if err := client.Wait(); err != nil {
						return err
					}
					if err := client.RecordJob(); err != nil {
						log.WithError(err).Error("job not recorded. If this was a submission, it was not recorded.")
						return err
					}
				}()
			}
			return nil
		}

		execPool := tunny.NewFunc(concurrentCount, runClient)
		defer execPool.Close()
		execPool.Process()

		return nil
	},
}
