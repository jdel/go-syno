package syno // import jdel.org/go-syno/syno

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func executablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func sliceOfStringsContains(s []string, query string) bool {
	for _, a := range s {
		if a == query {
			return true
		}
	}
	return false
}

func trimDotSlash(s string) string {
	if s[:2] == "./" {
		s = s[2:]
	}
	return s
}

func writeToFile(r io.Reader, f string, m os.FileMode) error {
	os.MkdirAll(filepath.Dir(f), m)
	outputFile, err := os.OpenFile(f, os.O_CREATE|os.O_WRONLY, m)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, r)
	return err
}

func parseBoolOrYes(value string) (bool, error) {
	if value == "yes" {
		value = "true"
	}
	return strconv.ParseBool(value)
}

func sliceOfStringsItemMatches(s []string, query string) []string {
	var output []string
	for _, i := range s {
		if match, _ := regexp.MatchString(query, i); match {
			output = append(output, i)
		}
	}
	return output
}
