package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/blang/semver"
)

// CGen application
type CGen struct {
	BaseDir      string
	TemplatesDir string
	Generator    Generator
}

func (app *CGen) init() (err error) {
	var mode os.FileMode = 0700

	usr, err := user.Current()
	if err != nil {
		return err
	}

	app.BaseDir = path.Join(usr.HomeDir, ".cgen")
	app.TemplatesDir = path.Join(app.BaseDir, "generators")
	app.Generator = Generator{}

	if _, err := os.Stat(app.BaseDir); os.IsNotExist(err) {
		os.Mkdir(app.BaseDir, mode)
	}

	if _, err := os.Stat(app.TemplatesDir); os.IsNotExist(err) {
		os.Mkdir(app.TemplatesDir, mode)
	}

	return nil
}

func (app *CGen) install(url string) (err error) {
	// what to name the generator dir.
	as := strings.TrimSuffix(path.Base(url), path.Ext(url))
	dir := path.Join(app.TemplatesDir, as)

	cmd := exec.Command("git", "clone", url, dir)
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanErr := bufio.NewScanner(stderr)
	for scanErr.Scan() {
		fmt.Println(scanErr.Text())
	}

	scanOut := bufio.NewScanner(stdout)
	for scanOut.Scan() {
		fmt.Println(scanOut.Text())
	}

	return cmd.Wait()
	// if err != nil {
	//     // something went wrong
	// }

	// _, err = git.PlainClone(dir, false, &git.CloneOptions{
	// 	URL:      url,
	// 	Progress: os.Stdout,
	// })
}

func (app *CGen) update(name string) (err error) {
	return errors.New("this endpoint has not yet been created, want to contribute?")
}

func (app *CGen) listInstalled() (installed []string, err error) {
	files, err := ioutil.ReadDir(app.TemplatesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		installed = append(installed, file.Name())
	}

	if len(installed) == 0 {
		return nil, errors.New("no generators are installed, would you like to add some? Try: `cgen -install <url>`")
	}

	return installed, err
}

func (app *CGen) bump(place string) (version string, err error) {
	place = strings.ToLower(strings.TrimSpace(place))
	if out, err := exec.Command("git", "describe", "--tags", "--always", "--dirty", "--abbrev=0").Output(); err != nil {
		return "", err
	} else {
		version = strings.TrimSpace(string(out))

		// check to make sure git repository is not dirty before performing a bump
		if strings.Contains(version, "dirty") {
			return "", fmt.Errorf("UncommittedChanges: please stash or commit the current changes before bumping the version.")
		}

		v, err := semver.Make(version)

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

		// bump using git tag
		cmd := exec.Command("git", "tag", "-a", v.String(), "-m", fmt.Sprintf("[bump] cgen -bump %s", place))
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		scanErr := bufio.NewScanner(stderr)
		for scanErr.Scan() {
			fmt.Println(scanErr.Text())
		}

		scanOut := bufio.NewScanner(stdout)
		for scanOut.Scan() {
			fmt.Println(scanOut.Text())
		}

		cmd.Wait()
		return v.String(), err
	}
}
