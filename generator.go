package main

import (
	"bytes"
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
	RunAfter  []string    `yaml:"run_after"`
}

// Generator struct
type Generator struct {
	Name            string `json:"Name"`
	Source          string `json:"Source"`
	Destination     string `json:"Destination"`
	QuestionsFile   string `json:"QuestionsFile"`
	AnswersFile     string `json:"AnswersFile"`
	TemplateFiles   string `json:"TemplateFiles"`
	TemplateName    string `json:"TemplateName"`
	TemplatesDir    string `json:"TemplatesDir"`
	TemplateHelpers template.FuncMap
	Config          *Config                `json:"Config"`
	Answers         map[string]interface{} `json:"Answers"`
	Options         struct {
		StaticOnly     bool `json:"StaticOnly"`
		PerformUpgrade bool `json:"PerformUpgrade"`
	}
}

func (gen *Generator) init(params GeneratorParams) error {
	params.toJSON()
	// set options
	gen.Options.StaticOnly = params.StaticOnly
	gen.Options.PerformUpgrade = params.PerformUpgrade

	// todo: validate inputs, that files exist etc
	// default destination to current working directory or use project name
	// check to see if an answers file exists in current dir
	answerFile := path.Join(params.Destination, ".cgen.yaml")
	if gen.Options.PerformUpgrade {
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
		gen.TemplateName = params.Tempate
	}

	// path to generators
	gen.Name = params.Name
	gen.TemplatesDir = params.TemplatesDir
	gen.Source = path.Join(gen.TemplatesDir, gen.TemplateName)

	gen.Destination = params.Destination
	gen.AnswersFile = path.Join(gen.Destination, ".cgen.yaml")
	gen.QuestionsFile = path.Join(gen.Source, "config.yaml")
	gen.TemplateFiles = path.Join(gen.Source, "template")
	gen.Config = &Config{}

	gen.TemplateHelpers = template.FuncMap{
		"ToTitle": strings.Title,
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"Replace": strings.Replace,
		"MkdirAll": func(relativePath string) (err error) {
			err = os.MkdirAll(path.Join(gen.Destination, relativePath), 0700)
			return err
		},
		"Touch": func(relativePath string) (err error) {
			_, err = os.Create(path.Join(gen.Destination, relativePath))
			return err
		},
	}

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

	// save output to projct_root/.cgen.yaml
	ans, err := gen.save()
	if err := ioutil.WriteFile(gen.Destination+"/.cgen.yaml", ans, 0644); err != nil {
		return err
	}

	if !gen.Options.PerformUpgrade {
		// run scripts in config.run_after array.
		if err := gen.runAfter(); err != nil {
			return err
		}
	}
	return err
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
	// identify template files
	isTemplate := filepath.Ext(inPath) == ".tmpl"

	// skip all directories
	if file.IsDir() {
		return nil
	}

	// skip template files if we are only generating statics.
	if isTemplate && gen.Options.StaticOnly == true {
		return nil
	}

	outPath := strings.Replace(inPath, gen.TemplateFiles, gen.Destination, 1)
	if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
		return err
	}

	if isTemplate {
		outPath = strings.Replace(outPath, filepath.Ext(outPath), "", 1)
		fmt.Printf("Processing Template File %s\n", inPath)
		// fmt.Printf("Generating Template: %s\n", inPath)

		templateFile, err := template.New(file.Name()).Funcs(gen.TemplateHelpers).ParseFiles(inPath)
		if err != nil {
			return err
		}

		generated, err := os.Create(outPath)
		if err != nil {
			return err
		}

		fmt.Printf("Writing To: %s: \n", outPath)
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
		fmt.Printf("%s: %q\n", q.Name, res)
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
	case "constant":
		answer = q.Default
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

// run all commands found in config.yaml run_after prop.
func (gen *Generator) runAfter() (err error) {
	for _, cmd := range gen.Config.RunAfter {
		var command bytes.Buffer

		cmdTemplate, err := template.New("cmd").Funcs(gen.TemplateHelpers).Parse(cmd)
		if err != nil {
			return err
		}
		if err := cmdTemplate.Execute(&command, gen.Answers); err != nil {
			return err
		}

		fmt.Printf("RunningCommand: %s \n", command.String())

		split := strings.Split(command.String(), " ")
		name := split[0]
		arguments := split[1:len(split)]

		// execute and break on error.
		if err := execute(name, arguments...); err != nil {
			return err
		}
	}
	return err
}
