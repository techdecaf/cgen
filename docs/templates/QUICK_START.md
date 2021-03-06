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
{{`
// {{.Name}} project was generated by robots at
// {{ .Timestamp }}
// using data from
// {{ .URL }}
{{- range .MyArray }}
    {{ printf "%q" . }},
{{- end }}
`}}
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
