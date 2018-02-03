package syno

import (
	"os"
	"testing"
)

var defaultModels Models

func init() {
	defaultModels = Models{
		&Model{
			Name: "Spaghetti",
		},
		&Model{
			Name: "Penne",
		},
		&Model{
			Name: "Farfalle",
		},
		&Model{
			Name: "Papardelle",
		},
	}
}

func cleanupModels(t *testing.T) {
	if err := os.RemoveAll(o.ModelsFile); err != nil {
		t.Skipf("Cannot cleanup %s", o.ModelsFile)
	}
}

func TestModelsGetFromFile(t *testing.T) {
	cleanupModels(t)
	if err := defaultModels.SaveModelsFile(); err != nil {
		t.Error(err)
	}
	if m, err := getModelsFromModelsFile(); err != nil {
		t.Error(err)
	} else if len(m) != len(defaultModels) {
		t.Errorf("Expected to read %d models from file but fot %d", len(defaultModels), len(m))
	}
}

func TestModelsGetFromFileReadError(t *testing.T) {
	cleanupModels(t)
	if err := defaultModels.SaveModelsFile(); err != nil {
		t.Error(err)
	}
	if err := os.RemoveAll(o.ModelsFile); err != nil {
		t.Skipf("Could not create file %s to prepare for test", o.ModelsFile)
	}
	if _, err := getModelsFromModelsFile(); err == nil {
		t.Error("Expected a read error")
	}
}

// func TestModelsGetFromFileBadYML(t *testing.T) {
// 	cleanupModels(t)
// 	f, err := os.OpenFile(o.ModelsFile, os.O_CREATE|os.O_WRONLY, 0755)
// 	defer f.Close()
// 	if err != nil {
// 		t.Skipf("Could not write file %s to prepare for test", o.ModelsFile)
// 	}
// 	f.Write([]byte(":"))
// 	f.Close()
// 	m, err := getModelsFromModelsFile()
// 	fmt.Println(m)
// 	if err == nil {
// 		t.Error("Expected a parsing error")
// 	}
// }

func TestModelsGetFromInternet(t *testing.T) {
	if m, err := getModelsFromInternet(); err != nil {
		t.Error(err)
	} else if len(m) == 0 {
		t.Error("Expected at least one model from the internet but got 0")
	}
}

func TestModelsFinished(t *testing.T) {
	cleanupModels(t)
}
