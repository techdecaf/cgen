package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/techdecaf/golog"
	"github.com/techdecaf/templates"
)

// ApplicationDirectory root directory of the cgen application
var ApplicationDirectory string

// TemplatesDirectory where all template generators are cached
var TemplatesDirectory string

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
	ApplicationDirectory = filepath.Join(usr.HomeDir, ".cgen")
	TemplatesDirectory = filepath.Join(ApplicationDirectory, "generators")

	app.BaseDir = ApplicationDirectory
	app.TemplatesDir = TemplatesDirectory
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

	GitCommand := templates.CommandOptions{
		UseStdOut: true,
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		GitCommand.Cmd = fmt.Sprintf("git clone '%s' '%s'", url, dir)
	} else {
		GitCommand.Cmd = fmt.Sprintf("cd '%s' && git pull", dir)
	}

	return templates.Run(GitCommand)
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
