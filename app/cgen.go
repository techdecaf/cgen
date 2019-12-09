package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/techdecaf/golog"
	"github.com/techdecaf/templates"
)

// Log cgen logger
var Log = golog.Log{
	Name: "cgen",
}

// CGen application
type CGen struct {
	BaseDir      string
	TemplatesDir string
	Generator    Generator
}

// Init new instance of CGEN
func (app *CGen) Init() (err error) {
	var mode os.FileMode = 0700

	usr, err := user.Current()
	if err != nil {
		return err
	}

	app.BaseDir = filepath.Join(usr.HomeDir, ".cgen")
	app.TemplatesDir = filepath.Join(app.BaseDir, "generators")
	app.Generator = Generator{}

	if _, err := os.Stat(app.BaseDir); os.IsNotExist(err) {
		os.Mkdir(app.BaseDir, mode)
	}

	if _, err := os.Stat(app.TemplatesDir); os.IsNotExist(err) {
		os.Mkdir(app.TemplatesDir, mode)
	}

	return nil
}

// Install a generator from git.
func (app *CGen) Install(url string) (out string, err error) {
	// what to name the generator dir.
	as := strings.TrimSuffix(filepath.Base(url), filepath.Ext(url))
	dir := filepath.Join(app.TemplatesDir, as)

	GitClone := templates.CommandOptions{
		Cmd:       fmt.Sprintf("git clone '%s' '%s'", url, dir),
		UseStdOut: true,
	}
	return templates.Run(GitClone)
}

// Update a project
func (app *CGen) Update(name string) (err error) {
	return errors.New("this endpoint has not yet been created, want to contribute?")
}

// ListInstalled generators from ~/.cgen/generators
func (app *CGen) ListInstalled() (installed []string, err error) {
	files, err := ioutil.ReadDir(app.TemplatesDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		installed = append(installed, file.Name())
	}

	if len(installed) == 0 {
		return nil, errors.New("no generators are installed, would you like to add some? Try: `cgen install <url>`")
	}

	return installed, err
}
