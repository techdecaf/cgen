package app

import (
	"testing"
)

func TestReadPropertiesFile(t *testing.T) {
	props, err := ReadPropertiesFile("./test-files/file.json.ptr")
	if err != nil {
		t.Error("Error while reading properties file")
	}

	if props["repository"] != "git@kochsource.io:KAES/generators/npm-module.git" || props["path"] != "" {
		t.Error("Error properties not loaded correctly")
	}

}
