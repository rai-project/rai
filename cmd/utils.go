package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Unknwon/com"
	"github.com/pkg/errors"
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
