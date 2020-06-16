package fs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
)

type HyphaStorage struct {
	RootWikiDir string
	hyphae      map[string]Hypha
}

func (hs *HyphaStorage) ByNameRevision(name string, revId string) (Revision, error) {
	for revName, rev := range hs.hyphae[name].Revisions {
		if revId == revName {
			return *rev, nil
		}
	}
	return Revision{}, errors.New("No such name/revision")
}

// Hypha name is rootWikiDir/{here}
func (hs *HyphaStorage) hyphaName(fullPath string) string {
	return fullPath[len(hs.RootWikiDir)+1:]
}

func (hs *HyphaStorage) Init(rwd string) {
	hs.RootWikiDir = rwd
	hs.findHyphae(rwd)
}

func (hs *HyphaStorage) findHyphae(fullPath string) {
	valid, revs, possibleSubhyphae, metaJsonPath, err := scanHyphaDir(fullPath)
	if err != nil {
		return
	}

	// First, let's process subhyphae
	for _, possibleSubhypha := range possibleSubhyphae {
		hs.findHyphae(possibleSubhypha)
	}

	// This folder is not a hypha itself, nothing to do here
	if !valid {
		return
	}

	// Template hypha struct. Other fields are default json values.
	h := Hypha{
		FullName:   hs.hyphaName(fullPath),
		Path:       fullPath,
		parentName: filepath.Dir(hs.hyphaName(fullPath)),
		// Children names are unknown now
	}

	metaJsonContents, err := ioutil.ReadFile(metaJsonPath)
	if err != nil {
		log.Printf("Error when reading `%s`; skipping", metaJsonPath)
		return
	}
	err = json.Unmarshal(metaJsonContents, &h)
	if err != nil {
		log.Printf("Error when unmarshaling `%s`; skipping", metaJsonPath)
		return
	}

	// Fill in every revision paths
	for id, paths := range revs {
		if r, ok := h.Revisions[id]; ok {
			for fType, fPath := range paths {
				switch fType {
				case "bin":
					r.BinaryPath = fPath
				case "txt":
					r.TextPath = fPath
				}
			}
		} else {
			log.Printf("Error when reading hyphae from disk: hypha `%s`'s meta.json provided no information about revision `%s`, but files %s are provided; skipping\n", h.FullName, id, paths)
		}
	}

	// Now the hypha should be ok, gotta send structs
	hs.hyphae[h.FullName] = h
}
