package syno // import jdel.org/go-syno/syno

import (
	"errors"

	yaml "gopkg.in/yaml.v2"
)

// Package represents a Synology package
type Package struct {
	FileName string `json:"-" yaml:"-"` // Expose FileName but never serialize it
	// Mandatory fields
	Name             string `json:"package,omitempty" yaml:"package,omitempty"`
	DisplayName      string `json:"dname,omitempty" yaml:"dname,omitempty"` // DSM 6 Field, now undocumented
	DisplayName7     string `json:"displayname,omitempty" yaml:"displayname,omitempty"`
	Version          string `json:"version,omitempty" yaml:"version,omitempty"`
	Firmware         string `json:"firmware,omitempty" yaml:"firmware,omitempty"` // Replaced by OsMinimumVersion after 6.1-14715
	OSMinimumVersion string `json:"os_min_ver,omitempty" yaml:"os_min_ver,omitempty"`
	Description      string `json:"desc,omitempty" yaml:"desc,omitempty"` // DSM 6 Field, now undocumented
	Description7     string `json:"description,omitempty" yaml:"description,omitempty"`
	Arch             string `json:"arch,omitempty" yaml:"arch,omitempty"`
	Maintainer       string `json:"maintainer,omitempty" yaml:"maintainer,omitempty"`
	PackageIcon      string `json:"package_icon,omitempty" yaml:"package_icon,omitempty"`         // DSM 6 Field, now undocumented
	PackageIcon256   string `json:"package_icon_256,omitempty" yaml:"package_icon_256,omitempty"` // DSM 6 Field, now undocumented
	// Optional fields
	MaintainerURL                 string `json:"maintainer_url,omitempty" yaml:"maintainer_url,omitempty"`
	Distributor                   string `json:"distributor,omitempty" yaml:"distributor,omitempty"`
	DistributorURL                string `json:"distributor_url,omitempty" yaml:"distributor_url,omitempty"`
	SupportURL                    string `json:"support_url,omitempty" yaml:"support_url,omitempty"`
	SupportCenter                 string `json:"support_center,omitempty" yaml:"support_center,omitempty"`
	SupportMove                   string `json:"support_move,omitempty" yaml:"support_move,omitempty"`
	Model                         string `json:"model,omitempty" yaml:"model,omitempty"`
	ExcludeArch                   string `json:"exclude_arch,omitempty" yaml:"exclude_arch,omitempty"`   //maybe handle that ?
	ExcludeModel                  string `json:"exclude_model,omitempty" yaml:"exclude_model,omitempty"` //maybe handle that ?
	Checksum                      string `json:"checksum,omitempty" yaml:"checksum,omitempty"`
	AdminPort                     string `json:"adminport,omitempty" yaml:"adminport,omitempty"`
	AdminURL                      string `json:"adminurl,omitempty" yaml:"adminurl,omitempty"`
	AdminProtocol                 string `json:"adminprotocol,omitempty" yaml:"adminprotocol,omitempty"`
	DSMUIDir                      string `json:"dsmuidir,omitempty" yaml:"dsmuidir,omitempty"` // DSM 6 Field, now undocumented
	DSMAppDir                     string `json:"dsmappdir,omitempty" yaml:"dsmappdir,omitempty"`
	DSMAppPage                    string `json:"dsmapppage,omitempty" yaml:"dsmapppage,omitempty"`
	DSMAppName                    string `json:"dsmappname,omitempty" yaml:"dsmappname,omitempty"`
	DSMAppLaunchName              string `json:"dsmapplaunchname,omitempty" yaml:"dsmapplaunchname,omitempty"`
	Changelog                     string `json:"changelog,omitempty" yaml:"changelog,omitempty"` // DSM 6 Field, now undocumented
	CheckPort                     bool   `json:"checkport,omitempty" yaml:"checkport,omitempty"`
	Startable                     bool   `json:"startable,omitempty" yaml:"startable,omitempty"` // Replaced by ConstrolStop after 6.1-14907
	ConstrolStop                  bool   `json:"ctl_stop,omitempty" yaml:"ctl_stop,omitempty"`
	ConstrolUninstall             bool   `json:"ctl_uninstall,omitempty" yaml:"ctl_uninstall,omitempty"`
	PreCheckStartStop             bool   `json:"precheckstartstop,omitempty" yaml:"precheckstartstop,omitempty"`
	HelpURL                       string `json:"helpurl,omitempty" yaml:"helpurl,omitempty"`
	Beta                          bool   `json:"beta,omitempty" yaml:"beta,omitempty"`
	ReportURL                     string `json:"report_url,omitempty" yaml:"report_url,omitempty"`
	InstallReboot                 bool   `json:"install_reboot,omitempty" yaml:"install_reboot,omitempty"`
	InstallDepPackages            string `json:"install_dep_packages,omitempty" yaml:"install_dep_packages,omitempty"`
	InstallConflictPackages       string `json:"install_conflict_packages,omitempty" yaml:"install_conflict_packages,omitempty"`
	InstallBreakPackages          string `json:"install_break_packages,omitempty" yaml:"install_break_packages,omitempty"`
	InstallOnColdStorage          bool   `json:"install_on_cold_storage,omitempty" yaml:"install_on_cold_storage,omitempty"`
	InstallReplacePackages        string `json:"install_replace_packages,omitempty" yaml:"install_replace_packages,omitempty"`
	InstUninstRestartServices     string `json:"instuninst_restart_services,omitempty" yaml:"instuninst_restart_services,omitempty"` // DSM 6 Field, now undocumented
	StartStopRestartServices      string `json:"startstop_restart_services,omitempty" yaml:"startstop_restart_services,omitempty"`
	InstallDepServices            string `json:"install_dep_services,omitempty" yaml:"install_dep_services,omitempty"`
	StartDepServices              bool   `json:"start_dep_services,omitempty" yaml:"start_dep_services,omitempty"`
	ExtractSize                   string `json:"extractsize,omitempty" yaml:"extractsize,omitempty"`
	SupportConfFolder             bool   `json:"support_conf_folder,omitempty" yaml:"support_conf_folder,omitempty"` // Deprecated after 6.0
	InstallType                   string `json:"install_type,omitempty" yaml:"install_type,omitempty"`
	SilentInstall                 bool   `json:"silent_install,omitempty" yaml:"silent_install,omitempty"`
	SilentUpgrade                 bool   `json:"silent_upgrade,omitempty" yaml:"silent_upgrade,omitempty"`
	SilentUninstall               bool   `json:"silent_uninstall,omitempty" yaml:"silent_uninstall,omitempty"`
	AutoUpgradeFrom               string `json:"auto_upgrade_from,omitempty" yaml:"auto_upgrade_from,omitempty"`
	OfflineInstall                bool   `json:"offline_install,omitempty" yaml:"offline_install,omitempty"`
	ThirdParty                    bool   `json:"thirdparty,omitempty" yaml:"thirdparty,omitempty"`
	OSMaximumVersion              string `json:"os_max_ver,omitempty" yaml:"os_max_ver,omitempty"`
	UseDeprecatedReplaceMechanism bool   `json:"use_deprecated_replace_mechanism,omitempty" yaml:"use_deprecated_replace_mechanism,omitempty"`
	// Package Center metadata
	Start                bool     `json:"start,omitempty" yaml:"start,omitempty"`
	Price                string   `json:"price,omitempty" yaml:"price,omitempty"`
	DownloadCount        string   `json:"download_count,omitempty" yaml:"download_count,omitempty"`
	ServiceDependencies  string   `json:"depsers,omitempty" yaml:"depsers,omitempty"`
	PackagesDependencies string   `json:"deppkgs,omitempty" yaml:"deppkgs,omitempty"`
	PackagesConflicts    string   `json:"conflictpkgs,omitempty" yaml:"conflictpkgs,omitempty"`
	RecentDownloadCount  string   `json:"recent_download_count,omitempty" yaml:"recent_download_count,omitempty"`
	Link                 string   `json:"link,omitempty" yaml:"link,omitempty"`
	Size                 string   `json:"size,omitempty" yaml:"size,omitempty"`
	MD5                  string   `json:"md5,omitempty" yaml:"md5,omitempty"`
	QuickInstall         bool     `json:"qinst,omitempty" yaml:"qinst,omitempty"`
	QuickStart           bool     `json:"qstart,omitempty" yaml:"qstart,omitempty"`
	QuickUpgrade         bool     `json:"qupdate,omitempty" yaml:"qupdate,omitempty"`
	Thumbnail            []string `json:"thumbnail,omitempty" yaml:"thumbnail,omitempty"`
	ThumbnailRetina      []string `json:"thumbnail_retina,omitempty" yaml:"thumbnail_retina,omitempty"`
	Snapshot             []string `json:"snapshot,omitempty" yaml:"snapshot,omitempty"`
	Category             string   `json:"category,omitempty" yaml:"category,omitempty"`
	SubCategory          string   `json:"subcategory,omitempty" yaml:"subcategory,omitempty"`
	Type                 string   `json:"type,omitempty" yaml:"type,omitempty"`
	// Unsupported fields (for go API usage only)
	I18nDisplayNames map[string]string `json:"i18n_dnames,omitempty" yaml:"i18n_dnames,omitempty"`
	I18nDescriptions map[string]string `json:"i18n_desc,omitempty" yaml:"i18n_desc,omitempty"`
}

func (p *Package) String() string {
	yamlPackage, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(yamlPackage)
}

// Packages is a slice of *Package
type Packages []*Package

func (p *Packages) String() string {
	yamlPackages, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(yamlPackages)
}

// sort.Sort Interface implementtion
func (p Packages) Len() int {
	return len(p)
}
func (p Packages) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}
func (p Packages) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// index returns the index of the first occurrence of synoPkg
// if it exists in synoPkgs.
// The comparison is done on the Name and Arch properties
func (p Packages) index(name, arch string) (int, error) {
	returnIndex := 0
	var err error
	err = errors.New("does not exist")

	for index, synoPkg := range p {
		if synoPkg.Name == name && synoPkg.Arch == arch {
			returnIndex = index
			err = nil
			break
		}
	}
	return returnIndex, err
}
