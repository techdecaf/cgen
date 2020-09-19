package app

import (
	"os"
)

// Confirm prompt
func Confirm(msg string) {
	cgen := &CGen{}
	confirm, err := cgen.Generator.Ask(Question{
		Name:    "Confirm",
		Type:    "bool",
		Prompt:  msg,
		Default: "false",
	})

	if err != nil {
		Log.Fatal("confirmation_prompt", err)
	}

	if confirm != "true" {
		os.Exit(0)
	}
}
