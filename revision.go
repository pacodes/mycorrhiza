package main

import (
	"github.com/bouncepaw/mycorrhiza/fs"
	"net/http"
)

type Hypha interface {
	AddChild(string)
	AsHtml(map[string]Hypha) (string, error)
	NewestRevision() string
	ParentName() string
}

type Revision interface {
	AsHtml(map[string]Hypha) (string, error)
	ActionGetBinary(http.ResponseWriter)
	ActionRaw(http.ResponseWriter)
	ActionZen(http.ResponseWriter, map[string]Hypha)
	ActionView(http.ResponseWriter, map[string]Hypha, func(map[string]Hypha, Revision, string) string)
}

func GetRevision(hyphae map[string]Hypha, hyphaName string, revId string) (Revision, bool) {
	for revName, rev := range hyphae[hyphaName].Revisions {
		if revId == revName {
			return *rev, true
		}
	}
	return fs.Revision{}, false
}
