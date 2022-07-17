package app

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/techdecaf/templates"
)

// BumpParams options for running bump
type BumpParams struct {
	Place   string
	Pattern string
	DryRun  bool
	GitPush bool
}

// Bump bump project versions
func Bump(bump BumpParams) (version string, err error) {
	place := strings.ToLower(strings.TrimSpace(bump.Place))
	pattern := bump.Pattern

	// check for uncommitted changes
	// if err := bump.checkForUncommittedChanges(); err != nil {
	//   return "", err
	// }

	// find the current version
	currentVersion, err := bump.getCurrentVersion()
	if err != nil {
		return "", err
	}

	v, err := VersionIncrement(bump.Place).Bump(currentVersion)
	if err != nil {
		return "", err
	}

	// format tag according to the pattern
	tag := fmt.Sprintf(pattern, v)
	msg := fmt.Sprintf("incrementing %s version", place)
	cmd := fmt.Sprintf("git tag -a %s -m '%s'", tag, msg)

	if bump.DryRun {
		return tag, err
	}

	GitTag := templates.CommandOptions{
		Cmd:       cmd,
		UseStdOut: true,
	}

	if out, err := templates.Run(GitTag); err != nil {
		return out, err
	}

	// push
	GitPush := templates.CommandOptions{
		Cmd:       "git push --follow-tags",
		UseStdOut: true,
	}

	if bump.GitPush {
		if out, err := templates.Run(GitPush); err != nil {
			return out, err
		}
	}

	return tag, err
}

func (bump BumpParams) getCurrentVersion() (version string, err error) {
	reEx := regexp.MustCompile(`(\d+\.\d+\.\d+)(-\d+)?`)

	// get most recent git tag by tagger date
	gitTags := templates.CommandOptions{
		Cmd:        "git tag --sort=taggerdate",
		UseStdOut:  false,
		TrimOutput: false,
	}

	gitTagsOut, err := templates.Run(gitTags)
	if err != nil {
		return "", err
	}

	tags := strings.Split(strings.TrimSpace(gitTagsOut), "\n")

	// parse the semver out of a string ignoring prefixes like `v0.2.3` or `version_0.2.6-rc.4`
	version = reEx.FindString(tags[len(tags)-1])
	// if no curent version is found, default to 0.0.0
	if version == "" {
		return "0.0.0", nil
	}

	return version, nil
}

func (bump BumpParams) checkForUncommittedChanges() (err error) {
	gitDescribe := templates.CommandOptions{
		Cmd:        "git describe --all --dirty --abbrev=0",
		UseStdOut:  false,
		TrimOutput: false,
	}

	gitDescribeOut, err := templates.Run(gitDescribe)
	if err != nil {
		return err
	}
	if strings.Contains(gitDescribeOut, "dirty") {
		return errors.New("uncommitted changes: please stash or commit the current changes before bumping the version")
	}

	return nil
}
