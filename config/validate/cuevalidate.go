package validate

import (
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"
	_ "embed"
)

//go:embed config.cue
var cueSchema []byte

// Validate validates dataYaml by cueSchema, returns error if
// validation is failed
func Validate(dataYaml []byte) error {
	cueContext := cuecontext.New()
	cueSchema := cueContext.CompileBytes(cueSchema)
	err := yaml.Validate(dataYaml, cueSchema)
	return err
}
