package color

import (
	"testing"
)

func TestDeleteColor(t *testing.T) {
	s := "foo"

	if s != DeleteColor(GreenStr(s)) {
		t.Errorf("ansi color should be deleted")
	}
}
