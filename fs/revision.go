package fs

import (
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	"log"
	"net/http"
)

type Revision struct {
	Id         int
	Tags       []string `json:"tags"`
	Name       string   `json:"name"`
	Comment    string   `json:"comment"`
	Author     string   `json:"author"`
	Time       int      `json:"time"`
	TextMime   string   `json:"text_mime"`
	BinaryMime string   `json:"binary_mime"`
	TextPath   string
	BinaryPath string
}

// During initialisation, it is guaranteed that r.BinaryMime is set to "" if the revision has no binary data.
func (r *Revision) hasBinaryData() bool {
	return r.BinaryMime != ""
}

func (r *Revision) urlOfBinary() string {
	return fmt.Sprintf("/%s?action=getBinary&rev=%d", r.Name, r.Id)
}

// TODO: use templates https://github.com/bouncepaw/mycorrhiza/issues/2
func (r *Revision) AsHtml(hyphae map[string]Hypha) (ret string, err error) {
	ret += `<article class="page">
`
	// TODO: support things other than images
	if r.hasBinaryData() {
		ret += fmt.Sprintf(`<img src="/%s" class="page__image"/>`, r.urlOfBinary())
	}

	contents, err := ioutil.ReadFile(r.TextPath)
	if err != nil {
		return "", err
	}

	// TODO: support more markups.
	// TODO: support mycorrhiza extensions like transclusion.
	switch r.TextMime {
	case "text/plain":
		ret += fmt.Sprintf(`<pre>%s</pre>`, contents)
	case "text/markdown":
		html := markdown.ToHTML(contents, nil, nil)
		ret += string(html)
	default:
		return "", errors.New("Unsupported mime-type: " + r.TextMime)
	}

	ret += `
</article>`
	return ret, nil
}

func (r *Revision) ActionGetBinary(w http.ResponseWriter) {
	fileContents, err := ioutil.ReadFile(r.urlOfBinary())
	if err != nil {
		log.Println("Failed to load binary data of", r.Name, r.Id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", r.BinaryMime)
	w.WriteHeader(http.StatusOK)
	w.Write(fileContents)
	log.Println("Serving binary data of", r.Name, r.Id)
}

func (r *Revision) ActionRaw(w http.ResponseWriter) {
	fileContents, err := ioutil.ReadFile(r.TextPath)
	if err != nil {
		log.Println("Failed to load text data of", r.Name, r.Id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", r.TextMime)
	w.WriteHeader(http.StatusOK)
	w.Write(fileContents)
	log.Println("Serving text data of", r.Name, r.Id)
}

func (r *Revision) ActionZen(w http.ResponseWriter, hyphae map[string]Hypha) {
	html, err := r.AsHtml(hyphae)
	if err != nil {
		log.Println("Failed to render", r.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, html)
}

func (r *Revision) ActionView(w http.ResponseWriter, hyphae map[string]Hypha, layoutFun func(map[string]Hypha, Revision, string) string) {
	html, err := r.AsHtml(hyphae)
	if err != nil {
		log.Println("Failed to render", r.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, layoutFun(hyphae, *r, html))
	log.Println("Rendering", r.Name)
}
