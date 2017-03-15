package main

import "github.com/rai-project/config"

var AppSecret = "-secret-"

func init() {
	config.DefaultAppSecret = AppSecret
}
