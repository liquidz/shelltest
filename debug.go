package main

import (
	"fmt"
)

var ShellTestDebugMode = false

func DebugPrint(from string, format string, a ...interface{}) {
	if !ShellTestDebugMode {
		return
	}

	pre := fmt.Sprintf("DEBUG(%s)", from)
	fmt.Printf(fmt.Sprintf("\x1b[%dm%-12s:\x1b[0m %v\n", ColorYellow, pre, format), a...)
}
