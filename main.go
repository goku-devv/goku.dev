package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type PageData struct {
	Sections []Section
}

func markdownToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
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

	if err := tmpl.ExecuteTemplate(w, "layout.html", PageData{Sections: sections}); err != nil {
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
