package syno // import jdel.org/go-syno/syno

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var pX8664,
	pCedarview1,
	pCedarview2,
	pCedarview3,
	pX86,
	pEvansport,
	p1,
	p2 *Package

func init() {
	testRoot = "tests"
	o.CacheDir = filepath.Join(testRoot, "cache")
	o.PackagesDir = filepath.Join(testRoot, "packages")
	o.ModelsFile = filepath.Join(testRoot, "models.yml")

	p1 = &Package{
		Name:        "package1",
		DisplayName: "Bleu",
		Arch:        "noarch",
		Firmware:    "3.0",
		Beta:        true,
	}
	p2 = &Package{
		Name:        "package2",
		DisplayName: "French Emmental",
		Arch:        "noarch",
		Firmware:    "6.1",
	}
	pX8664 = &Package{
		Name:        "x86-64-package",
		DisplayName: "Swiss Emmental",
		Arch:        "x86_64",
	}
	pCedarview1 = &Package{
		Name:        "cedarview-package1",
		DisplayName: "Brie",
		Arch:        "cedarview",
		Beta:        true,
	}
	pCedarview2 = &Package{
		Name:        "cedarview-package",
		DisplayName: "Cedarview Package 2",
		Arch:        "cedarview",
		Version:     "1.0",
	}
	pCedarview3 = &Package{
		Name:        "cedarview-package",
		DisplayName: "Cedarview Package 3",
		Arch:        "cedarview",
		Version:     "1.1",
	}
	pX86 = &Package{
		Name: "x86-package",
		Arch: "x86",
	}
	pEvansport = &Package{
		Name: "evansport-package",
		Arch: "evansport",
		Beta: true,
	}
}

func skipIfMissingTestPackage(t *testing.T, f string) *Package {
	var p *Package
	var err error
	if p, err = NewPackage(f); err != nil {
		t.Skip(err)
	}
	return p
}

func cleanupPackage(t *testing.T, p string) {
	packageCacheDir := filepath.Join(o.CacheDir, p)
	if err := os.RemoveAll(packageCacheDir); err != nil {
		t.Skipf("Cannot cleanup %s", packageCacheDir)
	}
	if packageFile := filepath.Join(o.PackagesDir, p); fileExists(fmt.Sprintf("%s.ignored", packageFile)) {
		if err := os.Rename(fmt.Sprintf("%s.ignored", packageFile), packageFile); err != nil {
			t.Skipf("Cannot unignore %s", packageFile)
		}
	}
}

func TestPackageFullPath(t *testing.T) {
	p := skipIfMissingTestPackage(t, "real-package.spk")

	if expectedFullPath := filepath.Join(o.PackagesDir, p.fileName); p.FullPath() != expectedFullPath {
		t.Errorf("Expected FullPath %s but got %s", expectedFullPath, p.FullPath())
	}
}

func TestX8664PackageFamilyContains(t *testing.T) {
	if !pX8664.familyOrArchMatch("x86") {
		t.Errorf("Expected %s to contain x86", pX8664.Arch)
	}
}

func TestX8664PackageFamilyNotContains(t *testing.T) {
	if pX8664.familyOrArchMatch("fake-arch") {
		t.Errorf("Expected %s to not contain fake-arch", pX8664.Arch)
	}
}

func TestPackageIndex(t *testing.T) {
	pp := Packages{pCedarview2, pCedarview3}
	var i int
	var err error
	if i, err = pp.index("cedarview-package", "cedarview"); err != nil {
		t.Error("Package not in packages")

	}
	if p := pp[i]; p != pCedarview2 {
		t.Errorf("Expected %+v but got %+v", pCedarview2, p)
	}
}

func TestPackageContainsFiles(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	if _, err := p.containsFiles("INFO"); err != nil {
		t.Error("Expected real-package.spk to contain an INFO file")
	}
	if ok, _ := p.containsFiles("INFO"); !ok {
		t.Error("Expected real-package.spk to contain an INFO file")
	}
}

func TestBadPackageContainsFiles(t *testing.T) {
	cleanupPackage(t, "bad-package.spk")
	p := Package{}
	p.fileName = "bad-package.spk"
	if _, err := p.containsFiles("INFO"); err == nil {
		t.Error("Expected bad-package.spk to not contain an INFO file")
	}
}

func TestDeadPackageContainsFiles(t *testing.T) {
	cleanupPackage(t, "dead-package.spk")
	p := Package{}
	p.fileName = "dead-package.spk"
	if _, err := p.containsFiles("INFO"); err == nil {
		t.Error("Expected dead-package.spk to not exist on disk")
	}
}

func TestPackageExtractFiles(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	var files []string
	var err error

	if files, err = p.extractFiles("INFO"); err != nil {
		t.Errorf("Could not extract INFO form real-package.spk: %s", err)

	}
	if l := len(files); l != 1 {
		t.Errorf("Expected 1 file but got %d", l)
	}
	for _, f := range files {
		if fullPath := filepath.Join(o.CacheDir, p.fileName, f); !fileExists(fullPath) {
			t.Errorf("Expected file %s to have been written at %s", f, fullPath)
		}
	}

	if files, err = p.extractFiles("PACKAGE_ICON.*"); err != nil {
		t.Errorf("Could not extract PACKAGE_ICON.* form real-package.spk: %s", err)

	}
	if l := len(files); l != 2 {
		t.Errorf("Expected 2 file but got %d", l)
	}
	for _, f := range files {
		if fullPath := filepath.Join(o.CacheDir, p.fileName, f); !fileExists(fullPath) {
			t.Errorf("Expected file %s to have been written at %s", f, fullPath)
		}
	}

	p = skipIfMissingTestPackage(t, "dot-slash-package.spk")

	if files, err = p.extractFiles("INFO"); err != nil {
		t.Errorf("Could not extract INFO form real-package.spk: %s", err)

	}
	if l := len(files); l != 1 {
		t.Errorf("Expected 1 file but got %d", l)
	}
	for _, f := range files {
		if fullPath := filepath.Join(o.CacheDir, p.fileName, f); !fileExists(fullPath) {
			t.Errorf("Expected file %s to have been written at %s", f, fullPath)
		}
	}
}
func TestPackageExtractFilesWriteError(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	previousCache := o.CacheDir
	defer func() {
		o.CacheDir = previousCache
	}()
	o.CacheDir = filepath.Join(testRoot, "badcache")
	if err := os.MkdirAll(o.CacheDir, 0666); err != nil {
		t.Skipf("Cannot create directory %s to prepare for test", o.CacheDir)
	}
	if _, err := p.extractFiles("INFO"); err == nil {
		t.Error("Expected a write error")
	}
}

func TestPackageGetSize(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	if s, err := p.getSize(); err != nil {
		t.Error(err)
	} else if s != "168960" {
		t.Errorf("Expected size of 168960 bytes but got %s", s)
	}
}

func TestPackageGetSizeReadError(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	os.Rename(p.FullPath(), fmt.Sprintf("%s.ignored", p.FullPath()))
	if _, err := p.getSize(); err == nil {
		t.Error("Expected a file read error")
	}
}

func TestPackageGetMD5(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	if m, err := p.getMD5(); err != nil {
		t.Error(err)
	} else if m != "364b86c18d3d849ccb40b8ae63f7fc53" {
		t.Errorf("Expected MD5 of 364b86c18d3d849ccb40b8ae63f7fc53 but got %s", m)
	}
}

func TestPackageGetMD5ReadError(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	os.Rename(p.FullPath(), fmt.Sprintf("%s.ignored", p.FullPath()))
	if _, err := p.getMD5(); err == nil {
		t.Error("Expected a file read error")
	}
}

func TestPackageGetNewlyExtractedInfo(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := Package{}
	p.fileName = "real-package.spk"
	if i, err := p.getOrExtractInfo(); err != nil {
		t.Error(err)
	} else if i != filepath.Join(o.CacheDir, "real-package.spk", "INFO") {
		t.Errorf("Expected INFO file to be extracted at %s but got %s", filepath.Join(o.CacheDir, "real-package.spk", "INFO"), i)
	} else if !fileExists(i) {
		t.Errorf("Expected INFO file to be newly extracted")
	}
}

func TestPackageGetAlreadyExtractedInfo(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	if err := os.MkdirAll(filepath.Join(o.CacheDir, "real-package.spk"), 0755); err != nil {
		t.Skip("Cannot create package cache directory to prepare test")
	}
	if _, err := os.Create(filepath.Join(o.CacheDir, "real-package.spk", "INFO")); err != nil {
		t.Skip("Cannot create INFO file to prepare test")
	}
	p := Package{}
	p.fileName = "real-package.spk"
	if i, err := p.getOrExtractInfo(); err != nil {
		t.Error(err)
	} else if i != filepath.Join(o.CacheDir, "real-package.spk", "INFO") {
		t.Errorf("Expected INFO file to be %s but got %s", filepath.Join(o.CacheDir, "real-package.spk", "INFO"), i)
		if _, err := os.Stat(i); os.IsNotExist(err) {
			t.Errorf("Expected INFO file to be on disk at %s but it doesn't exist", i)
		}
	} else if !fileExists(i) {
		t.Errorf("Expected INFO file to be already extracted")
	}
}

func TestBadPackageGetInfo(t *testing.T) {
	cleanupPackage(t, "bad-package.spk")
	p := Package{}
	p.fileName = "bad-package.spk"
	if i, err := p.getOrExtractInfo(); err == nil {
		t.Error(err)
	} else if i != "" {
		t.Errorf("Expected INFO to be empty but got %s", i)
	} else if fileExists(i) {
		t.Errorf("Expected INFO file not to exist on disk")
	}
}

func TestPackageExtractImages(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := Package{}
	p.fileName = "real-package.spk"
	i, err := p.extractImages()
	if err != nil {
		t.Error(err)
	}
	if len(i) != 2 {
		t.Errorf("Expected 2 image files extracted but got %d", len(i))
	}
	if !fileExists(filepath.Join(o.CacheDir, "real-package.spk", i[0])) {
		t.Errorf("Cannot find %s", filepath.Join(o.CacheDir, "real-package.spk", i[0]))
	}
	if !fileExists(filepath.Join(o.CacheDir, "real-package.spk", i[1])) {
		t.Errorf("Cannot find %s", filepath.Join(o.CacheDir, "real-package.spk", i[1]))
	}
}

func TestPackageExtractImagesDeadPackage(t *testing.T) {
	p := Package{}
	p.fileName = "dead-package.spk"
	if _, err := p.extractImages(); err == nil {
		t.Error("Expected a read error")
	}
}

func TestPackageGetImages(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	i, err := p.getImages()
	if err != nil {
		t.Error(err)
	} else if len(i) != 2 {
		t.Errorf("Expected 2 image files extracted but got %d", len(i))
	}
	if i1 := filepath.Join(o.CacheDir, p.fileName, i[0]); !fileExists(i1) {
		t.Errorf("Cannot find %s", i1)
	}
	if i2 := filepath.Join(o.CacheDir, p.fileName, i[1]); !fileExists(i2) {
		t.Errorf("Cannot find %s", i2)
	}
}

func TestParseBoolTrue(t *testing.T) {
	if b, _ := parseBoolOrYes("true"); b == false {
		t.Errorf("Expected true to be parsed as true but got false")
	}
}

func TestParseBoolYes(t *testing.T) {
	if b, _ := parseBoolOrYes("yes"); b == false {
		t.Errorf("Expected yes to be parsed as true but got false")
	}
}

func TestParseBoolInvalid(t *testing.T) {
	if b, _ := parseBoolOrYes("invalid"); b == true {
		t.Errorf("Expected invalid to be parsed as false but got true")
	}
}

func TestNewDebugPackage(t *testing.T) {
	if p := NewDebugPackage("Test"); p.Name == "" {
		t.Error("Expecting debug package to have a Name")
	} else if p.DisplayName == "" {
		t.Error("Expecting debug package to have a DisplayName")
	} else if p.Description != "Test" {
		t.Error("Expecting debug package to have a Description that reads Test")
	}
}

func TestPopulateImageFields(t *testing.T) {
	cleanupPackage(t, "real-package-screen.spk")
	p := skipIfMissingTestPackage(t, "real-package-screen.spk")
	err := p.populateImageFields()
	if err != nil {
		t.Error(err)
	}

	if strings.ToLower(p.Thumbnail[0]) != "package_icon.png" {
		t.Errorf("Expected 1st image to be called package_icon.png but got %s", p.Thumbnail[0])
	}

	if strings.ToLower(p.Thumbnail[1]) != "package_icon_256.png" {
		t.Errorf("Expected 2nd image to be called package_icon_256.png but got %s", p.Thumbnail[0])
	}

	if strings.ToLower(p.ThumbnailRetina[0]) != "package_icon_256.png" {
		t.Errorf("Expceted retina icon called package_icon_256.png but got %s", p.ThumbnailRetina[0])
	}

	if strings.ToLower(p.Snapshot[0]) != "package_screen_01.png" {
		t.Errorf("Expceted 1st screen to be called package_screen_01.png but got %s", p.Snapshot[0])
	}

	if strings.ToLower(p.Snapshot[1]) != "package_screen_02.png" {
		t.Errorf("Expceted 2nd screen to be called package_screen_02.png but got %s", p.Snapshot[0])
	}

	if strings.ToLower(p.Snapshot[2]) != "package_screen_03.png" {
		t.Errorf("Expceted 3rd sceen to be called package_screen_03.png but got %s", p.Snapshot[0])
	}
}

func TestPopulateImageFieldsNoImages(t *testing.T) {
	// this package has no images
	cleanupPackage(t, "real-package-wizard.spk")
	p := skipIfMissingTestPackage(t, "real-package-wizard.spk")
	err := p.populateImageFields()
	if err != nil {
		t.Error(err)
	}

	if l := len(p.Thumbnail); l != 0 {
		t.Errorf("Expected no icon but got %d", l)
	}

	if l := len(p.ThumbnailRetina); l != 0 {
		t.Errorf("Expceted no retina icon but got %d", l)
	}

	if l := len(p.Snapshot); l != 0 {
		t.Errorf("Expcetedno screenshots but got but got %d", l)
	}
}

func TestFinished(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	cleanupPackage(t, "dot-slash-package.spk")
}
