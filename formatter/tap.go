package formatter

import (
	"fmt"
	. "github.com/liquidz/shelltest/color"
	//. "github.com/liquidz/shelltest/eval"
	. "github.com/liquidz/shelltest/testcase"
)

type TapFormatter struct{}

func (f TapFormatter) Setup(suite TestSuite) string {
	l := len(suite.Tests)
	if l > 0 {
		return fmt.Sprintf("1..%d\n", len(suite.Tests))
	} else {
		return "0..0\n"
	}
}

func (f TapFormatter) Result(no int, tc TestCase, err error) string {
	comment := tc.Comment
	if comment == "" {
		comment = tc.Command
	}

	if err == nil {
		return fmt.Sprintf("%s %d - %s\n", GreenStr("ok"), no+1, comment)
	} else {
		return fmt.Sprintf("%s %d - %s\n", RedStr("not ok"), no+1, comment)
	}
}

func (f TapFormatter) TearDown(suite TestSuite, errors []error) string {
	return ""
}
