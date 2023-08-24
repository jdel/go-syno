package syno // import jdel.org/go-syno/syno

import (
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	yaml "gopkg.in/yaml.v2"
)

// Model represents a Synology model
type Model struct {
	Name    string `json:"name,omitempty"`
	CPU     string `json:"cpu,omitempty"`
	Cores   string `json:"cores,omitempty"`
	Threads string `json:"threads,omitempty"`
	FPU     string `json:"fpu,omitempty"`
	Arch    string `json:"arch,omitempty"`
	RAM     string `json:"ram,omitempty"`
}

// Models is a slice of Model
type Models []*Model

// Families contains families - arch mappings
var Families map[string][]string

func init() {
	// https://github.com/SynologyOpenSource/pkgscripts-ng/blob/master/include/pkg_util.sh#L107
	Families = make(map[string][]string)
	Families["x86_64"] = []string{"x86", "bromolow", "cedarview", "avoton", "braswell", "broadwell", "dockerx64", "kvmx64", "grantley", "denverton", "apollolake"}
	Families["i686"] = []string{"evansport"}
	Families["armv5"] = []string{"88f6281", "88f6282", "88f5281"}
	Families["armv7"] = []string{"alpine", "alpine4k", "ipq806x", "northstarplus"}
	Families["armv8"] = []string{"rtd1296"}
	Families["ppc"] = []string{"ppc854x", "ppc853x", "ppc824x", "powerpc", "qoriq"}
}

// GetModels returns Synology models from file or fall back to web craling
func GetModels(forceRefresh bool) (Models, error) {
	if modelsFileExists() && !forceRefresh {
		return getModelsFromModelsFile()
	}
	models, err := getModelsFromInternet()
	if len(models) != 0 && err == nil {
		models.SaveModelsFile()
	}
	return models, err
}

// FilterByName filters models by name
func (m Models) FilterByName(query string) Models {
	output := Models{}
	for _, synoModel := range m {
		if strings.Contains(strings.ToLower(synoModel.Name), strings.ToLower(query)) {
			output = append(output, synoModel)
		}
	}
	return output
}

// SaveModelsFile saves the model file to o.ModelsFile
func (m Models) SaveModelsFile() error {
	yamlModels, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	err = os.WriteFile(o.ModelsFile, yamlModels, 0755)
	if err != nil {
		return err
	}
	return nil
}

func modelsFileExists() bool {
	if _, err := os.Stat(o.ModelsFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func getModelsFromModelsFile() (Models, error) {
	var models Models
	bytes, err := os.ReadFile(o.ModelsFile)
	if err != nil {
		return nil, err
	}

	if yaml.Unmarshal(bytes, &models); err != nil {
		return nil, err
	}
	return models, nil
}

// getModelsFromInternet fetches Synology models from
// The official Synology kb
func getModelsFromInternet() (Models, error) {
	resp, err := http.Get("https://kb.synology.com/en-me/DSM/tutorial/What_kind_of_CPU_does_my_NAS_have")
	if err != nil && resp != nil && resp.StatusCode != 200 && resp.StatusCode != 302 {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// All the html we want is embedded in the last <script> tag
	var lastDeferedScript string
	doc.Find("script[defer]").Each(func(i int, s *goquery.Selection) {
		lastDeferedScript = s.Text()
	})

	// Extract the content capture group
	r := regexp.MustCompile(`"content":"(.*)","tags"`)
	match := r.FindStringSubmatch(lastDeferedScript)[1]

	// Let otto handle the double quoted JS from content
	vm := otto.New()
	inlineHtml, err := vm.Run(`"` + match + `"`)
	if err != nil {
		return nil, err
	}

	// New sub document to parse inline HTML
	subDoc, err := goquery.NewDocumentFromReader(strings.NewReader(inlineHtml.String()))
	if err != nil {
		return nil, err
	}

	var models Models

	subDoc.Find("tr").Each(func(i int, s *goquery.Selection) {
		tds := s.ChildrenFiltered("td")
		if tds.Size() == 7 {
			model := &Model{
				Name:    tds.Eq(0).Text(),
				CPU:     tds.Eq(1).Text(),
				Cores:   tds.Eq(2).Text(),
				Threads: tds.Eq(3).Text(),
				FPU:     tds.Eq(4).Text(),
				Arch:    tds.Eq(5).Text(),
				RAM:     tds.Eq(6).Text(),
			}
			models = append(models, model)
		}
	})
	return models, nil
}
