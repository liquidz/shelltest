package main

import (
	"io/ioutil"
)

var ReadFile = ioutil.ReadFile

func MockReadFile(content string, err error) {
	ReadFile = func(_ string) ([]byte, error) {
		return []byte(content), err
	}
}

func ResetMock() {
	ReadFile = ioutil.ReadFile
}
