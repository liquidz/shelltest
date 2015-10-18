package debug

import (
	"fmt"
	. "github.com/liquidz/shelltest/color"
)

var ShellTestDebugMode = false

func DebugPrint(from string, format string, a ...interface{}) {
	if !ShellTestDebugMode {
		return
	}

	pre := fmt.Sprintf("%-12s:", fmt.Sprintf("DEBUG(%s)", from))
	fmt.Printf(fmt.Sprintf("%s %v\n", YellowStr(pre), format), a...)
}
