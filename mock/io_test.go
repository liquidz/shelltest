package mock

import (
	"errors"
	"testing"
)

func TestMockReadFile(t *testing.T) {
	var (
		b   []byte
		err error
	)
	MockReadFile(
		ReadFileReturn{"hello", nil},
		ReadFileReturn{"", errors.New("error!")},
	)
	defer ResetMock()

	// first time
	b, err = ReadFile("foobar")
	if err != nil {
		t.Errorf("error occured: %v", err)
	}
	if string(b) != "hello" {
		t.Errorf("expected response is 'hello' but '%v'", string(b))
	}

	// second time
	b, err = ReadFile("barbaz")
	if err == nil {
		t.Errorf("error should be occured but nil")
	}
	if string(b) != "" {
		t.Errorf("expected response is empty but '%v'", string(b))
	}
}
