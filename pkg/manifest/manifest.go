package manifest

import (
	"os"

	"gopkg.in/yaml.v3"
)

var CommandManifest Manifest

type Manifest struct {
	Repos map[string]Repository `yaml:"repos"`
}

type Repository struct {
	Path string `yaml:"path"`
	DefaultBranch string `yaml:"default_branch,omitempty"`
}

func (manifest Manifest) Paths() (p []string) {
	for _, repo := range manifest.Repos {
		p = append(p, repo.Path)
	}

	return
}

func ReadManifest(file string) (*Manifest, error) {

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

func WriteManifest(manifest Manifest, file string) error {

	yamltext, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	err = os.WriteFile(file, yamltext, 0o644)

	return nil
}
