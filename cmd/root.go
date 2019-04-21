package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	_ "github.com/rai-project/logger/hooks" // include all logging hooks
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xlab/catcher"
)

var (
	appSecret       string
	workingDir      string
	jobQueueName    string
	buildFilePath   string
	isColor         bool
	isVerbose       bool
	isDebug         bool
	isRatelimit     bool
	submitionName   string
	outputDirectory string
	forceOutput     bool
	serverArch      string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "rai",
	Short:        "The client is used to submit jobs to the server.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if workingDir == "" {
			cwd, err := os.Getwd()
			if err == nil {
				cwd, err = filepath.Abs(cwd)
				if err == nil {
					workingDir = cwd
				}
			}
		}
		if jobQueueName == "" && ece408ProjectMode {
			jobQueueName = "rai_amd64_ece408"
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// create a new rai client
		client, err := newClient()
		if err != nil {
			return err
		}
		// destroy the client before exiting the function
		defer client.Disconnect()
		// run the client steps
		return runClient(client)
	},
}

func safeCall() (err error) {
	// First one will put the panic message into the error value;
	// second one will yield the message to the stderr without the stracktrace.
	defer catcher.Catch(
		catcher.RecvError(&err, isDebug),
		catcher.RecvDie(1, true),
		catcher.RecvWrite(os.Stderr, isVerbose),
	)

	// run the main executable
	err = RootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

func Execute() error {
	// make sure to capture panics
	defer catcher.Catch(
		catcher.RecvWrite(os.Stderr, isVerbose),
	)
	return safeCall()
}

var VersionCmd = cmd.VersionCmd

func init() {
	// The version information command
	VersionCmd.Run = func(c *cobra.Command, args []string) {
		cmd.VersionCmd.Run(c, args)
		if ece408ProjectMode {
			fmt.Println("ECE408ProjectMode: ", ece408ProjectMode)
		}
	}

	cobra.OnInitialize(initConfig, initColor)

	// add the commands
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)
	RootCmd.AddCommand(cmd.GendocCmd)
	RootCmd.AddCommand(cmd.CompletionCmd)
	RootCmd.AddCommand(cmd.BuildTimeCmd)

	cwd, err := os.Getwd()
	if err == nil {
		cwd, err = filepath.Abs(cwd)
	}
	if err != nil {
		cwd = ""
	}

	RootCmd.PersistentFlags().StringVarP(&workingDir, "path", "p", cwd,
		"Path to the directory you wish to submit. Defaults to the current working directory.")
	RootCmd.PersistentFlags().StringVar(&cwd, "build", "", "Path to the build file. Defaults to `cwd`/rai_build.yml file.")
	RootCmd.PersistentFlags().StringVarP(&jobQueueName, "queue", "q", "", "Name of the job queue. Infers queue from build file by default.")
	RootCmd.PersistentFlags().StringVarP(&appSecret, "secret", "s", "", "Pass in application secret.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", true, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")
	RootCmd.PersistentFlags().StringVarP(&outputDirectory, "output", "o", "", "Set to output directory.")
	RootCmd.PersistentFlags().BoolVar(&forceOutput, "force", false, "Toggle to force overwriting output directory.")
	RootCmd.PersistentFlags().BoolVar(&isRatelimit, "ratelimit", true, "Toggle rate limiter.")
	RootCmd.PersistentFlags().StringVarP(&serverArch, "arch", "", "ppc64le", "RAID server architecture.")
	if ece408ProjectMode {
		RootCmd.PersistentFlags().StringVar(&submitionName, "submit", "", "The kind of submission (m2, m3, final)")
	}

	RootCmd.MarkPersistentFlagRequired("path")

	// mark secret flag hidden
	RootCmd.PersistentFlags().MarkHidden("secret")
	RootCmd.PersistentFlags().MarkHidden("ratelimit")
	RootCmd.PersistentFlags().MarkHidden("queue")

	// bind the flags specified to the configuration file
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
