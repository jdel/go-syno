package syno_test

import (
	"os"
	"path/filepath"
	"testing"

	syno "github.com/jdel/go-syno"
	log "github.com/sirupsen/logrus"
)

var defaultModels syno.Models

func init() {
	defaultModels = syno.Models{
		&syno.Model{
			Name: "Spaghetti",
		},
		&syno.Model{
			Name: "Penne",
		},
		&syno.Model{
			Name: "Farfalle",
		},
		&syno.Model{
			Name: "Papardelle",
		},
	}
}

func cleanupModels(t *testing.T) {
	if err := os.RemoveAll(o.ModelsFile); err != nil {
		t.Skipf("Cannot cleanup %s", o.ModelsFile)
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func TestModelsFileExists(t *testing.T) {
	cleanupModels(t)
	syno.SetLogLevel(log.ErrorLevel)
	if fileExists(o.ModelsFile) {
		t.Errorf("Expected file %s to not exist on disk", o.ModelsFile)
	}
	if err := defaultModels.SaveModelsFile(); err != nil {
		t.Error(err)
	}
	if !fileExists(o.ModelsFile) {
		t.Errorf("Expected file %s to not exist on disk", o.ModelsFile)
	}
}

func TestGetModelsFromFile(t *testing.T) {
	cleanupModels(t)
	if err := defaultModels.SaveModelsFile(); err != nil {
		t.Skip("Could not create models file to prepare for test")
	}
	if m, err := syno.GetModels(false); err != nil {
		t.Error(err)
	} else if len(m) != len(defaultModels) {
		t.Errorf("Expected %d models but got %d", len(defaultModels), len(m))
	}
}

func TestGetModelsFromInternet(t *testing.T) {
	cleanupModels(t)
	if err := os.RemoveAll(o.ModelsFile); err != nil {
		t.Skip("Could not remove models file to prepare for test")
	}
	if m, err := syno.GetModels(false); err != nil {
		t.Error(err)
	} else if len(m) == len(defaultModels) {
		t.Errorf("Expected more than %d models but got %d", len(defaultModels), len(m))
	}
}

func TestGetModelsFromInternetForce(t *testing.T) {
	cleanupModels(t)
	if err := defaultModels.SaveModelsFile(); err != nil {
		t.Skip("Could not create models file to prepare for test")
	}
	if m, err := syno.GetModels(true); err != nil {
		t.Error(err)
	} else if len(m) == len(defaultModels) {
		t.Errorf("Expected more than %d models but got %d", len(defaultModels), len(m))
	}
}

func TestModelsFilterByName(t *testing.T) {
	if m := defaultModels.FilterByName("Pa"); len(m) != 2 {
		t.Errorf("Expected 2 models but got %d", len(m))
	}
	if m := defaultModels.FilterByName("invalid"); len(m) != 0 {
		t.Errorf("Expected 0 models but got %d", len(m))
	}
}

func TestSaveModelsFile(t *testing.T) {
	cleanupModels(t)
	if err := defaultModels.SaveModelsFile(); err != nil {
		t.Error(err)
	}
	if !fileExists(o.ModelsFile) {
		t.Errorf("Expected file %s to exist on disk", o.ModelsFile)
	}
	if m, err := syno.GetModels(false); err != nil {
		t.Error(err)
	} else if len(m) != len(defaultModels) {
		t.Errorf("Expected to read %d models from file but fot %d", len(defaultModels), len(m))
	}
}

func TestSaveModelsFileWriteError(t *testing.T) {
	cleanupModels(t)
	previousModelsFile := o.ModelsFile
	defer func() {
		o.ModelsFile = previousModelsFile
	}()
	o.ModelsFile = filepath.Join(testRoot, "badcache", "models.yml")
	if err := os.MkdirAll(o.CacheDir, 0666); err != nil {
		t.Skipf("Cannot create directory %s to prepare for test", o.CacheDir)
	}
	if err := defaultModels.SaveModelsFile(); err == nil {
		t.Error("Expected a write error")
	}
}

func TestModelsFinished(t *testing.T) {
	cleanupModels(t)
}
