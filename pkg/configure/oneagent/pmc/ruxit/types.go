package ruxit

import (
	"sort"
	"strings"
)

// ProcMap presents the content in a more easy to work with format. (a map of maps).
type ProcMap map[string]map[string]string

// ProcConf represents the response of the /deployment/installer/agent/processmoduleconfig endpoint from the Dynatrace Environment(v1) API.
type ProcConf struct {
	Properties []Property `json:"properties"`
	Revision   uint       `json:"revision"`
}

type Property struct {
	Section string `json:"section"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

func (pmc ProcConf) ToMap() ProcMap {
	sections := map[string]map[string]string{}
	for _, prop := range pmc.Properties {
		section := sections[prop.Section]
		if section == nil {
			section = map[string]string{}
		}

		section[prop.Key] = prop.Value
		sections[prop.Section] = section
	}

	return sections
}

// ToString creates the content of the configuration file, the sections and properties are printed in a sorted order, so it can be tested.
func (pmc ProcConf) ToString() string {
	rawMap := pmc.ToMap()

	var sections []string
	for section := range rawMap {
		sections = append(sections, section)
	}

	sort.Strings(sections)

	var content strings.Builder
	for _, section := range sections {
		content.WriteString("[" + section + "]")
		content.WriteString("\n")

		var props []string
		for prop := range rawMap[section] {
			props = append(props, prop)
		}

		sort.Strings(props)

		for _, prop := range props {
			content.WriteString(prop)
			content.WriteString(" ")
			content.WriteString(rawMap[section][prop])
			content.WriteString("\n")
		}

		content.WriteString("\n")
	}

	return content.String()
}

// Merge will return the merged ProcConf, the values in the input will take precedent, does not mutate the original.
func (pmc ProcConf) Merge(input ProcConf) ProcConf {
	source := pmc.ToMap()
	override := input.ToMap()

	for section, props := range override {
		_, ok := source[section]
		if !ok {
			source[section] = map[string]string{}
		}

		for key, value := range props {
			source[section][key] = value
		}
	}

	updated := FromMap(source)
	updated.Revision = input.Revision
	return updated
}
