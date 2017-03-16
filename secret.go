package main

import "github.com/rai-project/rai/cmd"

var AppSecret string

func init() {
	cmd.AppSecret = AppSecret
}
