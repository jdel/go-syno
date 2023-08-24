package syno // import jdel.org/go-syno/syno

import (
	"path/filepath"
)

// Options holds external context
type Options struct {
	PackagesDir string
	CacheDir    string
	ModelsFile  string
	Language    string
	MD5         bool
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
