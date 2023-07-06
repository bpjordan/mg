package main

import "os"
import "gopkg.in/yaml.v3"

type manifest struct {
	repos []struct {
		path           string
		default_branch string
	}
}

func readManifest(file string) (*manifest, error) {

	f, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	var manifest manifest

	if err := yaml.Unmarshal(f, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func writeManifest(manifest manifest) error {

	return nil
}
