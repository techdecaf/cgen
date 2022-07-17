package app

import (
	"encoding/json"
	"fmt"
)

// GeneratorParams struct
type GeneratorParams struct {
	ProjectName      string `json:"ProjectName"`      // name of this project
	TemplateName     string `json:"TemplateName"`     // selected cgen template
	ProjectDirectory string `json:"ProjectDirectory"` // destination directory for generated files
	// options
	PerformUpgrade bool `json:"PerformUpgrade"` // perform upgrade
	PromoteFile    bool `json:"PromoteFile"`    // run file promotion mode
	StaticOnly     bool `json:"StaticOnly"`     // only copy static files, no template interpolation
	Verbose        bool `json:"Verbose"`        // use verbose logging
	// IsExtension    bool   `json:"IsExtension"`    // template is running super exec
}

func (params *GeneratorParams) toJSON() error {
	out, err := json.Marshal(params)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
