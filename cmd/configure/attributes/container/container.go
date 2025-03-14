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

type Attributes struct {
	ImageInfo     `json:",inline"`
	ContainerName string `json:"k8s.container.name"`
}

func (attr Attributes) ToMap() (map[string]string, error) {
	return structs.ToMap(attr)
}

func ParseAttributes(rawAttributes []string) ([]Attributes, error) {
	var attributeList []Attributes

	for _, attr := range rawAttributes {
		parsedAttr, err := parse(attr)
		if err != nil {
			return nil, err
		}

		attributeList = append(attributeList, *parsedAttr)
	}

	return attributeList, nil
}

// ToArgs is a helper func to convert an []container.Attributes to a list of args that can be put into a Pod Template
func ToArgs(attributes []Attributes) ([]string, error) {
	var args []string

	for _, attr := range attributes {
		jsonAttr, err := json.Marshal(attr)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		args = append(args, fmt.Sprintf("--%s=%s", Flag, string(jsonAttr)))
	}

	return args, nil
}

func parse(rawAttribute string) (*Attributes, error) {
	var result Attributes

	err := json.Unmarshal([]byte(rawAttribute), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
