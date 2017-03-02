package cmd

import "github.com/rai-project/config"

var (
	appsecret = "-secret-"
)

func init() {
	config.DefaultAppSecret = appsecret
}
