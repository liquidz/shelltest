package eval

import (
	"errors"
	"fmt"
	. "github.com/liquidz/shelltest/color"
	. "github.com/liquidz/shelltest/debug"
	. "github.com/liquidz/shelltest/testcase"
	"strings"
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

	var exports []string
	for k, v := range suite.EnvMap {
		exports = append(exports, fmt.Sprintf("export %s=%s", k, v))
	}

	inch, outch, termch, err := startShell(shell, exports...)
	if err != nil {
		return []error{err}
	}

	for _, tc := range suite.Tests {
		inch <- tc.Command
	}
	inch <- "exit"

	testLen := len(suite.Tests)
	DebugPrint("eval", "testLen: %d", testLen)
loop:
	for i, n := 0, 0; ; i++ {
		select {
		case result := <-outch:
			if i >= testLen {
				DebugPrint("eval", "continue  %d", i)
				continue
			}

			// remove color codes
			result = DeleteColor(result)

			testcase := suite.Tests[i]
			expected := testcase.Expected
			if len(expected) == 0 {
				DebugPrint("eval", "skip tests[%d]: %v", i, result)
				continue
			}

			DebugPrint("eval", "expected[%d]: %v, actual: [%v]", i, expected, result)
			if !expected.IsExpected(result) {
				err := EvaluateError{No: i, Test: testcase, Result: result}
				callback(n, testcase, err)
				errs = append(errs, err)
			} else {
				callback(n, testcase, nil)
			}
			n++
		case <-termch:
			if i < testLen-1 {
				err := errors.New(fmt.Sprintf("too few command result"))
				callback(n, suite.Tests[testLen-1], err)
				errs = append(errs, err)
			}
			break loop
		}
	}

	return errs
}
