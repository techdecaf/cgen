package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
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

	_, err = git.PlainClone(path.Join(app.TemplatesDir, as), false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	return err
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
