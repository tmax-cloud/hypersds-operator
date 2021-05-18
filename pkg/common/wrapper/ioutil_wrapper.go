package wrapper

import (
	"io/ioutil"
	"os"
)

type IoUtilInterface interface {
	ReadFile(fileName string) ([]byte, error)
	WriteFile(fileName string, data []byte, fileMode os.FileMode) error
}

type ioUtilStruct struct {
}

func (i *ioUtilStruct) ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func (i *ioUtilStruct) WriteFile(fileName string, data []byte, fileMode os.FileMode) error {
	return ioutil.WriteFile(fileName, data, fileMode)
}

var IoUtilWrapper IoUtilInterface

func init() {
	IoUtilWrapper = &ioUtilStruct{}
}
