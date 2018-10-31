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
	"strings"
	"time"

	"github.com/manifoldco/promptui"
)

// Question struct for questions file.
type Question struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Prompt  string   `json:"prompt"`
	Default string   `json:"default"`
	Options []string `json:"options,omitempty"`
}

// Generator struct
type Generator struct {
	Name          string
	Source        string
	Destination   string
	QuestionsFile string
	TemplateFiles string
	Questions     []*Question
	Answers       map[string]string
}

func (gen *Generator) init(name, src string) error {
	// todo: validate inputs, that files exist etc
	// default destination to current working directory or use project name
	gen.Source = src
	gen.Name = name
	gen.Destination = path.Join(".", gen.Name)
	gen.QuestionsFile = path.Join(gen.Source, "questions.json")
	gen.TemplateFiles = path.Join(gen.Source, "template")

	gen.Answers["TimeStamp"] = time.Now().UTC().Format(time.RFC3339)

	// check for required project structure
	if _, err := os.Stat(gen.TemplateFiles); os.IsNotExist(err) {
		return fmt.Errorf("%s does not have the required template directory, please check the README file", gen.Source)
	}

	if _, err := os.Stat(gen.QuestionsFile); os.IsNotExist(err) {
		log.Printf("%s does not have a questions.json file, so it may not actually be a cgen template...", gen.Source)
	}

	return nil
}

func (gen *Generator) exec() error {
	if err := gen.prompt(); err != nil {
		return err
	}
	err := filepath.Walk(gen.TemplateFiles, gen.walkFiles)

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
		var templateFile = template.Must(template.ParseFiles(inPath))

		generated, err := os.Create(outPath)
		if err != nil {
			return err
		}

		if err := templateFile.Execute(generated, gen.Answers); err != nil {
			return err
		}
	} else {
		gen.copy(inPath, outPath)
	}
	// if the file name starts with an _ then parse it as a template.
	// else read the file as is and spit it out
	// todo: how would we update the template later?
	return nil // no errors
}

func (gen *Generator) prompt() error {
	var questions = []*Question{}

	jsonFile, err := os.Open(gen.QuestionsFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &questions)

	for _, q := range questions {
		res, err := gen.ask(*q)
		if err != nil {
			return err
		}
		fmt.Printf("You choose %q\n", res)
	}

	return nil
}

func (gen *Generator) ask(q Question) (answer string, err error) {
	// init answers
	if gen.Answers == nil {
		gen.Answers = make(map[string]string)
	}

	switch q.Type {
	case "string":
		prompt := promptui.Prompt{
			Label:   q.Prompt,
			Default: q.Default,
		}
		answer, err = prompt.Run()
	case "bool":
		prompt := promptui.Prompt{
			Label:     q.Prompt,
			IsConfirm: true,
			Default:   "n",
		}
		answer, err = prompt.Run()
		if answer == "y" {
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

	// append answer to the answers map.
	gen.Answers[q.Name] = answer
	return answer, nil
}
