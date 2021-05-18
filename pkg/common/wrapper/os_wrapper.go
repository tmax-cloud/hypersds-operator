package wrapper

import (
	"os"
)

type OsInterface interface {
	MkdirAll(path string, perm os.FileMode) error
}

type osStruct struct {
}

func (o *osStruct) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

var OsWrapper OsInterface

func init() {
	OsWrapper = &osStruct{}
}
