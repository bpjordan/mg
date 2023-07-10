package manifest

import (
	"os"

	"gopkg.in/yaml.v3"
)

var CommandManifest Manifest

type Manifest struct {
	Repos []Repository `yaml:"repos"`
}

type Repository struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Home string `yaml:"home,omitempty"`
}

func (manifest Manifest) Paths() (p []string) {
	for _, repo := range manifest.Repos {
		p = append(p, repo.Path)
	}

	return
}

func (manifest Manifest) Names() (p []string) {
	for _, repo := range manifest.Repos {
		p = append(p, repo.Name)
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

// Allow repositories to be entered as just a single string
// with the path, or as a struct with all fields
func (r *Repository) UnmarshalYAML(value *yaml.Node) error {

	var path string
	if err := value.Decode(&path); err == nil {
		r.Name = path
		r.Path = path
		return nil
	}

	type rawRepoStruct Repository
	if err := value.Decode((*rawRepoStruct)(r)); err != nil {
		return err
	}

	return nil
}
