package container

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
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

type ImageInfo struct {
	Registry    string `json:"container_image.registry"`
	Repository  string `json:"container_image.repository"`
	Tag         string `json:"container_image.tags"`
	ImageDigest string `json:"container_image.digest"`
}

type Attributes struct {
	ImageInfo     `json:",inline"`
	ContainerName string `json:"k8s.container.name"`
}

func ParseAttributes() ([]Attributes, error) {
	var attributeList []Attributes
	for _, attr := range attributes {
		parsedAttr, err := parseAttributes(attr)
		if err != nil {
			return nil, err
		}

		attributeList = append(attributeList, *parsedAttr)
	}
	return attributeList, nil
}

func parseAttributes(rawAttribute string) (*Attributes, error) {
	logrus.Infof("Starting to parse container attributes for: %s", rawAttribute)

	var result Attributes

	err := json.Unmarshal([]byte(rawAttribute), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
