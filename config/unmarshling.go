package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
)

// UnmarshalYAML used for Link parsing from yaml config by gopkg.in/yaml.v3
func (d *Link) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode {
		return errors.New("link data isn't a list element")
	}

	linkNodes := []yaml.Node{}
	node.Decode(&linkNodes)
	if len(linkNodes) != 2 {
		return fmt.Errorf("link expect two elements, but count of elements are %v",
			len(linkNodes))
	}

	err := linkNodes[0].Decode(&d.Target)
	if err != nil {
		return err
	}

	linkNodes[1].Decode(&d.LinkPath)
	if err != nil {
		return err
	}

	return nil
}
