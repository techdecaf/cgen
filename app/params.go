package app

import (
	"encoding/json"
	"fmt"
)

// GeneratorParams struct
type GeneratorParams struct {
	Name           string `json:"Name"`           // name of this project
	TemplatesDir   string `json:"TemplatesDir"`   // directory of all cgen templates
	Tempate        string `json:"Tempate"`        // selected cgen template
	Pointer        string `json:"Pointer"`        // pointer file to different repository
	Destination    string `json:"Destination"`    // destination directory for generated files
	PerformUpgrade bool   `json:"PerformUpgrade"` // perform upgrade
	PromoteFile    bool   `json:"PromoteFile"`    // run file promotion mode
	StaticOnly     bool   `json:"StaticOnly"`     // only copy static files, no template interpolation
	Verbose        bool   `json:"Verbose"`        // use verbose logging
}

func (params *GeneratorParams) toJSON() error {
	json, err := json.Marshal(params)
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}
