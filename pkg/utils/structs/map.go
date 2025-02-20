package structs

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ToMap converts *SIMPLE* structs into a map[string]string.
// Obviously it only works for structs that can be represented as a map[string]string.
// So it will not work with more complicated structs.
func ToMap[T any](input T) (map[string]string, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	contentMap := map[string]string{}

	err = json.Unmarshal(raw, &contentMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return contentMap, nil
}

// FromMap converts *SIMPLE* map[string]string into a struct.
// Obviously it only works for structs that can be represented as a map[string]string.
// So it will not work with more complicated structs.
func FromMap[T any](input map[string]string) (*T, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var content T

	err = json.Unmarshal(raw, &content)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &content, nil
}
