package main

import (
	"errors"
	"fmt"
	"strings"
)

type Callback func(TestCase, error)

func Evaluate(shell string, suite TestSuite, callback Callback) []error {
	var (
		errs []error
	)

	in, out, term, err := startShell(shell)
	if err != nil {
		return []error{err}
	}

	for _, tc := range suite.Tests {
		in <- tc.Command
	}
	in <- "exit"

	testLen := len(suite.Tests)
loop:
	for i := 0; ; i++ {
		select {
		case result := <-out:
			if i >= testLen {
				continue
			}

			testcase := suite.Tests[i]
			expected := testcase.Expected
			if len(expected) == 0 {
				callback(testcase, nil)
				continue
			}

			DebugPrint("expected: %v, actual: %v", expected, result)
			if !expected.IsExpected(result) {
				err := errors.New(strings.TrimSpace(fmt.Sprintf(`
%d) %v
  expected: %v
  actual  : %v
`, i, testcase.Command, expected, result)))
				callback(testcase, err)
				errs = append(errs, err)
			}

			callback(testcase, nil)
		case <-term:
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
