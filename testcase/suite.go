package testcase

import (
	"fmt"
	"strings"
)

type TestSuite struct {
	Before TestCases
	After  TestCases
	Tests  TestCases
	EnvMap map[string]string
}

func (ts *TestSuite) Append(test TestCase) error {
	if test.IsEmpty() {
		return nil
	}

	ts.Tests = append(ts.Tests, test)
	return nil
}

func (ts *TestSuite) String() string {
	s := fmt.Sprintf(`
[before]
  %v
[tests]
  %v
[after]
  %v
	`, ts.Before, ts.Tests, ts.After)

	return strings.TrimSpace(s)
}
