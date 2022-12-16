// config Validates and parses user yaml config to Config structure.
package config

import (
	"github.com/backdround/deploy-configs/internal/config/validate"
	"gopkg.in/yaml.v3"

	"fmt"
	"reflect"
)

// fullConfigData represents all user instances parsed from user yaml
type fullConfigData struct {
	Instances map[string]Config `yaml:"instances"`
}

// Get validates, parses user yaml data and returns config for given instance.
func Get(dataYaml []byte, instance string) (*Config, error) {
	// Validates yaml config
	err := validate.Validate(dataYaml)
	if err != nil {
		return nil, err
	}

	// Parses yaml config
	fullConfig := fullConfigData{
		Instances: make(map[string]Config),
	}
	err = yaml.Unmarshal(dataYaml, &fullConfig)
	if err != nil {
		return nil, err
	}

	// Gets config for given instance
	config, ok := fullConfig.Instances[instance]
	if !ok {
		availableInstances := reflect.ValueOf(fullConfig.Instances).MapKeys()
		err := fmt.Errorf("There is no instance %q in %v", instance, availableInstances)
		return nil, err
	}

	return &config, nil
}
