package debug

import (
	"fmt"
)

var color = 33 // yellow
var ShellTestDebugMode = false

func DebugPrint(from string, format string, a ...interface{}) {
	if !ShellTestDebugMode {
		return
	}

	pre := fmt.Sprintf("DEBUG(%s)", from)
	fmt.Printf(fmt.Sprintf("\x1b[%dm%-12s:\x1b[0m %v\n", color, pre, format), a...)
}
