package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/techdecaf/templates"
	"gopkg.in/yaml.v2"
)

// ProjectState maintains last generated state by project
type ProjectState struct {
	Name      string                 `json:"name"`
	Template  string                 `yaml:"template"`
	Version   string                 `yaml:"version"`
	Timestamp string                 `yaml:"timestamp"`
	Answers   map[string]interface{} `yaml:"answers"`
}

// Project struct
type Project struct {
	Name      string
	Directory string
	State     ProjectState
	// private properties
	helpers   templates.Functions
	variables templates.Variables
	StateFile string
}

// Init a new instance of Generator
func (proj *Project) Init() error {
	state := Project{}.State

	if proj.Directory == "" {
		return fmt.Errorf("can not init a new project without a directory")
	}

	proj.StateFile = filepath.Join(proj.Directory, ".cgen.yaml")

	// if stateFile exists read it in and load saved values
	if _, err := os.Stat(proj.StateFile); err == nil {
		readState, err := os.Open(proj.StateFile)
		if err != nil {
			return err
		}

		defer readState.Close()
		byteValue, _ := ioutil.ReadAll(readState)
		yaml.Unmarshal(byteValue, &state)
		// load new state into current state
		proj.State = state

		for k, v := range proj.State.Answers {
			proj.AppendAnswer(k, fmt.Sprintf("%v", v))
		}
	}

	if proj.State.Name != "" {
		proj.Name = proj.State.Name
	} else {
		proj.State.Name = proj.Name
	}

	proj.variables.Init()
	proj.LoadHelpers()

	return nil
}

// AppendAnswer to proj.State.Answers map
func (proj *Project) AppendAnswer(key, val string) (answer string) {
	if proj.State.Answers == nil {
		proj.State.Answers = make(map[string]interface{})
	}
	// append answer to the answers map.
	switch val {
	case "true":
		proj.State.Answers[key] = "true"
	case "false":
		proj.State.Answers[key] = ""
	default:
		proj.State.Answers[key] = val
	}

	proj.variables.Set(templates.Variable{
		Key:         key,
		Value:       val,
		OverrideEnv: true,
	})

	return val
}

// SaveState yaml output
func (proj *Project) SaveState() (err error) {
	state, err := yaml.Marshal(&proj.State)
	if err != nil {
		return err
	}

	fmt.Println(string(state))

	if err := ioutil.WriteFile(proj.StateFile, state, 0644); err != nil {
		return err
	}

	return err
}

// LoadHelpers adds additional helper functions
func (proj *Project) LoadHelpers() {
	helpers := &proj.variables.Functions

	// Custom Generator Functions
	helpers.Add("MkdirAll", func(dir string) string {
		path := filepath.Join(proj.Directory, dir)
		if err := os.MkdirAll(path, 0700); err != nil {
			Log.Info("error", fmt.Sprintf("%v", err))
		}
		return ""
	})

	helpers.Add("Touch", func(file string) string {
		path := filepath.Join(proj.Directory, file)
		if _, err := os.Create(path); err != nil {
			Log.Info("error", fmt.Sprintf("%v", err))
		}
		return ""
	})
	proj.helpers = *helpers
}

func (proj *Project) toJSON() error {
	out, err := json.Marshal(proj)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
