package ruxit

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// ProcConf represents the response of the /deployment/installer/agent/processmoduleconfig endpoint from the Dynatrace Environment(v1) API.
type ProcConf struct {
	InstallPath *string    `json:"-"`
	Properties  []Property `json:"properties"`
	Revision    uint       `json:"revision"`
}

type Property struct {
	Section string `json:"section"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

// ToString creates the content of the configuration file, the sections and properties are printed in a sorted order, so it can be tested.
func (pc ProcConf) ToString() string {
	if pc.InstallPath != nil {
		pm := pc.ToMap()
		pm = pm.SetupReadonly(*pc.InstallPath)

		return pm.ToString()
	}

	return pc.ToMap().ToString()
}

// Merge returns the merged ProcConf, the values in the input will take precedent, does not mutate the original.
func (pc ProcConf) Merge(input ProcConf) ProcConf {
	source := pc.ToMap()
	override := input.ToMap()

	updated := FromMap(source.Merge(override))
	updated.Revision = input.Revision
	updated.InstallPath = input.InstallPath

	return updated
}

func (pc ProcConf) ToMap() ProcMap {
	sections := map[string]map[string]string{}
	for _, prop := range pc.Properties {
		section := sections[prop.Section]
		if section == nil {
			section = map[string]string{}
		}

		section[prop.Key] = prop.Value
		sections[prop.Section] = section
	}

	return sections
}

// ProcMap presents the content in a more easy to work with format. (a map of maps).
type ProcMap map[string]map[string]string

var (
	redundantEntries = map[string][]string{
		"general": {"logDir", "dataStorageDir"},
	}
	additionalEntries = ProcMap{
		"general": map[string]string{
			"storage": "\"/var/lib/dynatrace/oneagent\"", // TODO: Make configurable?
		},
	}
)

func (pm ProcMap) SetupReadonly(installPath string) ProcMap {
	for key, values := range redundantEntries {
		_, ok := pm[key]
		if !ok {
			continue
		}

		for _, value := range values {
			delete(pm[key], value)
		}
	}

	for section, entries := range pm {
		for entry, value := range entries {
			volume := filepath.VolumeName(value)
			fmt.Printf("%s", volume)

			if strings.HasPrefix(entry, "libraryPath") {
				sanitizedEntry := strings.ReplaceAll(value, "../", "")
				sanitizedEntry, found := strings.CutPrefix(sanitizedEntry, "\"")

				if found {
					pm[section][entry] = "\"" + filepath.Join(installPath, "agent", sanitizedEntry)
				} else {
					pm[section][entry] = filepath.Join(installPath, "agent", sanitizedEntry)
				}
			}
		}
	}

	return pm.Merge(additionalEntries)
}

func (pm ProcMap) ToString() string {
	var sections []string
	for section := range pm {
		sections = append(sections, section)
	}

	sort.Strings(sections)

	var content strings.Builder
	for _, section := range sections {
		content.WriteString("[" + section + "]")
		content.WriteString("\n")

		var props []string
		for prop := range pm[section] {
			props = append(props, prop)
		}

		sort.Strings(props)

		for _, prop := range props {
			content.WriteString(prop)
			content.WriteString(" ")
			content.WriteString(pm[section][prop])
			content.WriteString("\n")
		}

		content.WriteString("\n")
	}

	return content.String()
}

func (pm ProcMap) Merge(override ProcMap) ProcMap {
	for section, props := range override {
		_, ok := pm[section]
		if !ok {
			pm[section] = map[string]string{}
		}

		for key, value := range props {
			pm[section][key] = value
		}
	}

	return pm
}
