package wrapper

import (
	"io/ioutil"
	"os"
)

// IoUtilInterface is wrapper interface of ioutil package
type IoUtilInterface interface {
	// ReadFile is wrapper function for Readfile function of ioutil package
	ReadFile(fileName string) ([]byte, error)
	// WriteFile is wrapper function for Writefile function of ioutil package
	WriteFile(fileName string, data []byte, fileMode os.FileMode) error
}

// ioUtilStruct is struct to be used instead of ioutil package
type ioUtilStruct struct {
}

// ReadFile executes ReadFile function of ioutil package
func (i *ioUtilStruct) ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

// WriteFile executes WriteFile function of ioutil package
func (i *ioUtilStruct) WriteFile(fileName string, data []byte, fileMode os.FileMode) error {
	return ioutil.WriteFile(fileName, data, fileMode)
}

// IoUtilWrapper is to be used instead of iotuil package
var IoUtilWrapper IoUtilInterface

func init() {
	IoUtilWrapper = &ioUtilStruct{}
}
