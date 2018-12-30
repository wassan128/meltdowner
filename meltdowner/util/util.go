package util

import (
	"fmt"
	"os"
)

func ExitIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("[Error]: %s", err))
	os.Exit(1)
}

func ExitIfFalse(status bool) {
	if status == false {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", "[Error] Validation failed")
	os.Exit(1)
}

func WarningIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[32;1m%s\x1b[0m\n", fmt.Sprintf("[Warning]: %s", err))
}

func Info(msg string) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf("[Info] %s", msg))
}

