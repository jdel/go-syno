package syno // import jdel.org/go-syno/syno

import (
	"errors"
)

// Package represents a Synology package
type Package struct {
	fileName string
	// Mandatory fields
	Name             string            `json:"package,omitempty"`
	DisplayName      string            `json:"dname,omitempty"`
	I18nDisplayNames map[string]string `json:"i18n_dnames,omitempty"`
	Version          string            `json:"version,omitempty"`
	Firmware         string            `json:"firmware,omitempty"` // Replaced by OsMinimumVersion after 6.1-14715
	OSMinimumVersion string            `json:"os_min_ver,omitempty"`
	Description      string            `json:"desc,omitempty"`
	I18nDescriptions map[string]string `json:"i18n_desc,omitempty"`
	Arch             string            `json:"arch,omitempty"`
	Maintainer       string            `json:"maintainer,omitempty"`
	PackageIcon      string            `json:"package_icon,omitempty"`
	PackageIcon256   string            `json:"package_icon_256,omitempty"`
	// Optional fields
	MaintainerURL             string `json:"maintainer_url,omitempty"`
	Distributor               string `json:"distributor,omitempty"`
	DistributorURL            string `json:"distributor_url,omitempty"`
	SupportURL                string `json:"support_url,omitempty"`
	SupportCenter             string `json:"support_center,omitempty"`
	Model                     string `json:"model,omitempty"`
	ExcludeArch               string `json:"exclude_arch,omitempty"` //maybe handle that ?
	Checksum                  string `json:"checksum,omitempty"`
	AdminPort                 string `json:"adminport,omitempty"`
	AdminURL                  string `json:"adminurl,omitempty"`
	AdminProtocol             string `json:"adminprotocol,omitempty"`
	DSMUIDir                  string `json:"dsmuidir,omitempty"`
	DSMAppDir                 string `json:"dsmappdir,omitempty"`
	Changelog                 string `json:"changelog,omitempty"`
	CheckPort                 bool   `json:"checkport,omitempty"`
	Startable                 bool   `json:"startable,omitempty"` // Replaced by ConstrolStop after 6.1-14907
	ConstrolStop              bool   `json:"ctl_stop,omitempty"`
	ConstrolUninstall         bool   `json:"ctl_uninstall,omitempty"`
	PreCheckStartStop         bool   `json:"precheckstartstop,omitempty"`
	HelpURL                   string `json:"helpurl,omitempty"`
	Beta                      bool   `json:"beta,omitempty"`
	ReportURL                 string `json:"report_url,omitempty"`
	InstallReboot             bool   `json:"install_reboot,omitempty"`
	InstallDepPackages        string `json:"install_dep_packages,omitempty"`
	InstallConflictPackages   string `json:"install_conflict_packages,omitempty"`
	InstallBreakPackages      string `json:"install_break_packages,omitempty"`
	InstallReplacePackages    string `json:"install_replace_packages,omitempty"`
	InstUninstRestartServices string `json:"instuninst_restart_services,omitempty"`
	StartStopRestartServices  string `json:"startstop_restart_services,omitempty"`
	InstallDepServices        string `json:"install_dep_services,omitempty"`
	StartDepServices          bool   `json:"start_dep_services,omitempty"`
	ExtractSize               string `json:"extractsize,omitempty"`
	SupportConfFolder         bool   `json:"support_conf_folder,omitempty"` // Deprecated after 6.0
	InstallType               string `json:"install_type,omitempty"`
	SilentInstall             bool   `json:"silent_install,omitempty"`
	SilentUpgrade             bool   `json:"silent_upgrade,omitempty"`
	SilentUninstall           bool   `json:"silent_uninstall,omitempty"`
	AutoUpgradeFrom           string `json:"auto_upgrade_from,omitempty"`
	OfflineInstall            bool   `json:"offline_install,omitempty"`
	ThirdParty                bool   `json:"thirdparty,omitempty"`
	OSMaximumVersion          bool   `json:"os_max_ver,omitempty"`
	// Package Center metadata
	Start                bool     `json:"start,omitempty"`
	Price                string   `json:"price,omitempty"`
	DownloadCount        string   `json:"download_count,omitempty"`
	ServiceDependencies  string   `json:"depsers,omitempty"`
	PackagesDependencies string   `json:"deppkgs,omitempty"`
	PackagesConflicts    string   `json:"conflictpkgs,omitempty"`
	RecentDownloadCount  string   `json:"recent_download_count,omitempty"`
	Link                 string   `json:"link,omitempty"`
	Size                 string   `json:"size,omitempty"`
	MD5                  string   `json:"md5,omitempty"`
	QuickInstall         bool     `json:"qinst,omitempty"`
	QuickStart           bool     `json:"qstart,omitempty"`
	QuickUpgrade         bool     `json:"qupdate,omitempty"`
	Thumbnail            []string `json:"thumbnail,omitempty"`
	ThumbnailRetina      []string `json:"thumbnail_retina,omitempty"`
	Snapshot             []string `json:"snapshot,omitempty"`
	Category             string   `json:"category,omitempty"`
	SubCategory          string   `json:"subcategory,omitempty"`
	Type                 string   `json:"type,omitempty"`
}

// Packages is a slice of Package
type Packages []*Package

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
