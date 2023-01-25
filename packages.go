package syno // import jdel.org/go-syno/syno

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/calmh/versions"
	ini "gopkg.in/ini.v1"
)

// NewPackage creates a new package from the spk file
func NewPackage(synoPackageName string) (*Package, error) {
	var err error
	synoPkg := Package{}
	synoPkg.fileName = synoPackageName
	synoPkg.Thumbnail = make([]string, 0)
	synoPkg.ThumbnailRetina = make([]string, 0)
	synoPkg.Snapshot = make([]string, 0)

	if err = synoPkg.pupulateINFOFields(); err != nil {
		return nil, err
	}
	err = synoPkg.populateImageFields()
	synoPkg.populatePackageCenterFields()

	return &synoPkg, err
}

// NewDebugPackage creates a new debug package from the description string
func NewDebugPackage(description string) *Package {
	return &Package{
		Name:        "debug",
		DisplayName: "GoSSPKS Debug",
		Arch:        "noarch",
		Firmware:    "1",
		Version:     "v1",
		Beta:        true,
		Maintainer:  "gosspks",
		Description: description,
		Thumbnail:   make([]string, 0),
	}
}

// FullPath returns the package full path on FS
func (p *Package) FullPath() string {
	return filepath.Join(o.PackagesDir, p.fileName)
}

// ExistsOnDisk returns true if the package file exists on disk
func (p *Package) ExistsOnDisk() bool {
	return fileExists(p.FullPath())
}

// FilterByArch filters synopkgs where Arch = query
func (p Packages) FilterByArch(query string) Packages {
	output := Packages{}

	for _, synoPkg := range p {
		if synoPkg.Arch == query || sliceOfStringsContains(strings.Fields(synoPkg.Arch), query) || synoPkg.familyOrArchMatch(query) || synoPkg.Arch == "noarch" {
			output = append(output, synoPkg)
		}
	}
	return output
}

// FilterByFirmware filters synopkgs where Version >= query
func (p Packages) FilterByFirmware(query string) Packages {
	output := Packages{}
	for _, synoPkg := range p {
		if synoPkg.Firmware <= query {
			output = append(output, synoPkg)
		}
	}
	return output
}

// FilterOutBeta returns synopkgs except beta packages
func (p Packages) FilterOutBeta() Packages {
	output := Packages{}
	for _, synoPkg := range p {
		if !synoPkg.Beta {
			output = append(output, synoPkg)
		}
	}
	return output
}

// SearchByName filters synopkgs name contains query
func (p Packages) SearchByName(query string) Packages {
	output := Packages{}
	for _, synoPkg := range p {
		if strings.Contains(strings.ToLower(synoPkg.DisplayName), strings.ToLower(query)) {
			output = append(output, synoPkg)
		}
	}
	return output
}

// OnlyShowLastVersion overrides an existing package
// with a new one if the version is greater
// the comparison is done on the Name property
func (p Packages) OnlyShowLastVersion() Packages {
	output := Packages{}
	for _, synoPkg := range p {
		if pkgIndex, err := output.index(synoPkg.Name, synoPkg.Arch); err != nil {
			// Not found
			output = append(output, synoPkg)
		} else if versions.Compare(synoPkg.Version, output[pkgIndex].Version) > 0 {
			// Newer, overwrite
			output[pkgIndex] = synoPkg
		}
	}
	return output
}

func (p *Package) pupulateINFOFields() error {
	infoINI, err := p.parseInfo()
	if err != nil {
		return err
	}

	p.populateMandatoryFields(infoINI)
	p.populateI18nFields(infoINI)
	p.populateOptionalFields(infoINI)
	p.populateQFields(infoINI)
	return err
}

func (p *Package) populateMandatoryFields(infoINI *ini.File) {
	p.Name = infoINI.Section("").Key("package").Value()
	p.Version = infoINI.Section("").Key("version").Value()
	p.Firmware = infoINI.Section("").Key("firmware").Value()
	p.Arch = infoINI.Section("").Key("arch").Value()
	p.Maintainer = infoINI.Section("").Key("maintainer").Value()
}

func (p *Package) populateOptionalFields(infoINI *ini.File) {
	p.MaintainerURL = infoINI.Section("").Key("maintainer_url").Value()
	p.Distributor = infoINI.Section("").Key("distributor").Value()
	p.DistributorURL = infoINI.Section("").Key("distributor_url").Value()
	p.SupportURL = infoINI.Section("").Key("support_url").Value()
	p.Model = infoINI.Section("").Key("model").Value()
	p.ExcludeArch = infoINI.Section("").Key("exclude_arch").Value()
	p.Changelog = infoINI.Section("").Key("changelog").Value()
	p.Checksum = infoINI.Section("").Key("checksum").Value()
	p.AdminPort = infoINI.Section("").Key("adminport").Value()
	p.AdminURL = infoINI.Section("").Key("adminurl").Value()
	p.AdminProtocol = infoINI.Section("").Key("adminprotocol").Value()
	p.DSMUIDir = infoINI.Section("").Key("dsmuidir").Value()
	p.DSMAppDir = infoINI.Section("").Key("dsmappdir").Value()
	p.CheckPort, _ = parseBoolOrYes(infoINI.Section("").Key("checkport").Value())
	p.Startable, _ = parseBoolOrYes(infoINI.Section("").Key("startable").Value())
	p.PreCheckStartStop, _ = parseBoolOrYes(infoINI.Section("").Key("precheckstartstop").Value())
	p.HelpURL = infoINI.Section("").Key("helpurl").Value()
	p.Beta, _ = parseBoolOrYes(infoINI.Section("").Key("beta").Value())
	p.ReportURL = infoINI.Section("").Key("report_url").Value()
	p.InstallReboot, _ = parseBoolOrYes(infoINI.Section("").Key("install_reboot").Value())
	p.InstallDepPackages = infoINI.Section("").Key("install_dep_packages").Value()
	p.InstallConflictPackages = infoINI.Section("").Key("install_conflict_packages").Value()
	p.InstUninstRestartServices = infoINI.Section("").Key("instuninst_restart_services").Value()
	p.StartStopRestartServices = infoINI.Section("").Key("startstop_restart_services").Value()
	p.InstallDepServices = infoINI.Section("").Key("install_dep_services").Value()
	p.StartDepServices, _ = parseBoolOrYes(infoINI.Section("").Key("start_dep_services").Value())
	p.ExtractSize = infoINI.Section("").Key("extractsize").Value()
	p.SupportConfFolder, _ = parseBoolOrYes(infoINI.Section("").Key("support_conf_folder").Value())
	p.InstallType = infoINI.Section("").Key("install_type").Value()
	p.SilentInstall, _ = parseBoolOrYes(infoINI.Section("").Key("silent_install").Value())
	p.SilentUpgrade, _ = parseBoolOrYes(infoINI.Section("").Key("silent_upgrade").Value())
	p.SilentUninstall, _ = parseBoolOrYes(infoINI.Section("").Key("silent_uninstall").Value())
	p.AutoUpgradeFrom = infoINI.Section("").Key("auto_upgrade_from").Value()
	p.OfflineInstall, _ = parseBoolOrYes(infoINI.Section("").Key("offline_install").Value())
	p.ThirdParty, _ = parseBoolOrYes(infoINI.Section("").Key("thirdparty").Value())
}

func (p *Package) populateQFields(infoINI *ini.File) {
	wizardFile, _ := p.containsFiles("WIZARD_UIFILES")
	p.QuickInstall = determineQfieldValue(infoINI.Section("").Key("qinst").Value(), wizardFile)
	p.QuickStart = determineQfieldValue(infoINI.Section("").Key("qstart").Value(), wizardFile)
	p.QuickUpgrade = determineQfieldValue(infoINI.Section("").Key("qupgrade").Value(), wizardFile)
}

func (p *Package) populateI18nFields(infoINI *ini.File) {
	p.DisplayName = infoINI.Section("").Key("displayname").Value()
	if value := infoINI.Section("").Key(fmt.Sprintf("displayname_%s", o.Language)).Value(); value != "" {
		p.DisplayName = value
	}

	p.Description = infoINI.Section("").Key("description").Value()
	if value := infoINI.Section("").Key(fmt.Sprintf("description_%s", o.Language)).Value(); value != "" {
		p.Description = value
	}
}

func (p *Package) populateImageFields() error {
	extractedImages, err := p.getOrExtractImages()
	p.Thumbnail = sliceOfStringsItemMatches(extractedImages, "(?i)(_thumb_[0-9]+|_icon.*).png$")
	p.ThumbnailRetina = sliceOfStringsItemMatches(extractedImages, "(?i)(_thumb_256|_icon_256).png$")
	p.Snapshot = sliceOfStringsItemMatches(extractedImages, "(?i)_screen_[0-9]+.png$$")
	return err
}

func (p *Package) populatePackageCenterFields() {
	// JSON gosspks fields
	// synoPkg.Price = ""
	// synoPkg.DownloadCount = ""
	// synoPkg.RecentDownloadCount = ""
	// synoPkg.Link = ""
	// synoPkg.Category = ""
	// synoPkg.SubCategory = ""
	// synoPkg.Type = ""
	p.Size, _ = p.getSize()
	if o.MD5 {
		p.MD5, _ = p.getMD5()
	}
	p.Start = true
}

func determineQfieldValue(iniValue string, wizardFIles bool) bool {
	if value, err := parseBoolOrYes(iniValue); err == nil {
		return value
	} else if wizardFIles {
		return false
	}
	return true
}

func archIsInFamilyAndMatches(arch, query string) bool {
	if archs := Families[arch]; archs != nil {
		return sliceOfStringsContains(archs, query)
	}
	return false
}

func (p *Package) familyOrArchMatch(query string) bool {
	return archIsInFamilyAndMatches(p.Arch, query) || archIsInFamilyAndMatches(query, p.Arch)
}
