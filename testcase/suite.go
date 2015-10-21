package testcase

import (
	"fmt"
	"strings"
)

type TestSuite struct {
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
[tests]
  %v
	`, ts.Tests)

	return strings.TrimSpace(s)
}

func mergeMap(m1, m2 map[string]string) map[string]string {
	res := map[string]string{}
	for k, v := range m1 {
		res[k] = v
	}
	for k, v := range m2 {
		res[k] = v
	}
	return res
}

func (ts1 TestSuite) Merge(ts2 TestSuite) TestSuite {
	res := TestSuite{EnvMap: mergeMap(ts1.EnvMap, ts2.EnvMap)}

	for _, tc := range ts1.Tests {
		res.Tests = append(res.Tests, tc)
	}
	for _, tc := range ts2.Tests {
		res.Tests = append(res.Tests, tc)
	}

	return res
}
