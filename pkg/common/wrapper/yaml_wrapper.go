package wrapper

import (
	"io"

	"gopkg.in/yaml.v2"
)

type YamlInterface interface {
	Unmarshal(in []byte, out interface{}) (err error)
	Marshal(in interface{}) (out []byte, err error)
	NewDecoder(r io.Reader) YamlDecoderInterface
}
type YamlDecoderInterface interface {
	Decode(v interface{}) (err error)
}
type yamlStruct struct {
}

func (y *yamlStruct) Unmarshal(in []byte, out interface{}) (err error) {
	return yaml.Unmarshal(in, out)
}
func (y *yamlStruct) Marshal(in interface{}) (out []byte, err error) {
	return yaml.Marshal(in)
}

func (y *yamlStruct) NewDecoder(r io.Reader) YamlDecoderInterface {
	return yaml.NewDecoder(r)
}

var YamlWrapper YamlInterface

func init() {
	YamlWrapper = &yamlStruct{}
}
