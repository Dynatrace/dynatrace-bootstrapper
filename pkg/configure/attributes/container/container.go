package container

import (
	"encoding/json"
	"strings"

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

func ParseAttributes() (*Attributes, error) {
	return parseAttributes(attributes)
}

func parseAttributes(rawAttributes []string) (*Attributes, error) {
	logrus.Infof("Starting to parse container attributes for: %s", rawAttributes)

	rawMap := make(map[string]string)

	for _, attr := range rawAttributes {
		parts := strings.Split(attr, "=")
		if len(parts) == 2 {
			rawMap[parts[0]] = parts[1]
		}
	}

	raw, err := json.Marshal(rawMap)
	if err != nil {
		return nil, err
	}

	var result Attributes

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
