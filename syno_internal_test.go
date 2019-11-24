package syno // import jdel.org/go-syno/syno

import (
	"os"
	"path/filepath"
)

var testRoot string

func init() {
	testRoot = "tests"
	o.PackagesDir = filepath.Join(testRoot, "packages")
	o.CacheDir = filepath.Join(testRoot, "cache")
	o.ModelsFile = filepath.Join(testRoot, "cache/models.yml")
	os.MkdirAll(o.CacheDir, 0755)
}
