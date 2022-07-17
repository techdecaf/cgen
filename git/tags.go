package git

import (
	"fmt"
	"regexp"
)

// Tags for a git repository
type Tags []string

// Filter for a new repository
func (tags Tags) Filter(expression string) []string {
	var versions []string
	reEx := regexp.MustCompile(expression)

	for _, tag := range tags {
		version := reEx.FindString(tag)
		if version != "" {
			versions = append(versions, version)
		}
	}

	return versions
}

// Latest for a new repository
func (tags Tags) Latest() (string, error) {
	tag := tags.Filter(`v(\d+\.\d+.\d+)(-\d+)?$`)
	if len(tag) == 0 {
		return "", fmt.Errorf("no valid versions were found in git tags")
	}
	return tag[len(tag)-1], nil
}

// LatestStable for a new repository
func (tags Tags) LatestStable() (string, error) {
	tag := tags.Filter(`v(\d+\.\d+.\d+)$`)
	if len(tag) == 0 {
		return "", fmt.Errorf("no valid versions were found in git tags")
	}
	return tag[len(tag)-1], nil
}
