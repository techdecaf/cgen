package app

import (
	"errors"
	"strconv"

	"github.com/blang/semver/v4"
)

// VersionIncrement major, minor, patch, pre
type VersionIncrement string

// IsValid - Version Increment
func (increment VersionIncrement) IsValid() error {
	switch increment {
	case "major", "minor", "patch", "pre-release":
		return nil
	}
	return errors.New("Invalid VersionIncrement wanted major, minor, patch or pre")
}

// Bump the incoming version by version increment
func (increment VersionIncrement) Bump(version string) (incremented string, err error) {
	var v semver.Version
	var preRelease = []semver.PRVersion{}

	// ensure valid version increment
	if err = increment.IsValid(); err != nil {
		return "", err
	}

	// get semver
	if v, err = semver.Parse(version); err != nil {
		return "", err
	}

	switch increment {
	case "major":
		v.IncrementMajor()
		v.Pre = []semver.PRVersion{}
	case "minor":
		v.IncrementMinor()
		v.Pre = []semver.PRVersion{}
	case "patch":
		v.IncrementPatch()
		v.Pre = []semver.PRVersion{}
	case "pre-release":
		if len(v.Pre) > 0 {
			i, _ := strconv.Atoi(v.Pre[0].String())
			newPre, _ := semver.NewPRVersion(strconv.Itoa(i + 1))
			preRelease = append(preRelease, newPre)
		} else {
			newPre, _ := semver.NewPRVersion("1")
			preRelease = append(preRelease, newPre)
		}
	default:
	}

	// set the pre-release version
	v.Pre = preRelease
	return v.String(), nil
}
