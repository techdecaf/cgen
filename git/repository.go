package git

import (
	"strings"

	"github.com/techdecaf/templates"
)

// Repository represents a directory containing a git repository
type Repository string

// ListTags for a new repository
func (repo Repository) ListTags() (Tags, error) {
	gitTags := templates.CommandOptions{
		Cmd:        "git --no-pager tag --sort='v:refname'",
		Dir:        string(repo),
		UseStdOut:  false,
		TrimOutput: true,
	}

	out, err := templates.Run(gitTags)
	if err != nil {
		return Tags{}, err
	}

	tags := strings.Split(strings.TrimSpace(out), "\n")
	return Tags(tags), nil
}

// Pull from remote repository
func (repo Repository) Pull() error {
	GitPull := templates.CommandOptions{
		Cmd:       "git pull --all",
		Dir:       string(repo),
		UseStdOut: true,
	}
	_, err := templates.Run(GitPull)
	return err
}

// ToString returns the string representation of the git repository
func (repo Repository) ToString() string {
	return string(repo)
}
