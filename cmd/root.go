package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/rai-project/client"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	_ "github.com/rai-project/logger/hooks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appSecret    string
	workingDir   string
	jobQueueName string
	isColor      bool
	isVerbose    bool
	isDebug      bool
	isRatelimit  bool
	isSubmission bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "rai",
	Short:        "The client is used to submit jobs to the server.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if workingDir == "" || !com.IsDir(workingDir) {
			fmt.Printf("Error:: the directory specified = %s was not found. "+
				"Use the --path option to specify the directory you want to build.\n", workingDir)
			return errors.New("Invalid directory")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := []client.Option{
			client.Directory(workingDir),
			client.Stdout(os.Stdout),
			client.Stderr(os.Stderr),
			client.JobQueueName(jobQueueName),
		}
		if !isRatelimit {
			opts = append(opts, client.DisableRatelimit())
		}
		if isSubmission {
			opts = append(opts, client.IsSubmission(true))
		}

		client, err := client.New(opts...)

		if err != nil {
			return err
		}
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
		if err := client.RecordIfSubmission(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initColor)

	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)
	RootCmd.AddCommand(cmd.GendocCmd)
	RootCmd.AddCommand(cmd.CompletionCmd)
	RootCmd.AddCommand(cmd.BuildTimeCmd)

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	RootCmd.PersistentFlags().StringVarP(&workingDir, "path", "p", cwd,
		"Path to the directory you wish to submit. Defaults to the current working directory.")
	RootCmd.PersistentFlags().StringVarP(&jobQueueName, "queue", "q", "rai", "Name of the job queue. Defaults to 'rai'")
	RootCmd.PersistentFlags().StringVarP(&appSecret, "secret", "s", "", "Pass in application secret.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", true, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")
	RootCmd.PersistentFlags().BoolVar(&isRatelimit, "ratelimit", true, "Toggle debug mode.")
	RootCmd.PersistentFlags().BoolVar(&isSubmission, "submit", false, "mark as a final submission")

	RootCmd.MarkPersistentFlagRequired("path")

	// mark secret flag hidden
	RootCmd.PersistentFlags().MarkHidden("secret")
	RootCmd.PersistentFlags().MarkHidden("ratelimit")

	// viper.BindPFlag("app.secret", RootCmd.PersistentFlags().Lookup("secret"))
	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("app.color", RootCmd.PersistentFlags().Lookup("color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	opts := []config.Option{
		config.AppName("rai"),
		config.ColorMode(isColor),
		config.ConfigString(configContent),
	}
	if appSecret != "" {
		opts = append(opts, config.AppSecret(appSecret))
	}
	config.Init(opts...)
}

func initColor() {
	color.NoColor = !isColor
}
