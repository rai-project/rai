package cmd

import (
	"os"
	"path/filepath"

	"github.com/rai-project/client"
	log "github.com/rai-project/logger"
)

func newClient(extraOpts ...client.Option) (*client.Client, error) {
	if wd, err := filepath.Abs(workingDir); err == nil {
		workingDir = sanitize(wd)
	}

	opts := []client.Option{
		client.Directory(workingDir),
		client.Stdout(os.Stdout),
		client.Stderr(os.Stderr),
		client.JobQueueName(jobQueueName),
	}
	if !isRatelimit {
		opts = append(opts, client.DisableRatelimit())
	}

	if outputDirectory {
		opts = append(opts, client.OutputDirectory(outputDirectory, forceOutput))
	}

	if buildFilePath != "" {
		absPath, err := filepath.Abs(buildFilePath)
		if err != nil {
			buildFilePath = absPath
		}
		opts = append(opts, client.BuildFilePath(absPath))
	}

	opts = extraClientOptions()

	opts = append(opts, extraOpts...)

	return client.New(opts...)
}

func runClient(client *client.Client) error {
	if err := client.Validate(); err != nil {
		return err
	}
	if err := client.Authenticate(); err != nil {
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
	if err := client.Wait(); err != nil {
		return err
	}
	if err := client.RecordJob(); err != nil {
		log.WithError(err).Error("job not recorded. If this was a submission, it was not recorded.")
		return err
	}
	return nil
}
