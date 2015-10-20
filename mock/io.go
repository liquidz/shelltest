package mock

import (
	"io/ioutil"
)

type ReadFileReturn struct {
	Content string
	Err     error
}

var ReadFile = ioutil.ReadFile

func MockReadFile(returns ...ReadFileReturn) {
	i := 0
	ReadFile = func(_ string) ([]byte, error) {
		content := []byte(returns[i].Content)
		err := returns[i].Err
		i++
		return content, err
	}
}

func ResetMock() {
	ReadFile = ioutil.ReadFile
}
