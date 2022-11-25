package config

import (
	"gopkg.in/yaml.v3"
	"github.com/backdround/deploy-configs/config/validate"

	"fmt"
	"reflect"
)

// fullConfigData used for parse given instances by user yaml
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
	yaml.Unmarshal(dataYaml, &fullConfig)

	// Gets config for given instance
	config, ok := fullConfig.Instances[instance]
	if !ok {
		availableInstances := reflect.ValueOf(fullConfig.Instances).MapKeys()
		err := fmt.Errorf("There is no instance %q in %v", instance, availableInstances)
		return nil, err
	}

	return &config, nil
}
