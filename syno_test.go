package syno_test // import jdel.org/go-syno/syno_test

import (
	"os"
	"path/filepath"

	syno "jdel.org/go-syno"
)

var testRoot string
var o = syno.GetOptions()

func init() {
	testRoot = "tests"
	o.PackagesDir = filepath.Join(testRoot, "packages")
	o.CacheDir = filepath.Join(testRoot, "cache")
	o.ModelsFile = filepath.Join(testRoot, "cache/models.yml")
	o.Language = "enu"
	o.MD5 = true
	os.MkdirAll(o.CacheDir, 0755)
}
