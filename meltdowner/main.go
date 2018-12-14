package main

import (
	"fmt"
	"os"

	"github.com/wassan128/meltdowner/meltdowner/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

