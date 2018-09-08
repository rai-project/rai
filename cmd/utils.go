package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Unknwon/com"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/rai-project/auth/provider"
	"github.com/rai-project/config"
	"gopkg.in/yaml.v2"
)

var (
	DefaultProfilePath string
)

func checkWorkingDir() error {
	if workingDir == "" {
		fmt.Printf("Error:: the project directory cannot be empty. " +
			"Use the --path option to specify the directory you want to build.\n")
		return errors.New("Invalid empty directory")
	}
	if d, err := filepath.Abs(workingDir); err == nil {
		workingDir = d
	}
	if !com.IsDir(workingDir) {
		fmt.Printf("Error:: the directory specified = %s was not found. "+
			"Use the --path option to specify the directory you want to build.\n", workingDir)
		return errors.New("Invalid directory")
	}
	return nil
}

func whoami() error {
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
	buf, err := yaml.Marshal(prof.Info())
	if err != nil {
		return err
	}
	fmt.Print(string(buf))
	return nil
}

func init() {
	config.AfterInit(func() {
		homeDir, err := homedir.Dir()
		if err != nil {
			return
		}

		DefaultProfilePath = filepath.Join(homeDir, "."+config.App.Name+"_profile")
	})
}
