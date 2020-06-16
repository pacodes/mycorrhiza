package fs

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

const (
	HyphaPattern    = `[^\s\d:/?&\\][^:?&\\]*`
	hyphaPattern    = `[^\s\d:/?&\\][^:?&\\]*`
	revisionPattern = `[\d]+`
	revTxtPattern   = revisionPattern + `\.txt`
	revBinPattern   = revisionPattern + `\.bin`
	metaJsonPattern = `meta\.json`
	RevQuery        = `{rev:` + revisionPattern + `}`
	HyphaUrl        = `/` + hyphaPattern
)

var (
	leadingInt = regexp.MustCompile(`^[-+]?\d+`)
)

func matchNameToEverything(name string) (hyphaM bool, revTxtM bool, revBinM bool, metaJsonM bool) {
	simpleMatch := func(s string, p string) bool {
		m, _ := regexp.MatchString(p, s)
		return m
	}
	switch {
	case simpleMatch(name, revTxtPattern):
		revTxtM = true
	case simpleMatch(name, revBinPattern):
		revBinM = true
	case simpleMatch(name, metaJsonPattern):
		metaJsonM = true
	case simpleMatch(name, hyphaPattern):
		hyphaM = true
	}
	return
}

func stripLeadingInt(s string) string {
	return leadingInt.FindString(s)
}

func hyphaDirRevsValidate(dto map[string]map[string]string) (res bool) {
	for k, _ := range dto {
		switch k {
		case "0":
			delete(dto, "0")
		default:
			res = true
		}
	}
	return res
}

func scanHyphaDir(fullPath string) (valid bool, revs map[string]map[string]string, possibleSubhyphae []string, metaJsonPath string, err error) {
	nodes, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return // implicit return values
	}

	for _, node := range nodes {
		hyphaM, revTxtM, revBinM, metaJsonM := matchNameToEverything(node.Name())
		switch {
		case hyphaM && node.IsDir():
			possibleSubhyphae = append(possibleSubhyphae, filepath.Join(fullPath, node.Name()))
		case revTxtM && !node.IsDir():
			revId := stripLeadingInt(node.Name())
			revs[revId]["txt"] = filepath.Join(fullPath, node.Name())
		case revBinM && !node.IsDir():
			revId := stripLeadingInt(node.Name())
			revs[revId]["bin"] = filepath.Join(fullPath, node.Name())
		case metaJsonM && !node.IsDir():
			metaJsonPath = filepath.Join(fullPath, "meta.json")
			// Other nodes are ignored. It is not promised they will be ignored in future versions
		}
	}

	valid = hyphaDirRevsValidate(revs)

	return // implicit return values
}
