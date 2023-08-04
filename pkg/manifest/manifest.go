package manifest

var CommandManifest Manifest

type Manifest struct {
	UnfilteredRepos []Repository `yaml:"repos" mapstructure:"repos"`
	Groups map[string][]string `yaml:"groups,omitempty" mapstructure:"groups"`
}

type Repository struct {
	Name string `yaml:"name" mapstructure:"name"`
	Path string `yaml:"path" mapstructure:"path"`
	Home string `yaml:"home,omitempty" mapstructure:"home"`
}

func (manifest Manifest) Repos() []Repository {
	return manifest.UnfilteredRepos
}

func (manifest Manifest) Paths() (p []string) {
	for _, repo := range manifest.UnfilteredRepos {
		p = append(p, repo.Path)
	}

	return
}

func (manifest Manifest) Names() (p []string) {
	for _, repo := range manifest.UnfilteredRepos {
		p = append(p, repo.Name)
	}

	return
}

func (r *Repository) UnmarshalText(text []byte) error {
	path := string(text)
	r.Name = path
	r.Path = path

	return nil
}

