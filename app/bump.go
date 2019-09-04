package app

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/blang/semver"
)

// Bump bump project versions
func Bump(place, pattern string) (version string, err error) {
	place = strings.ToLower(strings.TrimSpace(place))

	out, err := exec.Command("git", "describe", "--tags", "--always", "--dirty", "--abbrev=0").Output()
	if err != nil {
		return "", err
	}
	version = strings.TrimSpace(string(out))

	// check to make sure git repository is not dirty before performing a bump
	//TODO: catch git with no commit history
	if strings.Contains(version, "dirty") {
		return "", fmt.Errorf("uncommitted changes: please stash or commit the current changes before bumping the version")
	}

	v, _ := semver.Make(version)

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
	tag := fmt.Sprintf(pattern, v)

	// bump using git tag
	// cmd := exec.Command("git", "tag", "-a", tag, "-m", fmt.Sprintf("cgen bump %s", place))
	// stderr, _ := cmd.StderrPipe()
	// stdout, _ := cmd.StdoutPipe()
	// cmd.Start()

	// scanErr := bufio.NewScanner(stderr)
	// for scanErr.Scan() {
	// 	fmt.Println(scanErr.Text())
	// }

	// scanOut := bufio.NewScanner(stdout)
	// for scanOut.Scan() {
	// 	fmt.Println(scanOut.Text())
	// }

	// cmd.Wait()
	return tag, err
}
