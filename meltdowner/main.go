package main

import (
	"github.com/wassan128/meltdowner/meltdowner/cmd"
	"github.com/wassan128/meltdowner/meltdowner/util"
)

func main() {
	err := cmd.RootCmd.Execute()
	util.ExitIfError(err)
}

