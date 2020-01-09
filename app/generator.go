package app

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/mholt/archiver/v3"
	"github.com/techdecaf/templates"
	"github.com/techdecaf/utils"
	yaml "gopkg.in/yaml.v2"
)

// Output struct
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
	TemplateHelpers templates.Functions
	Variables       templates.Variables
	Config          *Config                `json:"Config"`
	Answers         map[string]interface{} `json:"Answers"`
	Options         struct {
		StaticOnly     bool `json:"StaticOnly"`
		PerformUpgrade bool `json:"PerformUpgrade"`
		PromoteFile    bool `json:"PromoteFile"`
		Verbose        bool `json:"Verbose"`
	}
}

// Init a new instance of Generator
func (gen *Generator) Init(params GeneratorParams) error {
	// set options
	gen.Options.StaticOnly = params.StaticOnly
	gen.Options.PerformUpgrade = params.PerformUpgrade
	gen.Options.PromoteFile = params.PromoteFile
	gen.Options.Verbose = params.Verbose

	if gen.Options.Verbose {
		params.toJSON()
	}

	// variables applied in this order
	// 1. cli options
	// 2. answer file
	// 3. environment variables
	// 4. user prompt

	// todo: validate inputs, that files exist etc
	// default destination to current working directory or use project name
	// check to see if an answers file exists in current dir
	answerFile := filepath.Join(params.Destination, ".cgen.yaml")
	if gen.Options.PerformUpgrade || gen.Options.PromoteFile {

		if gen.Options.PerformUpgrade {
			Log.Info("init", "running in Upgrade mode")
		}

		if gen.Options.PromoteFile {
			Log.Info("init", "running in PromoteFile mode")
		}
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

		gen.TemplateName = update.Template

		for k, v := range update.Answers {
			gen.AppendAnswer(k, fmt.Sprintf("%v", v))
			// handle `Name` as a special case, and set params from answer file.
			if k == "Name" {
				params.Name = fmt.Sprintf("%v", v)
			}
		}
	} else {
		gen.TemplateName = params.Tempate
	}

	// path to generators
	gen.Name = params.Name
	gen.TemplatesDir = params.TemplatesDir
	gen.Source = filepath.Join(gen.TemplatesDir, gen.TemplateName)

	gen.Destination = params.Destination
	gen.AnswersFile = filepath.Join(gen.Destination, ".cgen.yaml")
	gen.QuestionsFile = filepath.Join(gen.Source, "config.yaml")
	gen.TemplateFiles = filepath.Join(gen.Source, "template")

	gen.Config = &Config{}
	gen.Variables.Init()
	gen.TemplateHelpers = gen.LoadHelpers()

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

	gen.AppendAnswer("Name", gen.Name)
	gen.AppendAnswer("TemplateVersion", gen.Config.Version)
	gen.AppendAnswer("Timestamp", time.Now().UTC().Format(time.RFC3339))

	gen.Variables.Set(templates.Variable{
		Key:         "PWD",
		Value:       gen.Destination,
		OverrideEnv: true,
	})

	return nil
}

// Pull from remote repository
func (gen *Generator) Pull() error {
  Log.Info("pull", fmt.Sprintf("performing git pull in: %s", gen.Source))
  GitPull := templates.CommandOptions{
    Cmd:       "git pull",
    Dir: gen.Source,
		UseStdOut: true,
	}
  _, err := templates.Run(GitPull)
  return err
}

// Exec run the generator
func (gen *Generator) Exec() error {
	if err := gen.Prompt(); err != nil {
		return err
	}

	if err := filepath.Walk(gen.TemplateFiles, gen.WalkFiles); err != nil {
		return err
	}

	// save output to projct_root/.cgen.yaml
	ans, err := gen.Save()
	if err := ioutil.WriteFile(gen.Destination+"/.cgen.yaml", ans, 0644); err != nil {
		return err
	}

	if !gen.Options.PerformUpgrade {
		// run scripts in config.run_after array.
		if err := gen.RunAfter(); err != nil {
			return err
		}
	}
	return err
}

// Copy from src to dest
func (gen *Generator) Copy(src, dst string) error {
	Log.Info("copy", fmt.Sprintf("reading: %s", src))

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

	Log.Info("copy", fmt.Sprintf("writing: %s", dst))
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

type AppConfigProperties map[string]string

func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	config := AppConfigProperties{}

	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}

// WalkFiles files as part of the generator
func (gen *Generator) WalkFiles(inPath string, file os.FileInfo, err error) error {
	// identify template files
	isTemplate := filepath.Ext(inPath) == ".tmpl"
	isPtr := filepath.Ext(inPath) == ".ptr"

	// skip all directories
	if file.IsDir() {
		return nil
	}

	// skip template files if we are only generating statics.
	if isTemplate && gen.Options.StaticOnly == true {
		return nil
	}

	Log.Info("walk_files", fmt.Sprintf("source %s", inPath))

	outPath := strings.Replace(inPath, gen.TemplateFiles, gen.Destination, 1)

	if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
		return err
	}

	if isTemplate {
		Log.Info("walk_files", fmt.Sprintf("expanding template %s", outPath))
		outPath = strings.Replace(outPath, filepath.Ext(outPath), "", 1)

		generated, err := templates.ExpandFile(inPath, gen.TemplateHelpers)
		if err != nil {
			return err
		}

		Log.Info("walk_files", fmt.Sprintf("writing %s", outPath))
		return utils.WriteFile(outPath, generated)
	}

	if isPtr {
		var cmd *exec.Cmd
		Log.Info("walk_files", fmt.Sprintf("expanding pointer %s", outPath))
		gen.Copy(inPath, outPath)
		props, err := ReadPropertiesFile(outPath)
		if err != nil {
			return err
		}
		baseName := strings.Replace(outPath, filepath.Ext(outPath), "", 1)
		filename := filepath.Base(baseName)
		gitCmd := fmt.Sprintf("git archive --remote=%s HEAD:%s %s -o %s.tar", props["repository"], props["path"], filename, baseName)

		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", gitCmd)
		} else {
			cmd = exec.Command("bash", "-c", gitCmd)
		}
		if _, err = cmd.CombinedOutput(); err != nil {
			return err
		}

		if err = archiver.Unarchive(fmt.Sprintf("%s.tar", baseName), fmt.Sprintf("%s-tmp", baseName)); err != nil {
			return err
		}

		if err = os.Rename(fmt.Sprintf("%s-tmp/%s", baseName, filename), baseName); err != nil {
			return err
		}

		if err = os.Remove(outPath); err != nil {
			return err
		}
		if err = os.Remove(fmt.Sprintf("%s-tmp", baseName)); err != nil {
			return err
		}
		if err = os.Remove(fmt.Sprintf("%s.tar", baseName)); err != nil {
			return err
		}

		return nil
	}

	return gen.Copy(inPath, outPath)
}

// Prompt user to respond in the console.
func (gen *Generator) Prompt() error {
	for _, q := range gen.Config.Questions {
		res, err := gen.Ask(*q)
		Log.Info("prompt", fmt.Sprintf("%s: %q\n", q.Name, res))
		if err != nil {
			return err
		}
	}

	return nil
}

// Ask a question
func (gen *Generator) Ask(q Question) (answer string, err error) {

	if val := os.Getenv(q.Name); val != "" {
		return gen.AppendAnswer(q.Name, val), nil
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
			answer = ""
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

	return gen.AppendAnswer(q.Name, answer), nil
}

// AppendAnswer to gen.Answers map
func (gen *Generator) AppendAnswer(key, val string) (answer string) {
	// Log.Info("answer_file", fmt.Sprintf("%v: %v", key, val))
	if gen.Answers == nil {
		gen.Answers = make(map[string]interface{})
	}
	// append answer to the answers map.
	switch val {
	case "true":
		gen.Answers[key] = "true"
	case "false":
		gen.Answers[key] = ""
	default:
		gen.Answers[key] = val
	}

	gen.Variables.Set(templates.Variable{
		Key:         key,
		Value:       val,
		OverrideEnv: true,
	})

	return val
}

// Save yaml output
func (gen *Generator) Save() (out []byte, err error) {
	output := Output{}
	output.Answers = gen.Answers
	output.Template = gen.TemplateName

	res, err := yaml.Marshal(output)
	fmt.Println(string(res))
	return res, err
}

// RunAfter runs all commands found in config.yaml run_after prop.
func (gen *Generator) RunAfter() (err error) {
	for _, cmd := range gen.Config.RunAfter {
		var command string

		if command, err = templates.Expand(cmd, gen.TemplateHelpers); err != nil {
			return err
		}

		Log.Info("run_after", command)

		Command := templates.CommandOptions{
			Cmd:       command,
			Dir:       gen.Destination,
			UseStdOut: true,
		}

		if _, err := templates.Run(Command); err != nil {
			return err
		}
	}
	return err
}

// LoadHelpers adds additional helper functions
func (gen *Generator) LoadHelpers() templates.Functions {
	helpers := &gen.Variables.Functions

	// Custom Generator Functions
	helpers.Add("MkdirAll", func(dir string) string {
		path := filepath.Join(gen.Destination, dir)
		if err := os.MkdirAll(path, 0700); err != nil {
			Log.Info("error", fmt.Sprintf("%v", err))
		}
		return ""
	})

	helpers.Add("Touch", func(file string) string {
		path := filepath.Join(gen.Destination, file)
		if _, err := os.Create(path); err != nil {
			Log.Info("error", fmt.Sprintf("%v", err))
		}
		return ""
	})

	return *helpers
}
