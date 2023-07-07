package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	repos map[string]Repository
}

type Repository struct {
	path string
	default_branch *string
}

func (manifest Manifest) paths() (p []string) {
	for _, repo := range manifest.repos {
		p = append(p, repo.path)
	}

	return
}

func readManifest(file string) (*Manifest, error) {

	f, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	var manifest Manifest

	if err := yaml.Unmarshal(f, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func writeManifest(manifest Manifest) error {

	yamltext, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	err = os.WriteFile(*GlobalOptions.ManifestPath, yamltext, 0o644)

	return nil
}
