package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

// todo: write tests
// version - This is converted to the git tag at compile time using the make build command.
var VERSION string

func main() {
	var app = &CGen{}
	if err := app.init(); err != nil {
		log.Fatal(err)
	}
	// CLI Flags
	name := flag.String("name", "", "what would you like to name your new project")
	project := flag.String("tmpl", "", "specify a which template you would like to use.")

	// utilities
	install := flag.String("install", "", "install a generator using a git clone compatable url cgen -install <url>")
	bump := flag.String("bump", "", "bumps the {major | minor | patch | pre-release string} version of the current directory using git tags.")

	staticOnly := flag.Bool("static-only", false, "does not generate template files (most commonly used with update)")
	doList := flag.Bool("list", false, "lists all installed generators")
	doUpgrade := flag.Bool("upgrade", false, "attempts to update the current directory, if it's already a cgen project")
	doVersion := flag.Bool("version", false, "prints cgen version number")
	flag.Parse()

	// Utility functions
	if *doVersion != false {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// install handler
	if *install != "" {
		if err := app.install(*install); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *bump != "" {
		ver, err := app.bump(*bump)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ver)
		os.Exit(0)
	}

	installedGenerators, err := app.listInstalled()
	if err != nil {
		log.Fatal(err)
	}

	if *doList != false {
		for _, template := range installedGenerators {
			fmt.Println(path.Join(app.TemplatesDir, template))
		}
		os.Exit(0)
	}

	// Main Package
	// resolve current directory.
	pwd, _ := filepath.Abs(".")
	thisDir, err := os.Stat(pwd)
	if err != nil {
		log.Fatal(err)
	}

	// check to see if directory is dirty.
	if files, err := ioutil.ReadDir("./"); err != nil {
		log.Fatal(err)
	} else {
		if len(files) != 0 {
			fmt.Println("WARNING: This directory is not empty.")
		}
	}

	// PERFORM UPGRADE
	if *doUpgrade == true {
		*name = thisDir.Name()
		*project = "PerformUpgrade"

		confirm, err := app.Generator.ask(Question{
			Name:    "Confirm",
			Type:    "bool",
			Prompt:  fmt.Sprintf("Are you sure you want to upgrade [%s]", *name),
			Default: "false",
		})

		if err != nil {
			log.Fatal(err)
		}

		if confirm == "false" {
			os.Exit(0)
		}
	}

	// PERFORM PROJECT GENERATION
	if *project == "" {
		*project, err = app.Generator.ask(Question{
			Name:    "Template",
			Type:    "select",
			Prompt:  "Pick a template.",
			Options: installedGenerators,
		})
	}

	if *name == "" {
		*name, err = app.Generator.ask(Question{
			Name:    "Name",
			Type:    "string",
			Prompt:  "What do you want to call your project ",
			Default: thisDir.Name(),
		})
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := app.Generator.init(*name, *project, app.TemplatesDir, *doUpgrade, *staticOnly); err != nil {
		log.Fatal(err)
	}
	// app.Generator.toJSON()

	if err := app.Generator.exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(app.Generator)
}
