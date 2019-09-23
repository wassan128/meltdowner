package main

import (
    "os"

    "github.com/wassan128/meltdowner/meltdowner/cmd"
)

func main() {
    cmd.RootCmd.SetOutput(os.Stdout)
    if err := cmd.RootCmd.Execute(); err != nil {
        cmd.RootCmd.SetOutput(os.Stderr)
        cmd.RootCmd.Println(err)
        os.Exit(1)
    }
}

