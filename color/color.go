package color

import (
	"fmt"
)

const (
	ColorRed    = 31
	ColorGreen  = 32
	ColorYellow = 33
)

var NoColor = false

func colorize(color int, s string) string {
	if NoColor {
		return s
	} else {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
	}
}

func RedStr(s string) string {
	return colorize(ColorRed, s)
}

func GreenStr(s string) string {
	return colorize(ColorGreen, s)
}

func YellowStr(s string) string {
	return colorize(ColorYellow, s)
}
