package main

import (
	"flag"
	"fmt"
	"log"
	"path"
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
	gitURL := flag.String("install", "", "install a generator using a git clone compatable url cgen -install <url>")
	project := flag.String("tmpl", "", "specify a which template you would like to use.")
	name := flag.String("name", "", "what would you like to name your new project")
	doList := flag.Bool("list", false, "lists all installed generators")
	doVersion := flag.Bool("version", false, "prints cgen version number")
	flag.Parse()

	if *doVersion != false {
		fmt.Println(VERSION)
		return
	}

	if *doList != false {
		if installed, err := app.listInstalled(); err != nil {
			log.Fatal(err)
		} else {
			for _, template := range installed {
				fmt.Println(path.Join(app.TemplatesDir, template))
			}
			return
		}
	}

	// install handler
	if *gitURL != "" {
		if err := app.install(*gitURL); err != nil {
			log.Fatal(err)
		}
		return
	}

	installedGenerators, err := app.listInstalled()
	if err != nil {
		log.Fatal(err)
	}

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
			Prompt:  "What do you want to call your project",
			Default: "temp-project",
		})
	}

	if err != nil {
		log.Fatal(err)
	}

	if err := app.Generator.init(*name, path.Join(app.TemplatesDir, *project)); err != nil {
		log.Fatal(err)
	}
	// app.Generator.toJSON()
	if err := app.Generator.exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(app.Generator)
}
