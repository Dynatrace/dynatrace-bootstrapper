// Package ruxit provides parsing and conversion utilities for ruxit agent process module config files.
package ruxit

import (
	"bufio"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const (
	// expectedSectionParts defines the expected number of parts when parsing section matches.
	expectedSectionParts = 2
	// expectedConfigParts defines the expected number of parts when splitting config line.
	expectedConfigParts = 2
)

// sectionRegexp matches section headers like [general].
var sectionRegexp = regexp.MustCompile(`\[(.*)\]`)

// FromMap converts a ProcMap to a ProcConf struct.
func FromMap(procMap ProcMap) ProcConf {
	var result ProcConf

	for section, props := range procMap {
		for key, value := range props {
			result.Properties = append(result.Properties, Property{
				Section: section,
				Key:     key,
				Value:   value,
			})
		}
	}

	return result
}

// FromJSON parses a ProcConf from a JSON reader.
func FromJSON(reader io.Reader) (ProcConf, error) {
	var result ProcConf

	raw, err := io.ReadAll(reader)
	if err != nil {
		return result, errors.WithStack(err)
	}

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return result, errors.WithStack(err)
	}

	return result, nil
}

// FromConf creates the ProcConf struct from a valid ruxitagentproc.conf config file.
func FromConf(reader io.Reader) (ProcConf, error) {
	var result []Property

	const whiteSpace = "\t\n\v\f\r "

	scanner := bufio.NewScanner(reader)
	currentSection := ""

	for scanner.Scan() {
		line := scanner.Text()
		header := confSectionHeader(line)

		switch {
		case header != "":
			currentSection = strings.Trim(header, whiteSpace)
		case line != "" && !strings.HasPrefix(line, "#"):
			splitLine := strings.Split(line, " ")
			prop := Property{
				Section: currentSection,
				Key:     strings.Trim(splitLine[0], whiteSpace),
			}

			if len(splitLine) == expectedConfigParts {
				prop.Value = strings.Trim(splitLine[1], whiteSpace)
			}

			result = append(result, prop)
		}
	}

	return ProcConf{Properties: result}, scanner.Err()
}

// confSectionHeader extracts the section header from a line, if present.
func confSectionHeader(line string) string {
	matches := sectionRegexp.FindStringSubmatch(line)

	if len(matches) == expectedSectionParts {
		return matches[1]
	}

	return ""
}
