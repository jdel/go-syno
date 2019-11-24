package syno_test // import jdel.org/go-syno/syno_test

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	syno "jdel.org/go-syno"
)

func TestOptionsGetOptions(t *testing.T) {
	newOptions := syno.GetOptions()
	if newOptions.PackagesDir != filepath.Join("tests", "packages") {
		t.Errorf("Expected packages dir to be %s but got %s", filepath.Join("tests", "packages"), newOptions.PackagesDir)
	}
	if newOptions.CacheDir != filepath.Join("tests", "cache") {
		t.Errorf("Expected cache dir to be %s but got %s", filepath.Join("tests", "cache"), newOptions.CacheDir)
	}
	if newOptions.ModelsFile != filepath.Join("tests", "cache", "models.yml") {
		t.Errorf("Expected models file to be %s but got %s", filepath.Join("tests", "models.yml"), newOptions.ModelsFile)
	}
	if newOptions.Language != "enu" {
		t.Errorf("Expected language to be %s but got %s", "enu", newOptions.Language)
	}
	if !newOptions.MD5 {
		t.Errorf("Expected MD5 to be %t but got %t", true, newOptions.MD5)
	}
}

func TestOptionsLogLevel(t *testing.T) {
	cleanupModels(t)
	var b bytes.Buffer
	syno.SetLogLevel(syno.LogDebugLevel)
	syno.SetLogOutput(&b)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	syno.GetModels(false)
	if len(b.String()) == 0 {
		t.Errorf("Expecting loglevel Debug to write to buffer but it didn't")
	}
}

func TestOptionsSetOptions(t *testing.T) {
	previousOptions := *o
	defer func() {
		syno.SetOptions(previousOptions)
	}()

	syno.SetOptions(syno.Options{
		CacheDir:    "c",
		PackagesDir: "p",
		ModelsFile:  "m",
		Language:    "l",
		MD5:         false,
	})

	if o.CacheDir != "c" {
		t.Errorf("Expected cache dir to be c but got %s", o.CacheDir)
	}
	if o.PackagesDir != "p" {
		t.Errorf("Expected packages dir to be  but got %s", o.PackagesDir)
	}
	if o.ModelsFile != "m" {
		t.Errorf("Expected models file to be m but got %s", o.ModelsFile)
	}
	if o.Language != "l" {
		t.Errorf("Expected language to be l but got %s", o.Language)
	}
	if o.MD5 {
		t.Errorf("Expected MD5 to be %t but got %t", false, o.MD5)
	}
}

func TestOptoinsFinished(t *testing.T) {
	cleanupModels(t)
}
