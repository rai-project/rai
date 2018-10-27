package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/rai-project/client"
	log "github.com/rai-project/logger"
	"github.com/xlab/closer"
)

func newClient(inputOpts ...client.Option) (*client.Client, error) {
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

	if outputDirectory != "" {
		opts = append(opts, client.OutputDirectory(outputDirectory, forceOutput))
	}

	if buildFilePath != "" {
		absPath, err := filepath.Abs(buildFilePath)
		if err != nil {
			buildFilePath = absPath
		}
		opts = append(opts, client.BuildFilePath(absPath))
	}

	opts = extraClientOptions(opts)

	opts = append(opts, inputOpts...)

	clnt, err := client.New(opts...)

	if err == nil {
		closer.Bind(func() {
			clnt.Disconnect()
		})
	}

	return clnt, err
}

func runClient(client *client.Client) error {

	if !com.IsDir(workingDir) {
		fmt.Printf("Error:: the directory specified = %s was not found. "+
			"Use the --path option to specify the directory you want to build.\n", workingDir)
		return errors.New("Invalid directory")
	}

	// validate the rai_build.yml file and user privileges
	if err := client.Validate(); err != nil {
		return err
	}
	// authenticate the user, but connecting it to the
	// various backend and creating session tokens
	if err := client.Authenticate(); err != nil {
		return err
	}
	// subscribe to the redis queue. the redis queue
	// is used to gather stdout/stderr from the server
	if err := client.Subscribe(); err != nil {
		return err
	}
	// upload the user directory to the storage server
	// the client first creates an archive stream and
	// uploads that stream to the storage server
	if err := client.Upload(); err != nil {
		return err
	}
	// publish the job to the queue server
	if err := client.Publish(); err != nil {
		return err
	}
	//
	if err := client.Connect(); err != nil {
		return err
	}
	// wait until we receive an end signal
	if err := client.Wait(); err != nil {
		return err
	}
	// we record the job into the database.
	// this is used to store information such as
	// ranking
	if err := client.RecordJob(); err != nil {
		log.WithError(err).Error("job not recorded. If this was a submission, it was not recorded.")
		return err
	}
	return nil
}
