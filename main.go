package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type PageData struct {
	Content template.HTML
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

	// Read markdown file
	mdContent, err := os.ReadFile("profile.md")
	if err != nil {
		http.Error(w, "Error reading markdown file", http.StatusInternalServerError)
		log.Printf("Error reading profile.md: %v", err)
		return
	}

	// Convert markdown to HTML
	htmlContent := markdownToHTML(mdContent)

	// Parse template
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	// Prepare data
	data := PageData{
		Content: template.HTML(htmlContent),
	}

	// Render template
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Home page handler
	http.HandleFunc("/", homeHandler)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
