package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/techdecaf/cgen/git"
	"gopkg.in/yaml.v2"
)

// Template - the config.yaml
type Template struct {
  CgenVersion     string        `yaml:"cgen_version"`
  Extends         string        `yaml:"extends"`
  Variables       yaml.MapSlice `yaml:"variables"`
  Questions       []*Question   `yaml:"questions"`
  TemplateFiles   []string      `yaml:"template_files"`
	RunAfter        []string      `yaml:"run_after"`
  RunBefore       []string      `yaml:"run_before"`
  // private properties
  Name                  string
  LatestVersion         string
  LatestStableVersion   string
  Directory             git.Repository
  Files                 string
  ConfigFile             string
}

// Init values for template
func (template *Template) Init() error {
  template.Directory = git.Repository(filepath.Join(TemplatesDirectory, template.Name))

  if template.Directory.ToString() == "" {
    return fmt.Errorf("can not init a new template without a directory")
  }
  template.Files = filepath.Join(template.Directory.ToString(), "template")
  template.ConfigFile = filepath.Join(template.Directory.ToString(), "config.yaml")

  // check for required project structure
	if _, err := os.Stat(template.Files); os.IsNotExist(err) {
		return fmt.Errorf("%s does not have the required template directory, please check the README file", template.Directory)
  }

  if _, err := os.Stat(template.ConfigFile); os.IsNotExist(err) {
		return fmt.Errorf("%s does not have a config.yaml file, so it may not actually be a cgen template", template.Directory)
	}

  configYAML, err := os.Open(template.ConfigFile)
	if err != nil {
		return err
	}
	defer configYAML.Close()
	byteValue, _ := ioutil.ReadAll(configYAML)
  yaml.Unmarshal(byteValue, &template)

  tags, err := template.Directory.ListTags()
  if err != nil {
    return err
  }

  latest, err := tags.Latest()
  if err != nil {
    return err
  }

  stable, err := tags.LatestStable()
  if err != nil {
    return err
  }

  template.LatestVersion = latest
  template.LatestStableVersion = stable

  return nil
}

// FileIsTemplate checks to see if file has been marked as a template and should be expanded
func(template *Template) FileIsTemplate(filePath string) (bool,error) {
  if filepath.Ext(filePath) == ".tmpl" {
    return true, nil
  }

  relPath, err := filepath.Rel(template.Files, filePath)
  if err != nil {
    return false, err
  }

  for _, match := range template.TemplateFiles {
    if match == relPath {
        return true, nil
    }
  }

  return false, err
}

func (template *Template) toJSON() error {
	out, err := json.Marshal(template)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}