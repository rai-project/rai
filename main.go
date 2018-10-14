package main

import (
	"github.com/rai-project/rai/cmd"
	"github.com/xlab/closer"
)

func cleanup() {
}

func main() {
	closer.Bind(cleanup)
	closer.Checked(cmd.Execute, false)
}
