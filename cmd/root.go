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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	workingDir string
	isColor    bool
	isVerbose  bool
	isDebug    bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "rai",
	Short:        "The client is used to submit jobs to the server.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if workingDir == "" || !com.IsDir(workingDir) {
			fmt.Printf("Error:: the directory specified = %s was not found.\n", workingDir)
			return errors.New("Invalid directory")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client.New(
			client.Directory(workingDir),
			client.Stdout(os.Stdout),
			client.Stderr(os.Stderr),
		)
		if err != nil {
			return err
		}
		if err := client.Validate(); err != nil {
			return err
		}
		if err := client.Init(); err != nil {
			return err
		}
		if err := client.Connect(); err != nil {
			return err
		}
		defer client.Disconnect()
		return nil
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initColor)

	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	RootCmd.PersistentFlags().StringVarP(&workingDir, "directory", "p", cwd,
		"Path to the directory you wish to submit. Defaults to the current working directory.")
	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", color.NoColor, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")

	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	config.Init(
		config.AppName("rai"),
		config.AppSecret(appsecret),
		config.ConfigString(_escFSMustString(false, "/.rai_config.yaml")),
	)
}

func initColor() {
	color.NoColor = !isColor
}
