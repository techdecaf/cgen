package app

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
	"github.com/techdecaf/templates"
)

// BumpParams options for running bump
type BumpParams struct {
	Place   string
	Pattern string
	DryRun  bool
}

// Bump bump project versions
func Bump(params BumpParams) (version string, err error) {
	reEx := regexp.MustCompile(`(\d+\.\d+\.\d+)`)

	place := strings.ToLower(strings.TrimSpace(params.Place))
	pattern := params.Pattern

	out, err := templates.Run(templates.CommandOptions{
		Cmd:        "git describe --tags --always --dirty --abbrev=0",
		UseStdOut:  false,
		TrimOutput: false,
	})

	version = strings.TrimSpace(string(out))

	// check to make sure git repository is not dirty before performing a bump
	//TODO: catch git with no commit history
	if strings.Contains(version, "dirty") {
		return "", fmt.Errorf("uncommitted changes: please stash or commit the current changes before bumping the version")
	}

	v, _ := semver.Make(reEx.FindString(version))

	switch place {
	case "major":
		v.Major++
		v.Minor = 0
		v.Patch = 0
	case "minor":
		v.Minor++
		v.Patch = 0
	case "patch":
		v.Patch++
	default:
		v.Pre[0], err = semver.NewPRVersion(place)
	}

	// format tag according to the pattern
	tag := fmt.Sprintf(pattern, v.String())
	msg := fmt.Sprintf("cgen bump -l %s", place)
	cmd := fmt.Sprintf("git tag -a %s -m '%s'", tag, msg)

	if params.DryRun {
		return tag, err
	}

	templates.Run(templates.CommandOptions{
		Cmd:       cmd,
		UseStdOut: true,
	})

	return tag, err
}
