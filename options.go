package syno // import jdel.org/go-syno/syno

import (
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Options holds external context
type Options struct {
	PackagesDir string `json:"packages_dir,omitempty" yaml:"packages_dir,omitempty"`
	CacheDir    string `json:"cache_dir,omitempty" yaml:"cache_dir,omitempty"`
	ModelsFile  string `json:"models_file,omitempty" yaml:"models_file,omitempty"`
	Language    string `json:"language,omitempty" yaml:"language,omitempty"`
	MD5         bool   `json:"md5,omitempty" yaml:"md5,omitempty"`
}

func (o *Options) String() string {
	yamlOptions, err := yaml.Marshal(o)
	if err != nil {
		return ""
	}
	return string(yamlOptions)
}

var o Options

func init() {
	// Default options
	o = Options{
		PackagesDir: filepath.Join(executablePath(), "packages"),
		CacheDir:    filepath.Join(executablePath(), "cache"),
		ModelsFile:  filepath.Join(executablePath(), "models.yml"),
		Language:    "enu",
	}
}

// SetOptions sets global options
func SetOptions(opt Options) {
	o = opt
}

// GetOptions returns global options
func GetOptions() *Options {
	return &o
}
