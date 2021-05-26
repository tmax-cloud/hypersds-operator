package wrapper

import (
	"os"
)

// OsInterface is wrapper interface of os package
type OsInterface interface {
	// MkdirAll is wrapper function for MkdirAll function of os package
	MkdirAll(path string, perm os.FileMode) error
	// RemoveAll is wrapper function for RemoveAll function of os package
	RemoveAll(path string) error
}

// osStruct is struct to be used instead of os package
type osStruct struct {
}

// MkdirAll executes MkdirAll function of os package
func (o *osStruct) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// RemoveAll exctures RemoveAll function of os package
func (o *osStruct) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// OsWrapper is to be used instead of os package
var OsWrapper OsInterface

func init() {
	OsWrapper = &osStruct{}
}
