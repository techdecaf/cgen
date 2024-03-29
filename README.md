# cgen

<p align="center">
  <img
    alt="cgen"
    src="https://images.techdecaf.com/fit-in/100x/techdecaf/cgen_logo.png"
    width="100"
  />
</p>


- [cgen](#ciprojectname)
  - [Download and Install](#download-and-install)
  - [Quick Start](#quick-start)
  - [Contribution Guide](#contribution-guide)
  - [Credits](#credits)

## Download and Install

```bash
sh -c "$(curl -fsSL https://raw.github.com/techdecaf/cgen/master/install.sh)"
```

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.github.com/techdecaf/cgen/master/install.ps1'))
```

Download Links

- [windows](http://github.techdecaf.io/cgen/latest/windows/cgen.exe)
- [mac](http://github.techdecaf.io/cgen/latest/darwin/cgen)
- [linux](http://github.techdecaf.io/cgen/latest/linux/cgen)

To install cgen, use the provided script, simlink it or place it in any directory that is part of your path.
i.e. `/usr/local/bin` or `c:\windows`


## Quick Start

```text
You can use cgen to dynamically configure new projects based
   on your own standards and best practices. See the README.md to get started.

Usage:
  cgen [flags]
  cgen [command]

Available Commands:
  bump        Creates a new git tag with an increase in the current semantic version i.e. v1.0.2
  completion  Generates zsh completion scripts
  help        Help about any command
  install     Installs a new generator from a git repository
  list        prints a list of currently installed directories
  promote     promote a file from a project to your cgen template
  upgrade     this features is not currently supported pull request?

Flags:
  -h, --help                       help for cgen
      --ignore-version-tolerance   skips cgen version tolerance check (default true)
  -n, --name string                what do you want to call your newly generated project?
  -p, --path string                where you would like to generate your project. (default "/Users/ward/@decaf/cgen")
  -s, --static-only                does not generate template files (most commonly used with update)
  -t, --template string            specify a which template you would like to use.
      --var stringToString         overwrite environmental variables (default [])
      --verbose                    enable verbose log messages
  -v, --version                    prints the cgen version number

Use "cgen [command] --help" for more information about a command.
```

This project is designed to be a cross platform plugin-based project generator. Simply run `cgen` to get started!

### Installing a Template

cgen :heart_eyes: plugins, but it does not use package management, instead you can just reference any git repository that you have access to.

```bash
cgen install https://github.com/techdecaf/cgen-template
```

### Creating Your own Template Plugin

You can actually use `cgen` to create a `cgen` template :tada:

```bash
# install the cgen template generator for templates
cgen install https://github.com/techdecaf/cgen-template

# create a new directory for your template
mkdir my-new-template

# execute cgen, follow the prompts
cgen -t cgen-template
```

You can take a look at the [cgen-template](https://github.com/techdecaf/cgen-template) project for more information on use and details for how to create your own templates.

We use the go template engine to create your project, you can find detailed documentation here:

- [Go Template Documentation](https://golang.org/pkg/html/template/)
- [todo: link to examples](/examples)

```go

// {{.Name}} project was generated by robots at
// {{ .Timestamp }}
// using data from
// {{ .URL }}
{{- range .MyArray }}
    {{ printf "%q" . }},
{{- end }}

```

### Template Operators Operators

- eq - Returns the boolean truth of arg1 == arg2
- ne - Returns the boolean truth of arg1 != arg2
- lt - Returns the boolean truth of arg1 < arg2
- le - Returns the boolean truth of arg1 <= arg2
- gt - Returns the boolean truth of arg1 > arg2
- ge - Returns the boolean truth of arg1 >= arg2

### Updating a template

#### Full project generation

cgen creates an answer file in the root of your project, if you wish to upgrade your project with
a newer version of your installed template just `cd <project_dir>` and `cgn upgrade`.

#### Updating static files only

> static files: are any files that do not end in `.tmpl`

```bash
cd <my_project_dir>
cgen -upgrade -staticOnly
```

### Bumping a project version

Wait, what? Why does a generator do this?

we added a bump feature to cgen to help with your projects life cycle, frequently we end up using many different tools to change the version of a project depending on the language we are using. However we felt that git was the correct place to bump and release new versions of our code. So you can also use cgen to handle this for you.

To use run `cgen bump --level <major | minor | patch | pre-release string>` and cgen will update your git tags with a new semver.


## Contribution Guide

## Credits

### Logo

The logo for this project provided by [logomakr](https://logomakr.com)

### Sponsor

[![TechDecaf](https://images.techdecaf.com/fit-in/150x/techdecaf/logo_full.png)](https://techdecaf.com)

_Get back to doing what you do best, let us handle the rest._

