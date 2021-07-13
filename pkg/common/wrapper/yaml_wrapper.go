package wrapper

import (
	"io"

	"gopkg.in/yaml.v2"
)

// YamlInterface is wrapper interface of yaml package
type YamlInterface interface {
	// Unmarshal is wrapper function for Unmarshal function of yaml package
	Unmarshal(in []byte, out interface{}) (err error)
	// Marshal is wrapper function for Marshal function of yaml package
	Marshal(in interface{}) (out []byte, err error)
	// NewDecoder is wrapper function for NewDecoder function of yaml package
	NewDecoder(r io.Reader) YamlDecoderInterface
}

// YamlDecoderInterface is wrapper interface of Decode struct of yaml package
type YamlDecoderInterface interface {
	// Decode is wrapper function for Decode function of yaml package
	Decode(v interface{}) (err error)
}

// yamlStruct is struct to be used instead of yaml package
type yamlStruct struct {
}

// Unmarshal executes Unmarshal function of yaml package

func (y *yamlStruct) Unmarshal(in []byte, out interface{}) (err error) {
	return yaml.Unmarshal(in, out)
}

// Marshal executes Marshal function of yaml package
func (y *yamlStruct) Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}

// NewDecoder executes NewDecoder function of yaml package
func (y *yamlStruct) NewDecoder(r io.Reader) YamlDecoderInterface {
	return yaml.NewDecoder(r)
}

// YamlWrapper is to be used instead of yaml package
var YamlWrapper YamlInterface

func init() {
	YamlWrapper = &yamlStruct{}
}
