package ruxit

import (
	"bufio"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// example match: [general]
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

func FromJson(reader io.Reader) (ProcConf, error) {
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

// FromConf creates the ProcConf struct from an valid ruxitagentproc.conf config file.
func FromConf(reader io.Reader) (ProcConf, error) {
	var result []Property

	scanner := bufio.NewScanner(reader)
	currentSection := ""

	for scanner.Scan() {
		line := scanner.Text()
		header := confSectionHeader(line)

		switch {
		case header != "":
			currentSection = strings.Trim(header, "\t\n\v\f\r ")
		case line != "" && !strings.HasPrefix(line, "#"):
			splitLine := strings.Split(line, " ")
			prop := Property{
				Section: currentSection,
				Key:     strings.Trim(splitLine[0], "\t\n\v\f\r "),
			}

			if len(splitLine) == 2 {
				prop.Value = strings.Trim(splitLine[1], "\t\n\v\f\r ")
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

func confSectionHeader(confLine string) string {
	if matches := sectionRegexp.FindStringSubmatch(confLine); len(matches) != 0 {
		return matches[1]
	}

	return ""
}
