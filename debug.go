package main

import (
	"fmt"
)

var ShellTestDebugMode = false

func DebugPrint(format string, a ...interface{}) {
	if !ShellTestDebugMode {
		return
	}

	fmt.Printf(fmt.Sprintf("DEBUG: %v\n", format), a...)
}
