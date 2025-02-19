package container

import (
	"encoding/json"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
	"github.com/spf13/cobra"
)

const (
	flagKey = "attribute-container"
)

var (
	attributes []string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayVar(&attributes, flagKey, []string{}, "(Optional) Container-specific attributes in JSON format.")
}

type Attributes struct {
	ImageInfo     `json:",inline"`
	ContainerName string `json:"k8s.container.name"`
}

func (attr Attributes) ToMap() (map[string]string, error) {
	return structs.ToMap(attr)
}

func ParseAttributes() ([]Attributes, error) {
	return parseAttributes(attributes)
}

func parseAttributes(rawAttributes []string) ([]Attributes, error) {
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

func parse(rawAttribute string) (*Attributes, error) {
	var result Attributes

	err := json.Unmarshal([]byte(rawAttribute), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
