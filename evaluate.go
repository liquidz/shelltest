package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// c.f. http://www.commandlinefu.com/commands/view/3584/remove-color-codes-special-characters-with-sed
	AnsiColor = `\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`
)

type Callback func(TestCase, error)

type EvaludateError struct {
	n      int
	test   TestCase
	result string
}

func (e EvaludateError) Error() string {
	return strings.TrimSpace(fmt.Sprintf(`
%d) %v
   expected: %v
   actual  : %v
  `, e.n, e.test.Command, e.test.Expected, e.result))
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
				callback(testcase, nil)
				continue
			}

			DebugPrint("expected: %v, actual: %v", expected, result)
			if !expected.IsExpected(result) {
				err := EvaludateError{n: i, test: testcase, result: result}
				callback(testcase, err)
				errs = append(errs, err)
			} else {
				callback(testcase, nil)
			}
		case <-termch:
			if i < testLen-1 {
				err := errors.New(fmt.Sprintf("too few command result"))
				callback(suite.Tests[testLen-1], err)
				errs = append(errs, err)
			}
			break loop
		}
	}

	return errs
}
