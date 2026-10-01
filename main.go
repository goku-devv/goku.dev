package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// lastUpdated is injected at build time: -ldflags "-X main.lastUpdated=YYYY-MM-DD" (see Makefile).
var lastUpdated string

type PageData struct {
	Sections   []Section
	Updated    string
	CSSVersion string
}

// assetVersion is the cache-busting ?v= for a static file: its mtime, so a
// deploy that only uploads the file still reaches returning visitors.
func assetVersion(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return lastUpdated
	}
	return strconv.FormatInt(fi.ModTime().Unix(), 10)
}

func markdownToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	sections, err := LoadSections("content")
	if err != nil {
		http.Error(w, "Error loading content", http.StatusInternalServerError)
		log.Printf("LoadSections: %v", err)
		return
	}

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		http.Error(w, "Error parsing templates", http.StatusInternalServerError)
		log.Printf("ParseGlob: %v", err)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout.html", PageData{Sections: sections, Updated: lastUpdated, CSSVersion: assetVersion("static/style.css")}); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("ExecuteTemplate: %v", err)
		return
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler)
	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
