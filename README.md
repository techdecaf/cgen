<p align="center">
  <img alt="cgen" src="https://images.techdecaf.com/fit-in/100x/techdecaf/cgen_logo.png" width="100" />
</p>

# cgen project generator

This project is designed to be a cross platform plugin-based project generator.
Simply run `cgen` to get started!

- [cgen project generator](#cgen-project-generator)
  - [Download and Install cgen](#download-and-install-cgen)
    - [upgrading using curl](#upgrading-using-curl)
  - [Installing a template](#installing-a-template)
  - [Creating Your own Template Plugin](#creating-your-own-template-plugin)
  - [Template Operators Operators](#template-operators-operators)
  - [Updating a template](#updating-a-template)
    - [Full project generation](#full-project-generation)
    - [Updating static files only](#updating-static-files-only)
  - [Bumping a project version](#bumping-a-project-version)

```bash
Usage of cgen:
  -bump string
    # bumps the {major | minor | patch | pre-release string} version of the current directory using git tags.
  -install string
    # install a generator using a git clone compatible url cgen -install <url>
  -list
    # lists all installed generators
  -name string
    # what would you like to name your new project
  -static-only
    # does not generate template files (most commonly used with update)
  -tmpl string
    # specify a which template you would like to use.
  -upgrade
    # attempts to update the current directory, if it's already a cgen project
  -version
    # prints cgen version number
```

## Download and Install cgen

Download Links

- [windows](https://s3-us-west-2.amazonaws.com/github.techdecaf.io/cgen/latest/windows/cgen.exe)
- [mac](https://s3-us-west-2.amazonaws.com/github.techdecaf.io/cgen/latest/osx/cgen)
- [linux](https://s3-us-west-2.amazonaws.com/github.techdecaf.io/cgen/latest/linux/cgen)

To install cgen, simlink or place it in any directory that is part of your path.
i.e. `/usr/local/bin` or `c:\windows`

### upgrading using curl

```bash
# for linux simply replace `osx` with `linux`
# this will download cgen, make it executable and replace your existing binary with the upgraded version.
curl -o cgen https://s3-us-west-2.amazonaws.com/github.techdecaf.io/cgen/latest/osx/cgen  && chmod +x cgen && mv cgen $(which cgen)
```

> NOTE: you can also replace `latest` with any valid cgen version.

## Installing a template

cgen :heart_eyes: plugins, but it does not use package management, instead you can just reference any git repository that you have access to.

```bash
cgen -install https://github.com/techdecaf/cgen-template
```

## Creating Your own Template Plugin

You can actually use `cgen` to create a `cgen` template :tada:

```bash
# install the cgen template generator for templates
cgen -install https://github.com/techdecaf/cgen-template

# create a new directory for your template
mkdir my-new-template

# execute cgen, follow the prompts
cgen -tmpl cgen-template
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

## Template Operators Operators

- eq - Returns the boolean truth of arg1 == arg2
- ne - Returns the boolean truth of arg1 != arg2
- lt - Returns the boolean truth of arg1 < arg2
- le - Returns the boolean truth of arg1 <= arg2
- gt - Returns the boolean truth of arg1 > arg2
- ge - Returns the boolean truth of arg1 >= arg2

## Updating a template

### Full project generation

cgen creates an answer file in the root of your project, if you wish to upgrade your project with a newer version of your installed template just `cd <project_dir>` and `cgn -upgrade`.

### Updating static files only

> static files: are any files that do not end in `.tmpl`

```bash
cd <my_project_dir>
cgen -upgrade -staticOnly
```

## Bumping a project version

Wait, what? Why does a generator do this?

we added a bump feature to cgen to help with your projects life cycle, frequently we end up using many different tools to change the version of a project depending on the language we are using. However we felt that git was the correct place to bump and release new versions of our code. So you can also use cgen to handle this for you.

To use run `cgen -bump <major | minor | patch | pre-release string>` and cgen will update your git tags with a new semver.
