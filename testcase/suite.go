package testcase

import (
	"errors"
	"fmt"
	"strings"
)

type TestSuite struct {
	Before TestCases
	After  TestCases
	Tests  TestCases
	EnvMap map[string]string
}

func (ts *TestSuite) Append(section string, test TestCase) error {
	if test.IsEmpty() {
		return nil
	}

	switch section {
	case "before":
		ts.Before = append(ts.Before, test)
		return nil
	case "after":
		ts.After = append(ts.After, test)
		return nil
	case "test":
		ts.Tests = append(ts.Tests, test)
		return nil
	}
	return errors.New(fmt.Sprintf("unknown section: %v", section))
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
