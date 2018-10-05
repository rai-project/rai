// +build bench

package cmd

import (
	"sync"

	"github.com/k0kubun/pp"

	"github.com/Jeffail/tunny"
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
	benchCmd.PersistentFlags().IntVar(&concurrencyCount, "concurrency_count", 10, "Number of concurrent runs")
}

// benchmark the server
var benchCmd = &cobra.Command{
	Use:          "bench",
	Short:        "Bench",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		progress := pb.StartNew(iterationCount)
		defer progress.FinishPrint("finished benchmarking")

		var wg sync.WaitGroup

		runClient := func(arg interface{}) interface{} {
			defer wg.Done()
			wg.Add(1)
			go func() {
				defer wg.Done()
				runClient(client)
			}()
			return nil
		}

		execPool := tunny.NewFunc(concurrencyCount, runClient)
		defer execPool.Close()

		for ii := 0; ii < iterationCount; ii++ {
			wg.Add(1)
			go func() {
				execPool.Process(nil)
			}()
		}

		wg.Wait()

		pp.Println("done")

		return nil
	},
}
