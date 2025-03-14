package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateImageInfo(t *testing.T) {
	type testCase struct {
		title string
		expected    string
		in   ImageInfo
	}

	testCases := []testCase{
		{
			title: "empty URI",
			expected:    "",
			in:   ImageInfo{},
		},
		{
			title: "URI with tag",
			expected:    "registry.example.com/repository/image:tag",
			in: ImageInfo{
				Registry:    "registry.example.com",
				Repository:  "repository/image",
				Tag:         "tag",
				ImageDigest: "",
			},
		},
		{
			title: "URI with digest",
			expected:    "registry.example.com/repository/image@sha256:7173b809ca12ec5dee4506cd86be934c4596dd234ee82c0662eac04a8c2c71dc",
			in: ImageInfo{
				Registry:    "registry.example.com",
				Repository:  "repository/image",
				Tag:         "",
				ImageDigest: "sha256:7173b809ca12ec5dee4506cd86be934c4596dd234ee82c0662eac04a8c2c71dc",
			},
		},
		{
			title: "URI with digest and tag",
			expected:    "registry.example.com/repository/image:tag@sha256:7173b809ca12ec5dee4506cd86be934c4596dd234ee82c0662eac04a8c2c71dc",
			in: ImageInfo{
				Registry:    "registry.example.com",
				Repository:  "repository/image",
				Tag:         "tag",
				ImageDigest: "sha256:7173b809ca12ec5dee4506cd86be934c4596dd234ee82c0662eac04a8c2c71dc",
			},
		},
		{
			title: "URI with missing tag",
			expected:    "registry.example.com/repository/image",
			in: ImageInfo{
				Registry:   "registry.example.com",
				Repository: "repository/image",
			},
		},
		{
			title: "URI with docker.io (special case in certain libraries)",
			expected:    "docker.io/php:fpm-stretch",
			in: ImageInfo{
				Registry:   "docker.io",
				Repository: "php",
				Tag:        "fpm-stretch",
			},
		},
		{
			title: "URI with missing registry",
			expected:    "php:fpm-stretch",
			in: ImageInfo{
				Repository: "php",
				Tag:        "fpm-stretch",
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.title, func(t *testing.T) {
			assert.Equal(t, test.expected, test.in.ToURI())
		})
	}
}
