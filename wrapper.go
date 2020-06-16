package main

import (
	"net/http"
)

type HyphaStorage interface {
	ByName(string) (Hypha, error)
	ByNameRevision(string, string) (Revision, error)
	Init(string)
}

type Hypha interface {
	AddChild(string)
	AsHtml(map[string]Hypha) (string, error)
	Name() string
	NewestRevision() string
	ParentName() string
}

type Revision interface {
	ActionGetBinary(http.ResponseWriter)
	ActionRaw(http.ResponseWriter)
	ActionZen(http.ResponseWriter, HyphaStorage)
	ActionView(http.ResponseWriter, HyphaStorage, func(map[string]Hypha, Revision, string) string)
	AsHtml(HyphaStorage) (string, error)
	Name() string
}
