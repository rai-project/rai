package main

import "github.com/rai-project/config"

var AppSecret string

func init() {
	config.DefaultAppSecret = AppSecret
}
