package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// WriteStructToYAMLFile writes a struct to a YAML file.
// The file is created if it does not exist, and overwritten if it does.
func WriteStructToYAMLFile(filename string, data interface{}) error {

	yamlDataHeader := []byte("# PoeBuy config file\n# DO NOT EDIT\n\n")

	// Marshal the struct to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling to YAML: %v", err)
	}

	// Write the YAML data to a file
	err = os.WriteFile(filename, append(yamlDataHeader, yamlData...), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
