package main

import "github.com/rai-project/config"

var AppSecret string

func init() {
	config.BeforeInit(func() {
		config.SetAppSecret(AppSecret)
	})
}
