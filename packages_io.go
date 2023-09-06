package syno // import jdel.org/go-syno/syno

import (
	"archive/tar"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ini "gopkg.in/ini.v1"
)

func (p *Package) containsFiles(fileNamePattern string) (bool, error) {
	var err error
	var spkFile *os.File

	spkFile, err = os.Open(filepath.Join(o.PackagesDir, p.FileName))
	if err != nil {
		return false, err
	}
	defer spkFile.Close()

	tarReader := tar.NewReader(spkFile)

	for {
		fileHeader, eofErr := tarReader.Next()
		if eofErr == io.EOF {
			break
		} else if eofErr != nil {
			return false, eofErr
		}

		if match, _ := regexp.MatchString(fileNamePattern, fileHeader.Name); match {
			return true, nil
		}
	}
	return false, nil
}

func (p *Package) extractFiles(fileNamePattern string) ([]string, error) {
	var err error
	var spkFile *os.File
	var extractedFiles []string

	spkFile, err = os.Open(filepath.Join(o.PackagesDir, p.FileName))
	if err != nil {
		return extractedFiles, err
	}
	defer spkFile.Close()

	tarReader := tar.NewReader(spkFile)

	for {
		fileHeader, eofErr := tarReader.Next()
		if eofErr == io.EOF {
			break
		} else if eofErr != nil {
			return extractedFiles, eofErr
		}

		if match, _ := regexp.MatchString(fileNamePattern, fileHeader.Name); match && fileHeader.Typeflag == tar.TypeReg {
			fileName := trimDotSlash(fileHeader.Name)
			extractedFiles = append(extractedFiles, fileName)
			outputPath := filepath.Join(o.CacheDir, p.FileName, fileName)
			err = writeToFile(tarReader, outputPath, 0755)
		}
	}
	return extractedFiles, err
}

func (p *Package) getSize() (string, error) {
	var size int
	var err error

	spkFileInfo, err := os.Stat(filepath.Join(o.PackagesDir, p.FileName))
	if err != nil {
		return "", err
	}

	size = int(spkFileInfo.Size())

	return strconv.Itoa(size), err
}

func (p *Package) getMD5() (string, error) {
	var calculatedMD5 string
	var err error

	spkFile, err := os.Open(filepath.Join(o.PackagesDir, p.FileName))
	if err != nil {
		return "", err
	}
	defer spkFile.Close()

	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, spkFile); err != nil {
		return "", err
	}
	calculatedMD5 = hex.EncodeToString(md5Hash.Sum(nil))

	return calculatedMD5, err
}

func (p *Package) extractInfo() error {
	extractedFiles, err := p.extractFiles("INFO")
	if err != nil || len(extractedFiles) == 0 {
		if err := os.Rename(
			filepath.Join(o.PackagesDir, p.FileName), fmt.Sprintf("%s.ignored", filepath.Join(o.PackagesDir, p.FileName))); err != nil {
			return err
		}
		return fmt.Errorf("quatantined %s: No INFO or not a tar file", p.FileName)
	}
	return nil
}

// Returns the info file path and an eventual error
func (p *Package) getOrExtractInfo() (string, error) {
	var err error
	infoINIPath := filepath.Join(o.CacheDir, p.FileName, "INFO")
	if _, err = os.Stat(infoINIPath); err != nil {
		if err = p.extractInfo(); err != nil {
			return "", err
		}
	}
	return infoINIPath, err
}

func (p *Package) parseInfo() (*ini.File, error) {
	spkinfoFilePath, err := p.getOrExtractInfo()
	if err != nil {
		return nil, fmt.Errorf("cannot extract INFO for %s: %s", p.FileName, err)
	}
	infoINI, err := ini.InsensitiveLoad(spkinfoFilePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read INFO for %s: %s", p.FileName, err)
	}
	infoINI.BlockMode = false
	return infoINI, nil
}

func (p *Package) extractImages() ([]string, error) {
	var err error
	var images []string
	if images, err = p.extractFiles("(?i)[^/]*png"); err != nil {
		return images, err
	}
	return images, err
}

func (p *Package) getImages() ([]string, error) {
	var err error
	var files []os.DirEntry
	var images []string

	files, err = os.ReadDir(filepath.Join(o.CacheDir, p.FileName))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		filename := f.Name()
		if filepath.Ext(strings.ToLower(filename)) == ".png" {
			images = append(images, filename)
		}
	}

	return images, err
}

func (p *Package) getOrExtractImages() ([]string, error) {
	extractedImages, err := p.getImages()
	if len(extractedImages) == 0 {
		extractedImages, err = p.extractImages()
	}
	return extractedImages, err
}
