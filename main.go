package main

import (
	"fmt"
	"github.com/bouncepaw/mycorrhiza/fs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func RevInMap(m map[string]string) string {
	if val, ok := m["rev"]; ok {
		return val
	}
	return "0"
}

// handlers
func HandlerGetBinary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, err := hs.ByNameRevision(vars["hypha"], revno)
	if err != nil {
		return
	}
	rev.ActionGetBinary(w)
}

func HandlerRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, err := hs.ByNameRevision(vars["hypha"], revno)
	if err != nil {
		return
	}
	rev.ActionRaw(w)
}

func HandlerZen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, err := hs.ByNameRevision(vars["hypha"], revno)
	if err != nil {
		return
	}
	rev.ActionZen(w, hs)
}

func HandlerView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	revno := RevInMap(vars)
	rev, err := hs.ByNameRevision(vars["hypha"], revno)
	if err != nil {
		return
	}
	rev.ActionView(w, hs, HyphaPage)
}

func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerEdit(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerRewind(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerRename(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

func HandlerUpdate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	log.Println("Attempt to access an unimplemented thing")
}

var hs fs.HyphaStorage

const (
	revQuery = fs.RevQuery
	hyphaUrl = fs.HyphaUrl
)

func main() {
	if len(os.Args) == 1 {
		panic("Expected a root wiki pages directory")
	}
	rootWikiDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}
	hs.Init(rootWikiDir)

	// Start server code
	r := mux.NewRouter()

	r.Queries("action", "getBinary", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerGetBinary)
	r.Queries("action", "getBinary").Path(hyphaUrl).
		HandlerFunc(HandlerGetBinary)

	r.Queries("action", "raw", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerRaw)
	r.Queries("action", "raw").Path(hyphaUrl).
		HandlerFunc(HandlerRaw)

	r.Queries("action", "zen", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerZen)
	r.Queries("action", "zen").Path(hyphaUrl).
		HandlerFunc(HandlerZen)

	r.Queries("action", "view", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerView)
	r.Queries("action", "view").Path(hyphaUrl).
		HandlerFunc(HandlerView)

	r.Queries("action", "history").Path(hyphaUrl).
		HandlerFunc(HandlerHistory)

	r.Queries("action", "edit").Path(hyphaUrl).
		HandlerFunc(HandlerEdit)

	r.Queries("action", "rewind", "rev", revQuery).Path(hyphaUrl).
		HandlerFunc(HandlerRewind)

	r.Queries("action", "delete").Path(hyphaUrl).
		HandlerFunc(HandlerDelete)

	r.Queries("action", "rename", "to", fs.HyphaPattern).Path(hyphaUrl).
		HandlerFunc(HandlerRename)

	r.Queries("action", "update").Path(hyphaUrl).
		HandlerFunc(HandlerUpdate)

	r.HandleFunc(hyphaUrl, HandlerView)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Welcome to MycorrhizaWiki, feel free to do anything")
	})

	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
