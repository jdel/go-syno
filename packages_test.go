package syno_test // import jdel.org/go-syno/syno_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	syno "jdel.org/go-syno"
)

var pX8664,
	pCedarview1,
	pCedarview2,
	pCedarview3,
	pX86,
	pEvansport,
	pMulti,
	p1,
	p2 *syno.Package

func init() {
	p1 = &syno.Package{
		Name:        "package1",
		DisplayName: "Bleu",
		Arch:        "noarch",
		Firmware:    "3.0",
		Beta:        true,
	}
	p2 = &syno.Package{
		Name:        "package2",
		DisplayName: "French Emmental",
		Arch:        "noarch",
		Firmware:    "6.1",
	}
	pX8664 = &syno.Package{
		Name:        "x86-64-package",
		DisplayName: "Swiss Emmental",
		Arch:        "x86_64",
	}
	pCedarview1 = &syno.Package{
		Name:        "cedarview-package1",
		DisplayName: "Brie",
		Arch:        "cedarview",
		Beta:        true,
	}
	pCedarview2 = &syno.Package{
		Name:        "cedarview-package",
		DisplayName: "Cedarview Package 2",
		Arch:        "cedarview",
		Version:     "1.0",
	}
	pCedarview3 = &syno.Package{
		Name:        "cedarview-package",
		DisplayName: "Cedarview Package 3",
		Arch:        "cedarview",
		Version:     "1.1",
	}
	pMulti = &syno.Package{
		Name:        "multi-arch-package",
		DisplayName: "Multi Arch Package",
		Arch:        "x86_64 apollolake broadwell",
		Version:     "1.1",
	}
	pX86 = &syno.Package{
		Name: "x86-package",
		Arch: "x86",
	}
	pEvansport = &syno.Package{
		Name: "evansport-package",
		Arch: "evansport",
		Beta: true,
	}
}

func skipIfMissingTestPackage(t *testing.T, f string) *syno.Package {
	var p *syno.Package
	var err error
	if p, err = syno.NewPackage(f); err != nil {
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

func TestRealPackage(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	var p *syno.Package
	var err error
	if p, err = syno.NewPackage("real-package.spk"); err != nil {
		t.Error("Package real-package.spk should exist in tests/packages")
	}

	if p.Name != "real-package" {
		t.Error("Name should be real-package")
	}
	if p.DisplayName != "Real Package" {
		t.Error("Display Name should be Real Package")
	}
}

func TestRealPackageI18n(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	var p *syno.Package
	var err error
	previousLanguage := o.Language
	defer func() {
		o.Language = previousLanguage
	}()
	o.Language = "ita"
	if p, err = syno.NewPackage("real-package-i18n.spk"); err != nil {
		t.Error("Package real-package-i18n.spk should exist in tests/packages")
	}
	if p.DisplayName != "Real Package (ITA)" {
		t.Errorf("Display Name should be Real Package (ITA) but got %s", p.DisplayName)
	}
	if p.Description != "A Real Package (ITA)" {
		t.Errorf("Description should be A Real Package (ITA) but got %s", p.Description)
	}
}

func TestRealPackageWizard(t *testing.T) {
	cleanupPackage(t, "real-package-wizard.spk")
	var p *syno.Package
	var err error
	if p, err = syno.NewPackage("real-package-wizard.spk"); err != nil {
		t.Error("Package real-package-wizard.spk should exist in tests/packages")
	}

	if p.QuickInstall {
		t.Error("qinst should be false")
	}
	if p.QuickStart {
		t.Error("qstart should be false")
	}
	if p.QuickUpgrade {
		t.Error("qupgrade should be false")
	}
}

func TestRealPackageWizardQfields(t *testing.T) {
	cleanupPackage(t, "real-package-wizard-qfields.spk")
	var p *syno.Package
	var err error
	if p, err = syno.NewPackage("real-package-wizard-qfields.spk"); err != nil {
		t.Error("Package real-package-wizard-qfields.spk should exist in tests/packages")
	}

	if !p.QuickInstall {
		t.Error("qinst should be true")
	}
	if !p.QuickStart {
		t.Error("qstart should be true")
	}
	if !p.QuickUpgrade {
		t.Error("qupgrade should be true")
	}
}

func TestBadPackage(t *testing.T) {
	cleanupPackage(t, "bad-package.spk")
	if _, err := syno.NewPackage("bad-package.spk"); err == nil {
		t.Error("Package bad-package.spk should instantiate a Package")
	}
}

func TestDeadPackage(t *testing.T) {
	if _, err := syno.NewPackage("dead-package.spk"); err == nil {
		t.Error("Package dead-package.spk should not exist in tests/packages")
	}
}

func TestSortPackages(t *testing.T) {
	pp := syno.Packages{p2, p1}
	sort.Sort(pp)

	if name := pp[0].Name; name != "package1" {
		t.Errorf("Expected first package to be package1 but got %s", name)
	}
}

func TestRealPackageExistsOnDisk(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	p := skipIfMissingTestPackage(t, "real-package.spk")
	if !p.ExistsOnDisk() {
		t.Errorf("Expected %s to exist on disk", p.FullPath())
	}
}

func TestDeadPackageExistsOnDisk(t *testing.T) {
	p, _ := syno.NewPackage("dead-package.spk")
	if p != nil && p.ExistsOnDisk() {
		t.Errorf("Expected %s to not exist on disk", p.FullPath())
	}
}

func TestFilterByArch(t *testing.T) {
	pp := syno.Packages{p1, p2, pCedarview1, pCedarview2, pCedarview3, pEvansport, pX86, pX8664}
	ppf := pp.FilterByArch("fake-arch")
	if len(ppf) != 2 {
		t.Errorf("Expected 2 packages to match but got %+v", ppf)
	}
	ppf = pp.FilterByArch("cedarview")
	if len(ppf) != 6 {
		t.Errorf("Expected 6 packages to match but got %+v", ppf)
	}
	ppf = pp.FilterByArch("x86_64")
	if len(ppf) != 7 {
		t.Errorf("Expected 7 packages to match but got %+v", ppf)
	}
	ppf = pp.FilterByArch("evansport")
	if len(ppf) != 3 {
		t.Errorf("Expected 3 packages to match but got %+v", ppf)
	}
}

func TestFilterByArchMulti(t *testing.T) {
	pp := syno.Packages{pMulti}
	archs := strings.Fields(pMulti.Arch)
	for _, arch := range archs {
		ppf := pp.FilterByArch(arch)
		if len(ppf) != 1 {
			t.Errorf("Expected 1 package to match %q but got %+v", arch, ppf)
		}
	}
}

func TestFilterByFirmware(t *testing.T) {
	pp := syno.Packages{p1, p2}
	ppf := pp.FilterByFirmware("6.1")
	if len(ppf) != 2 {
		t.Errorf("Expected 2 packages to match but got %+v", ppf)
	}
	ppf = pp.FilterByFirmware("4.2")
	if len(ppf) != 1 {
		t.Errorf("Expected 1 packages to match but got %+v", ppf)
	}
	ppf = pp.FilterByFirmware("2.0")
	if len(ppf) != 0 {
		t.Errorf("Expected 0 packages to match but got %+v", ppf)
	}
}

func TestFilterOutBeta(t *testing.T) {
	pp := syno.Packages{p1, p2, pCedarview1, pCedarview2, pCedarview3, pEvansport, pX86, pX8664}
	ppf := pp.FilterOutBeta()
	if len(ppf) != 5 {
		t.Errorf("Expected 5 packages to match but got %+v", ppf)
	}
}

func TestFilterByName(t *testing.T) {
	pp := syno.Packages{p1, p2, pX8664, pCedarview1}
	ppf := pp.SearchByName("B")
	if len(ppf) != 2 {
		t.Errorf("Expected 2 packages to match but got %+v", ppf)
	}
	ppf = pp.SearchByName("Emmental")
	if len(ppf) != 2 {
		t.Errorf("Expected 2 packages to match but got %+v", ppf)
	}
	ppf = pp.SearchByName("ch Emmen")
	if len(ppf) != 1 {
		t.Errorf("Expected 1 packages to match but got %+v", ppf)
	}
	ppf = pp.SearchByName("Brie")
	if len(ppf) != 1 {
		t.Errorf("Expected 1 packages to match but got %+v", ppf)
	}
}
func TestOnlyShowLastVersion(t *testing.T) {
	pp := syno.Packages{pCedarview2, pCedarview3}
	ppf := pp.OnlyShowLastVersion()
	if len(ppf) != 1 {
		t.Errorf("Expected 1 packages to match but got %+v", ppf)
	}
}

func TestNewDebugPackage(t *testing.T) {
	if p := syno.NewDebugPackage("Test"); p.Name == "" {
		t.Error("Expecting debug package to have a Name")
	} else if p.DisplayName == "" {
		t.Error("Expecting debug package to have a DisplayName")
	} else if p.Description != "Test" {
		t.Error("Expecting debug package to have a Description that reads Test")
	}
}

func TestFinished(t *testing.T) {
	cleanupPackage(t, "real-package.spk")
	cleanupPackage(t, "real-package-i18n.spk")
	cleanupPackage(t, "dot-slash-package.spk")
	cleanupPackage(t, "bad-package.spk")
}
