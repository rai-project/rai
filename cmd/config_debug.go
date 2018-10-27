// +build debug

package cmd

import (
	"io/ioutil"
	"path/filepath"

	sourcepath "github.com/GeertJohan/go-sourcepath"
	"github.com/Unknwon/com"
)

var configContent string

func init() {
	configPath := filepath.Join(sourcepath.MustAbsoluteDir(), "..", "rai_config.yml")
	if !com.IsFile(configPath) {
		panic("unable to locate " + configPath)
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic("unable to read config file " + err.Error())
	}
	configContent = string(data)
}
