package main

import (
	"os"

	"meltdowner/src/interfaces"
)

func main() {
	interfaces.RootCmd.SetOutput(os.Stdout)
	if err := interfaces.RootCmd.Execute(); err != nil {
		interfaces.RootCmd.SetOutput(os.Stderr)
		interfaces.RootCmd.Println(err)
		os.Exit(1)
	}
}

