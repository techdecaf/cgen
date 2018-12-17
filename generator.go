package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	yaml "gopkg.in/yaml.v2"
)

type Output struct {
	Template string                 `yaml:"template"`
	Answers  map[string]interface{} `yaml:"answers"`
}

// Question struct for questions file.
type Question struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Prompt  string   `yaml:"prompt"`
	Default string   `yaml:"default"`
	Options []string `yaml:"options,omitempty"`
}

// Config - the config.yaml
type Config struct {
	Version   string      `yaml:"version"`
	From      string      `yaml:"from"`
	Questions []*Question `yaml:"questions"`
	Post      []string    `yaml:"post"`
}

// Generator struct
type Generator struct {
	Name          string
	Source        string
	Destination   string
	QuestionsFile string
	AnswersFile   string
	TemplateFiles string
	TemplateName  string
	TemplatesDir  string
	Config        *Config
	Answers       map[string]interface{}
}

func (gen *Generator) init(name, template, src string, upgrade bool) error {
	// todo: validate inputs, that files exist etc
	// default destination to current working directory or use project name

	// check to see if an answers file exists in current dir
	answerFile := path.Join(".", ".cgen.yaml")
	if upgrade {
		// ensure answer file exists
		if _, err := os.Stat(answerFile); err != nil {
			return err
		}
		update := Output{}

		answersYAML, err := os.Open(answerFile)
		if err != nil {
			return err
		}
		defer answersYAML.Close()
		byteValue, _ := ioutil.ReadAll(answersYAML)
		yaml.Unmarshal(byteValue, &update)

		gen.Answers = update.Answers
		gen.TemplateName = update.Template
	} else {
		gen.TemplateName = template
	}

	// path to generators
	gen.Name = name
	gen.TemplatesDir = src
	gen.Source = path.Join(gen.TemplatesDir, gen.TemplateName)

	gen.Destination = "."
	gen.AnswersFile = path.Join(gen.Destination, ".cgen.yaml")
	gen.QuestionsFile = path.Join(gen.Source, "config.yaml")
	gen.TemplateFiles = path.Join(gen.Source, "template")
	gen.Config = &Config{}

	// check for required project structure
	if _, err := os.Stat(gen.TemplateFiles); os.IsNotExist(err) {
		return fmt.Errorf("%s does not have the required template directory, please check the README file", gen.Source)
	}

	if _, err := os.Stat(gen.QuestionsFile); os.IsNotExist(err) {
		log.Printf("%s does not have a questions.yaml file, so it may not actually be a cgen template...", gen.Source)
	}

	configYAML, err := os.Open(gen.QuestionsFile)
	if err != nil {
		return err
	}
	defer configYAML.Close()
	byteValue, _ := ioutil.ReadAll(configYAML)
	yaml.Unmarshal(byteValue, &gen.Config)

	gen.appendAnswer("TemplateVersion", gen.Config.Version)
	gen.appendAnswer("Timestamp", time.Now().UTC().Format(time.RFC3339))
	return nil
}

func (gen *Generator) exec() error {
	if err := gen.prompt(); err != nil {
		return err
	}

	if err := filepath.Walk(gen.TemplateFiles, gen.walkFiles); err != nil {
		return err
	}

	ans, err := gen.save()
	if err := ioutil.WriteFile(gen.Destination+"/.cgen.yaml", ans, 0644); err != nil {
		return err
	}

	return err
}

func (gen *Generator) toJSON() error {
	json, err := json.Marshal(gen)
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}

func (gen *Generator) copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// create directory if it does not exist
	if err := os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func (gen *Generator) walkFiles(inPath string, file os.FileInfo, err error) error {
	// skip all directories
	if file.IsDir() {
		return nil
	}
	outPath := strings.Replace(inPath, gen.TemplateFiles, gen.Destination, 1)
	if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
		return err
	}

	if filepath.Ext(inPath) == ".tmpl" {
		outPath = strings.Replace(outPath, filepath.Ext(outPath), "", 1)
		fmt.Printf("Processing Template File %s\n", inPath)
		fmt.Printf("Generating Template: %s\n", inPath)

		var helpers = template.FuncMap{
			"ToTitle": strings.Title,
			"ToUpper": strings.ToUpper,
			"ToLower": strings.ToLower,
			"Replace": strings.Replace,
			"MkdirAll": func(p string) (err error) {
				err = os.MkdirAll(path.Join(gen.Destination, p), 0700)
				return err
			},
			"Touch": func(p string) (err error) {
				_, err = os.Create(path.Join(gen.Destination, p))
				return err
			},
		}

		templateFile, err := template.New(file.Name()).Funcs(helpers).ParseFiles(inPath)
		if err != nil {
			return err
		}

		generated, err := os.Create(outPath)
		if err != nil {
			return err
		}

		if err := templateFile.Execute(generated, gen.Answers); err != nil {
			return err
		}
	} else {
		gen.copy(inPath, outPath)
		fmt.Printf("Copying File: %s\n", outPath)
	}

	return nil // no errors
}

func (gen *Generator) prompt() error {
	for _, q := range gen.Config.Questions {
		res, err := gen.ask(*q)
		fmt.Printf("You choose %q\n", res)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gen *Generator) ask(q Question) (answer string, err error) {

	if val := os.Getenv(q.Name); val != "" {
		return gen.appendAnswer(q.Name, val), nil
	}

	if val := gen.Answers[q.Name]; val != nil {
		return fmt.Sprintf("%v", val), nil
	}

	switch q.Type {
	case "string":
		prompt := promptui.Prompt{
			Label:   q.Prompt,
			Default: q.Default,
		}
		answer, err = prompt.Run()
	case "bool":
		truthRE := "(?i)^true|y"

		if match, _ := regexp.MatchString(truthRE, q.Default); match {
			q.Default = "y"
		}

		prompt := promptui.Prompt{
			Label:     q.Prompt,
			IsConfirm: true,
			Default:   q.Default,
		}
		answer, _ = prompt.Run()
		if answer == "" {
			answer = q.Default
		}

		if match, _ := regexp.MatchString(truthRE, answer); match {
			answer = "true"
		} else {
			answer = "false"
		}
	case "select":
		prompt := promptui.Select{
			Label: q.Prompt,
			Items: q.Options,
		}
		_, answer, err = prompt.Run()
	default:
		return "", fmt.Errorf("invalid question type %s", q.Type)
	}

	if err != nil {
		return "", err
	}

	return gen.appendAnswer(q.Name, answer), nil
}

func (gen *Generator) appendAnswer(name, val string) (answer string) {
	if gen.Answers == nil {
		gen.Answers = make(map[string]interface{})
	}
	// append answer to the answers map.
	switch val {
	case "true":
		gen.Answers[name] = true
	case "false":
		gen.Answers[name] = false
	default:
		gen.Answers[name] = val
	}

	return val
}

func (gen *Generator) save() (out []byte, err error) {
	output := Output{}
	output.Answers = gen.Answers
	output.Template = gen.TemplateName

	res, err := yaml.Marshal(output)
	fmt.Println(string(res))
	return res, err
}
