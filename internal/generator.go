package app

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/mholt/archiver/v3"
	"github.com/techdecaf/templates"
	"github.com/techdecaf/utils"
)

// Question struct for questions file.
type Question struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Prompt  string   `yaml:"prompt"`
	Default string   `yaml:"default"`
	Options []string `yaml:"options,omitempty"`
}

// Generator struct
type Generator struct {
  Template         Template `json:"template"`
  Project          Project `json:"project"`
	Options          struct {
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
	if gen.Options.PerformUpgrade || gen.Options.PromoteFile {

		if gen.Options.PerformUpgrade {
			Log.Info("init", "running in upgrade mode")
		}

		if gen.Options.PromoteFile {
			Log.Info("init", "running in promote mode")
    }

    gen.Project = Project{
      Directory: params.ProjectDirectory,
    }

    // Initialize project
    if err := gen.Project.Init(); err != nil {
      return err
    }
    gen.Project.toJSON()


    // check to see if an tempate name was explicitly defined, else default to the .cgen answer file.
    // this prevents looping when running an upgrade that extends another template.
    if params.TemplateName == "" {
      params.TemplateName = gen.Project.State.Template
    }
    gen.Template = Template{
      Name: params.TemplateName,
    }

    // Initialize template
    if err := gen.Template.Init(); err != nil {
      return err
    }
    gen.Template.toJSON()

    return nil
  }


  gen.Template = Template{
    Name: params.TemplateName,
  }

  // Initialize template
  if err := gen.Template.Init(); err != nil {
    return err
  }
  gen.Template.toJSON()

  gen.Project = Project{
    Name: params.ProjectName,
    Directory: params.ProjectDirectory,
    State: ProjectState {
      Template: gen.Template.Name,
      Version: gen.Template.LatestStableVersion,
    },
  }

  // Initialize project
  if err := gen.Project.Init(); err != nil {
		return err
  }
  gen.Project.toJSON()

	gen.Project.AppendAnswer("Name", gen.Project.Name)
	// gen.Project.AppendAnswer("TemplateVersion", gen.Template.TemplateVersion)
  // gen.Project.AppendAnswer("Timestamp", time.Now().UTC().Format(time.RFC3339))

  gen.Project.variables.Set(templates.Variable{
		Key:         "PWD",
		Value:       gen.Project.Directory,
		OverrideEnv: true,
  })

	return nil
}

// Extends checks to see if this template extends another template and generates that first.
func (gen *Generator) Extends() error {
  var super = &CGen{}

  Log.Info("extends", gen.Template.Extends)

  // check to see if we need to initialize a super template
  if gen.Template.Extends == "" {
    Log.Info("extends", "this generator does not extend anything, skipping")
    return nil
  }

  // parse template name from git URL
  regex := regexp.MustCompile(`.*\/(.*)$`)
  matches := regex.FindStringSubmatch(gen.Template.Extends)
  if len(matches) != 2 {
    Log.Fatal(fmt.Sprintf("failed to identify template name from extends: %s", gen.Template.Extends), nil)
  }
  templateName := strings.ReplaceAll(matches[1], `.git`, "")
  Log.Info("extends", fmt.Sprintf("running super for %s", templateName))

  // install super generator
  if err := super.Init(); err != nil {
    Log.Fatal("super.Init() failed", err)
  }

  if _, err := super.Install(gen.Template.Extends); err != nil {
    Log.Fatal("super.Install failed", err)
  }

   // generate super template
    superParams := GeneratorParams{
			ProjectName:        gen.Project.Name,           // name of this projectTemplatesDirectory:   TemplatesDirectory,           // directory of all cgen templates
			TemplateName:        templateName,               // selected cgen template parsed from the git URL
			ProjectDirectory:   gen.Project.Directory,      // destination directory for generated files
			PerformUpgrade:     gen.Options.PerformUpgrade, // perform upgrade
			StaticOnly:         gen.Options.StaticOnly,     // only copy static files, no template interpolation
			Verbose:            gen.Options.Verbose,        // use verbose logging
    }

		if err := super.Generator.Init(superParams); err != nil {
      Log.Fatal("super_generator_init", err)
		}

		if err := super.Generator.Exec(); err != nil {
			Log.Fatal("super_generator_exec", err)
		}

  return nil
}

// Exec run the generator
func (gen *Generator) Exec() (err error) {
	if err := gen.Prompt(); err != nil {
		return err
  }

  if err := gen.Extends(); err != nil {
		return err
  }

  if !gen.Options.PerformUpgrade {
		// run scripts in config.run_after array.
		if err := gen.RunBefore(); err != nil {
			return err
		}
  }

	if err := filepath.Walk(gen.Template.Files, gen.WalkFiles); err != nil {
		return err
	}

  if err := gen.Project.SaveState(); err != nil {
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

// Pull from remote repository
func (gen *Generator) Pull() error {
  Log.Info("pull", fmt.Sprintf("performing git pull in: %s", gen.Template.Directory))
  GitPull := templates.CommandOptions{
    Cmd:       "git pull",
    Dir: gen.Template.Directory.ToString(),
		UseStdOut: true,
	}
  _, err := templates.Run(GitPull)
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

// ConfigProperties forgot
type ConfigProperties map[string]string

// ReadPropertiesFile reads props
func ReadPropertiesFile(filename string) (ConfigProperties, error) {
	config := ConfigProperties{}

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

	outPath := strings.Replace(inPath, gen.Template.Files, gen.Project.Directory, 1)

	if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
		return err
	}

	if isTemplate {
		Log.Info("walk_files", fmt.Sprintf("expanding template %s", outPath))
		outPath = strings.Replace(outPath, filepath.Ext(outPath), "", 1)

		generated, err := templates.ExpandFile(inPath, gen.Project.helpers)
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
	for _, q := range gen.Template.Questions {
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
		return gen.Project.AppendAnswer(q.Name, val), nil
	}

	if val := gen.Project.State.Answers[q.Name]; val != nil {
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

	return gen.Project.AppendAnswer(q.Name, answer), nil
}

// RunAfter runs all commands found in config.yaml run_after prop.
func (gen *Generator) RunAfter() (err error) {
	for _, cmd := range gen.Template.RunAfter {
		var command string

		if command, err = templates.Expand(cmd, gen.Project.helpers); err != nil {
			return err
		}

		Log.Info("run_after", command)

		Command := templates.CommandOptions{
			Cmd:       command,
			Dir:       gen.Project.Directory,
			UseStdOut: true,
		}

		if _, err := templates.Run(Command); err != nil {
			return err
		}
	}
	return err
}

// RunBefore runs all commands found in config.yaml run_before prop.
func (gen *Generator) RunBefore() (err error) {
	for _, cmd := range gen.Template.RunBefore {
		var command string

		if command, err = templates.Expand(cmd, gen.Project.helpers); err != nil {
			return err
		}

		Log.Info("run_before", command)

		Command := templates.CommandOptions{
			Cmd:       command,
			Dir:       gen.Project.Directory,
			UseStdOut: true,
		}

		if _, err := templates.Run(Command); err != nil {
			return err
		}
	}
	return err
}
