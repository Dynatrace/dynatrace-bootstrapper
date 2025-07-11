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
	twoSectionParts = 2
	twoConfigParts  = 2
)

var sectionRegexp = regexp.MustCompile(`\[(.*)\]`)

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

			if len(splitLine) == twoConfigParts {
				prop.Value = strings.Trim(splitLine[1], whiteSpace)
			}

			result = append(result, prop)
		}
	}

	if err := scanner.Err(); err != nil {
		return ProcConf{}, errors.WithStack(err)
	}

	return ProcConf{
		Properties: result,
		Revision:   0,
	}, nil
}

func confSectionHeader(line string) string {
	matches := sectionRegexp.FindStringSubmatch(line)

	if len(matches) == twoSectionParts {
		return matches[1]
	}

	return ""
}
