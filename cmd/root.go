package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/fatih/color"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	_ "github.com/rai-project/logger/hooks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appSecret     string
	workingDir    string
	jobQueueName  string
	buildFilePath string
	isColor       bool
	isVerbose     bool
	isDebug       bool
	isRatelimit   bool
	submit        string
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
		client, err := newClient()
		if err != nil {
			return err
		}
		return runClient(client)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

var VersionCmd = cmd.VersionCmd

func init() {
	VersionCmd.Run = func(c *cobra.Command, args []string) {
		cmd.VersionCmd.Run(c, args)
		fmt.Println("ProjectMode: ", projectMode)
	}

	cobra.OnInitialize(initConfig, initColor)

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
	RootCmd.PersistentFlags().StringVarP(&cwd, "build", "f", "", "Path to the build file. Defaults to `cwd`/rai_build.yml file.")
	RootCmd.PersistentFlags().StringVarP(&jobQueueName, "queue", "q", "", "Name of the job queue. Infers queue from build file by default.")
	RootCmd.PersistentFlags().StringVarP(&appSecret, "secret", "s", "", "Pass in application secret.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", true, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")
	RootCmd.PersistentFlags().BoolVar(&isRatelimit, "ratelimit", true, "Toggle debug mode.")
	if projectMode {
		RootCmd.PersistentFlags().StringVar(&submit, "submit", "", "mark the kind of submission (m2, m3, final)")
	}

	RootCmd.MarkPersistentFlagRequired("path")

	// mark secret flag hidden
	RootCmd.PersistentFlags().MarkHidden("secret")
	RootCmd.PersistentFlags().MarkHidden("ratelimit")
	RootCmd.PersistentFlags().MarkHidden("queue")

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
