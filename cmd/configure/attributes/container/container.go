package container

import (
	"encoding/json"
	"fmt"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
	"github.com/pkg/errors"
)

const (
	Flag = "attribute-container"
)

// Attributes represents container attributes including image info and container name.
type Attributes struct {
	ImageInfo     `json:",inline"`
	ContainerName string `json:"k8s.container.name,omitempty"`
}

func (attr Attributes) ToMap() (map[string]string, error) {
	return structs.ToMap(attr)
}

// ParseAttributes parses a slice of raw attribute strings into a slice of Attributes.
func ParseAttributes(rawAttributes []string) ([]Attributes, error) {
	var attributeList = make([]Attributes, 0, len(rawAttributes))

	for _, attr := range rawAttributes {
		parsedAttr, err := parse(attr)
		if err != nil {
			return nil, err
		}

		attributeList = append(attributeList, *parsedAttr)
	}

	return attributeList, nil
}

// ToArgs converts a slice of Attributes to a list of args for a Pod Template.
func ToArgs(attributes []Attributes) ([]string, error) {
	var args = make([]string, 0, len(attributes))

	for _, attr := range attributes {
		jsonAttr, err := json.Marshal(attr)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		args = append(args, fmt.Sprintf("--%s=%s", Flag, string(jsonAttr)))
	}

	return args, nil
}

// parse unmarshals a raw attribute string into an Attributes struct.
func parse(rawAttribute string) (*Attributes, error) {
	var result Attributes

	err := json.Unmarshal([]byte(rawAttribute), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
