package packagemanager

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	packageDir   = os.Getenv("SOURCE_BUILD_DIR")
	manifestPath = os.Getenv("SOURCE_MANIFEST_PATH")
)

func ReadManifest() ([]Package, error) {
	fileContent, err := os.ReadFile(manifestPath)
	if err != nil {
		errMsg := fmt.Sprintf("error reading the manifest file: %v", err)
		return nil, errors.New(errMsg)
	}

	var packs []Package

	err = yaml.Unmarshal(fileContent, &packs)
	if err != nil {
		errMsg := fmt.Sprintf("error unmarshalling yaml file to struct: %v", err)
		return nil, errors.New(errMsg)
	}

	return packs, nil
}

func UpdateManifest(packs *[]Package) error {
	yamlContent, err := yaml.Marshal(packs)
	if err != nil {
		errMsg := fmt.Sprintf("error marshalling the yaml file: %v", err.Error())
		return errors.New(errMsg)
	}

	err = os.WriteFile(manifestPath, yamlContent, os.ModePerm)
	if err != nil {
		errMsg := fmt.Sprintf("error writing yaml file: %v", err.Error())
		return errors.New(errMsg)
	}

	return nil
}
