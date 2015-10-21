package color

import (
	"fmt"
	"regexp"
)

const (
	ColorRed    = 31
	ColorGreen  = 32
	ColorYellow = 33
)

var NoColor = false

// c.f. http://www.commandlinefu.com/commands/view/3584/remove-color-codes-special-characters-with-sed
var ansiColorRegexp = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`)

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

func DeleteColor(s string) string {
	return ansiColorRegexp.ReplaceAllString(s, "")
}
