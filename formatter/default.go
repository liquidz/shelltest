package formatter

import (
	"fmt"
	. "github.com/liquidz/shelltest/color"
	. "github.com/liquidz/shelltest/debug"
	. "github.com/liquidz/shelltest/testcase"
)

type Formatter interface {
	Setup(TestSuite) string
	Result(int, TestCase, error) string
	TearDown(TestSuite, []error) string
}

type DefaultFormatter struct{}

func (f DefaultFormatter) Setup(suite TestSuite) string {
	return ""
}

func (f DefaultFormatter) Result(_ int, _ TestCase, err error) string {
	if err == nil {
		if !ShellTestDebugMode {
			return "."
		}
		return ""
	} else {
		s := fmt.Sprintf("\n%v\n\n", RedStr(err.Error()))
		if !ShellTestDebugMode {
			s += "x"
		}
		return s
	}
}

func (f DefaultFormatter) TearDown(suite TestSuite, errors []error) string {
	s := fmt.Sprintf("\n\n%v tests, %v failures\n", suite.Tests.Length(), len(errors))
	if len(errors) > 0 {
		return RedStr(s)
	}
	return GreenStr(s)
}
