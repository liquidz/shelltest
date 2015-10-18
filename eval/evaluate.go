package eval

import (
	"errors"
	"fmt"
	. "github.com/liquidz/shelltest/debug"
	. "github.com/liquidz/shelltest/testcase"
	"regexp"
	"strings"
)

const (
	// c.f. http://www.commandlinefu.com/commands/view/3584/remove-color-codes-special-characters-with-sed
	AnsiColor = `\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`
)

type Callback func(int, TestCase, error)

type EvaluateError struct {
	No     int
	Test   TestCase
	Result string
}

func (e EvaluateError) Error() string {
	comment := e.Test.Comment
	if comment == "" {
		comment = e.Test.Command
	}
	return strings.TrimSpace(fmt.Sprintf(`
%d) %v
   command : %v
   expected: %v
   actual  : %v
  `, e.No, comment, e.Test.Command, e.Test.Expected, e.Result))
}

func Evaluate(shell string, suite TestSuite, callback Callback) []error {
	var errs []error

	ansiColorRegexp := regexp.MustCompile(AnsiColor)

	inch, outch, termch, err := startShell(shell)
	if err != nil {
		return []error{err}
	}

	for _, tc := range suite.Tests {
		inch <- tc.Command
	}
	inch <- "exit"

	testLen := len(suite.Tests)
loop:
	for i := 0; ; i++ {
		select {
		case result := <-outch:
			if i >= testLen {
				continue
			}

			// remove color codes
			result = ansiColorRegexp.ReplaceAllString(result, "")

			testcase := suite.Tests[i]
			expected := testcase.Expected
			if len(expected) == 0 {
				DebugPrint("eval", "skip tests[%d]: %v", i, result)
				callback(i, testcase, nil)
				continue
			}

			DebugPrint("eval", "expected[%d]: %v, actual: [%v]", i, expected, result)
			if !expected.IsExpected(result) {
				err := EvaluateError{No: i, Test: testcase, Result: result}
				callback(i, testcase, err)
				errs = append(errs, err)
			} else {
				callback(i, testcase, nil)
			}
		case <-termch:
			if i < testLen-1 {
				err := errors.New(fmt.Sprintf("too few command result"))
				callback(i, suite.Tests[testLen-1], err)
				errs = append(errs, err)
			}
			break loop
		}
	}

	return errs
}
