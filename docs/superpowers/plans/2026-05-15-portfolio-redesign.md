# Portfolio Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the single-card layout at goku.dev with a multi-section bryllim.com-inspired portfolio (hero, about, experience, education, skills, projects, recommendations, gallery, contact), driven by per-section markdown files with optional YAML frontmatter.

**Architecture:** A new `content.go` exposes `LoadSections("content")` that globs `content/*.md`, parses optional YAML frontmatter, and renders body markdown via the existing `markdownToHTML`. The HTTP handler passes a `[]Section` slice to `templates/layout.html`, which dispatches each section to one of four partial templates (`section_hero`, `section_timeline`, `section_grouped`, `section_plain`) based on `layout:`. CSS extends the existing minimalist style sheet — same vars, same dark-mode toggle, same Inter typography — adding only what each new layout needs.

**Tech Stack:** Go 1.24 (`net/http`, `html/template`), `github.com/gomarkdown/markdown` (already in go.mod), `gopkg.in/yaml.v3` (new dep), vanilla CSS + a single inline `<script>`.

**Spec:** [docs/superpowers/specs/2026-05-15-portfolio-redesign-design.md](../specs/2026-05-15-portfolio-redesign-design.md)

---

### Task 1: Add yaml.v3 dependency and define Section types

**Files:**
- Modify: `go.mod`, `go.sum`
- Create: `content.go`

- [ ] **Step 1: Add yaml.v3 to go.mod**

Run:
```bash
go get gopkg.in/yaml.v3
```

Expected: go.mod gains a `gopkg.in/yaml.v3 v3.0.x` line; go.sum is updated.

- [ ] **Step 2: Create `content.go` with types and an empty `LoadSections`**

Create `content.go`:

```go
package main

import (
	"html/template"
)

type Section struct {
	Slug    string
	Layout  string // "hero" | "timeline" | "grouped" | "plain"
	Title   string
	HTML    template.HTML
	Hero    *HeroData
	Entries []TimelineEntry
	Groups  []SkillGroup
}

type HeroData struct {
	Name     string
	Role     string
	Location string
	Image    string
	Tagline  string
	CTAs     []CTA
}

type CTA struct {
	Label string
	URL   string
}

type TimelineEntry struct {
	Company string
	Dates   string
	Role    string
	URL     string
	Bullets []string
}

type SkillGroup struct {
	Name  string
	Items []string
}

// LoadSections is implemented in Task 3.
func LoadSections(dir string) ([]Section, error) {
	return nil, nil
}
```

- [ ] **Step 3: Verify it compiles**

Run:
```bash
go build ./...
```

Expected: builds with no output.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum content.go
git commit -m "add yaml.v3 dep and Section types"
```

---

### Task 2: Write failing tests for `LoadSections`

**Files:**
- Create: `content_test.go`

- [ ] **Step 1: Write the test file**

Create `content_test.go`:

```go
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0644); err != nil {
		t.Fatalf("writeFile %s: %v", name, err)
	}
}

func TestLoadSections_MissingDir(t *testing.T) {
	sections, err := LoadSections(filepath.Join(t.TempDir(), "does-not-exist"))
	if err != nil {
		t.Fatalf("want nil err, got %v", err)
	}
	if len(sections) != 0 {
		t.Fatalf("want empty, got %d sections", len(sections))
	}
}

func TestLoadSections_OrderByFilename(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "02-about.md", "## About\n\nHello.")
	writeFile(t, dir, "01-hero.md", "---\nlayout: hero\nname: goku\nrole: Software Engineer\nlocation: HCMC\nimage: /static/x.png\ntagline: hi\n---\n")
	writeFile(t, dir, "09-contact.md", "## Contact\n\nhi@x.com")

	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 3 {
		t.Fatalf("want 3 sections, got %d", len(sections))
	}
	if sections[0].Slug != "hero" || sections[1].Slug != "about" || sections[2].Slug != "contact" {
		t.Fatalf("wrong order: %s, %s, %s", sections[0].Slug, sections[1].Slug, sections[2].Slug)
	}
}

func TestLoadSections_Hero(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "01-hero.md", `---
layout: hero
name: goku
role: Software Engineer
location: Ho Chi Minh, Viet Nam
image: /static/goku.png
tagline: Backend dev.
ctas:
  - label: Email
    url: mailto:test@example.com
  - label: GitHub
    url: https://github.com/x
---
ignored body
`)
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Layout != "hero" {
		t.Fatalf("want one hero, got %+v", sections)
	}
	h := sections[0].Hero
	if h == nil {
		t.Fatal("hero data nil")
	}
	if h.Name != "goku" || h.Role != "Software Engineer" || h.Location != "Ho Chi Minh, Viet Nam" {
		t.Errorf("hero fields: %+v", h)
	}
	if h.Image != "/static/goku.png" || h.Tagline != "Backend dev." {
		t.Errorf("hero image/tagline: %+v", h)
	}
	if len(h.CTAs) != 2 || h.CTAs[0].Label != "Email" || h.CTAs[1].URL != "https://github.com/x" {
		t.Errorf("hero CTAs: %+v", h.CTAs)
	}
}

func TestLoadSections_Timeline(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "03-experience.md", `---
title: Experience
layout: timeline
entries:
  - company: Autonomous Inc.
    dates: Dec 2020 - Present
    role: Senior Software Engineer
    url: https://autonomous.ai
    bullets:
      - Did one thing.
      - Did another thing.
  - company: Old Co
    dates: 2018 - 2020
    role: Engineer
    bullets:
      - Worked there.
---
`)
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Layout != "timeline" {
		t.Fatalf("want one timeline, got %+v", sections)
	}
	s := sections[0]
	if s.Title != "Experience" {
		t.Errorf("title: %q", s.Title)
	}
	if len(s.Entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(s.Entries))
	}
	if s.Entries[0].Company != "Autonomous Inc." || s.Entries[0].URL != "https://autonomous.ai" {
		t.Errorf("entry[0]: %+v", s.Entries[0])
	}
	if len(s.Entries[0].Bullets) != 2 || s.Entries[0].Bullets[1] != "Did another thing." {
		t.Errorf("entry[0].Bullets: %v", s.Entries[0].Bullets)
	}
	if s.Entries[1].URL != "" {
		t.Errorf("entry[1] URL should be empty, got %q", s.Entries[1].URL)
	}
}

func TestLoadSections_Grouped(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "05-skills.md", `---
title: Skills
layout: grouped
groups:
  - name: Languages
    items: [Golang]
  - name: Databases
    items: [MySQL, Redis, MongoDB]
---
`)
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Layout != "grouped" {
		t.Fatalf("want one grouped, got %+v", sections)
	}
	s := sections[0]
	if s.Title != "Skills" || len(s.Groups) != 2 {
		t.Fatalf("title/groups: %q, %d", s.Title, len(s.Groups))
	}
	if s.Groups[1].Name != "Databases" || len(s.Groups[1].Items) != 3 || s.Groups[1].Items[2] != "MongoDB" {
		t.Errorf("groups[1]: %+v", s.Groups[1])
	}
}

func TestLoadSections_Plain(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "02-about.md", "## About\n\nI build things.")
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Layout != "plain" {
		t.Fatalf("want one plain, got %+v", sections)
	}
	html := string(sections[0].HTML)
	if !strings.Contains(html, "<h2") || !strings.Contains(html, "About") {
		t.Errorf("expected rendered h2 About, got %q", html)
	}
	if !strings.Contains(html, "I build things.") {
		t.Errorf("expected body text, got %q", html)
	}
}

func TestLoadSections_PlainWithExplicitFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "06-projects.md", "---\nlayout: plain\n---\n## Projects\n")
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Layout != "plain" {
		t.Fatalf("want one plain, got %+v", sections)
	}
}

func TestLoadSections_SkipsBadYAML(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "01-broken.md", "---\nlayout: hero\nname: [unterminated\n---\n")
	writeFile(t, dir, "02-ok.md", "## OK\n\nbody")
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Slug != "ok" {
		t.Fatalf("want only the ok section, got %+v", sections)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run:
```bash
go test ./...
```

Expected: every test that probes `LoadSections` behavior fails (the stub returns `nil, nil`), most commonly with `want N sections, got 0`. `TestLoadSections_MissingDir` is the only one that may already pass — that is fine.

- [ ] **Step 3: Commit**

```bash
git add content_test.go
git commit -m "add failing tests for LoadSections"
```

---

### Task 3: Implement `LoadSections`

**Files:**
- Modify: `content.go`

- [ ] **Step 1: Replace the stub with the real implementation**

Replace the body of `LoadSections` and add the helpers below. The full `content.go` after this edit:

```go
package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Section struct {
	Slug    string
	Layout  string
	Title   string
	HTML    template.HTML
	Hero    *HeroData
	Entries []TimelineEntry
	Groups  []SkillGroup
}

type HeroData struct {
	Name     string
	Role     string
	Location string
	Image    string
	Tagline  string
	CTAs     []CTA
}

type CTA struct {
	Label string
	URL   string
}

type TimelineEntry struct {
	Company string
	Dates   string
	Role    string
	URL     string
	Bullets []string
}

type SkillGroup struct {
	Name  string
	Items []string
}

type frontmatter struct {
	Layout  string          `yaml:"layout"`
	Title   string          `yaml:"title"`
	Name    string          `yaml:"name"`
	Role    string          `yaml:"role"`
	Location string         `yaml:"location"`
	Image   string          `yaml:"image"`
	Tagline string          `yaml:"tagline"`
	CTAs    []CTA           `yaml:"ctas"`
	Entries []TimelineEntry `yaml:"entries"`
	Groups  []SkillGroup    `yaml:"groups"`
}

func LoadSections(dir string) ([]Section, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Printf("content dir %q missing; rendering empty page", dir)
			return nil, nil
		}
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)

	var out []Section
	for _, name := range names {
		path := filepath.Join(dir, name)
		raw, err := os.ReadFile(path)
		if err != nil {
			log.Printf("read %s: %v; skipping", name, err)
			continue
		}
		s, err := parseSection(name, raw)
		if err != nil {
			log.Printf("parse %s: %v; skipping", name, err)
			continue
		}
		out = append(out, s)
	}
	return out, nil
}

func parseSection(filename string, raw []byte) (Section, error) {
	slug := slugFromFilename(filename)
	fm, body := splitFrontmatter(raw)

	if len(fm) == 0 {
		return Section{
			Slug:   slug,
			Layout: "plain",
			HTML:   template.HTML(markdownToHTML(body)),
		}, nil
	}

	var meta frontmatter
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return Section{}, err
	}

	s := Section{Slug: slug, Title: meta.Title}
	switch meta.Layout {
	case "hero":
		s.Layout = "hero"
		s.Hero = &HeroData{
			Name:     meta.Name,
			Role:     meta.Role,
			Location: meta.Location,
			Image:    meta.Image,
			Tagline:  meta.Tagline,
			CTAs:     meta.CTAs,
		}
	case "timeline":
		s.Layout = "timeline"
		s.Entries = meta.Entries
	case "grouped":
		s.Layout = "grouped"
		s.Groups = meta.Groups
	default:
		s.Layout = "plain"
		s.HTML = template.HTML(markdownToHTML(body))
	}
	return s, nil
}

func slugFromFilename(name string) string {
	base := strings.TrimSuffix(name, ".md")
	if i := strings.Index(base, "-"); i > 0 && allDigits(base[:i]) {
		return base[i+1:]
	}
	return base
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// splitFrontmatter returns (frontmatterYAML, body). If the file does not start
// with a `---` line, frontmatter is empty and the whole input is body.
func splitFrontmatter(raw []byte) ([]byte, []byte) {
	const delim = "---"
	if !bytes.HasPrefix(raw, []byte(delim)) {
		return nil, raw
	}
	rest := raw[len(delim):]
	// Require a newline after the opening ---
	nl := bytes.IndexByte(rest, '\n')
	if nl < 0 {
		return nil, raw
	}
	rest = rest[nl+1:]
	// Find closing --- on its own line
	end := bytes.Index(rest, []byte("\n"+delim))
	if end < 0 {
		return nil, raw
	}
	fm := rest[:end]
	body := rest[end+1+len(delim):]
	body = bytes.TrimLeft(body, "\r\n")
	return fm, body
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run:
```bash
go test ./...
```

Expected: all tests in `content_test.go` pass.

- [ ] **Step 3: Commit**

```bash
git add content.go
git commit -m "implement LoadSections frontmatter parser"
```

---

### Task 4: Wire `LoadSections` into the HTTP handler

**Files:**
- Modify: `main.go`

- [ ] **Step 1: Replace `homeHandler` and `PageData`**

The full `main.go` after edit:

```go
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
```

- [ ] **Step 2: Verify it builds**

Run:
```bash
go build ./...
```

Expected: builds with no output. (Page won't render fully yet — partials don't exist.)

- [ ] **Step 3: Commit**

```bash
git add main.go
git commit -m "load sections in homeHandler and parse all templates"
```

---

### Task 5: Update `layout.html` to dispatch sections + add `section_plain` partial

**Files:**
- Modify: `templates/layout.html`
- Create: `templates/section_plain.html`

- [ ] **Step 1: Create `templates/section_plain.html`**

```html
{{ define "section_plain" }}
<section class="section section-plain section-{{ .Slug }}">
    {{ .HTML }}
</section>
{{ end }}
```

- [ ] **Step 2: Replace the `<main>` block in `templates/layout.html`**

Find the existing block:

```html
        <main class="profile-main">
            <section class="profile-card">
                <div class="profile-img-wrapper">
                    <img src="/static/goku.png" alt="Description">
                </div>
                <div class="profile-content">
                    {{.Content}}
                </div>
            </section>
        </main>
```

Replace with:

```html
        <main class="profile-main">
            {{ range .Sections }}
                {{ if eq .Layout "hero" }}{{ template "section_hero" . }}
                {{ else if eq .Layout "timeline" }}{{ template "section_timeline" . }}
                {{ else if eq .Layout "grouped" }}{{ template "section_grouped" . }}
                {{ else }}{{ template "section_plain" . }}
                {{ end }}
            {{ end }}
        </main>
```

- [ ] **Step 3: Update the social-icon script scope**

In the same file, find:

```js
            document.querySelectorAll('.profile-content a').forEach(function(link) {
```

Replace with:

```js
            document.querySelectorAll('main a').forEach(function(link) {
```

- [ ] **Step 4: Add the three not-yet-existing partials as empty define blocks**

`html/template`'s `ParseGlob` will fail at runtime if `layout.html` references a template that doesn't exist. Create temporary empty stubs so the page renders during incremental work:

Create `templates/section_hero.html`:
```html
{{ define "section_hero" }}{{ end }}
```

Create `templates/section_timeline.html`:
```html
{{ define "section_timeline" }}{{ end }}
```

Create `templates/section_grouped.html`:
```html
{{ define "section_grouped" }}{{ end }}
```

These are replaced with real markup in Tasks 6–8.

- [ ] **Step 5: Smoke-test in the browser**

Create a quick test file so there is something to render:

```bash
mkdir -p content
printf '## Hello\n\nIt works.\n' > content/02-test.md
```

Start the server (assume port 8080 is free — kill any existing dev server first):

```bash
go run main.go
```

In another shell:

```bash
curl -s http://localhost:8080/ | grep -F 'It works.'
```

Expected: the line `<p>It works.</p>` appears in the output. Stop the server with Ctrl-C.

Remove the test file:

```bash
rm content/02-test.md
```

- [ ] **Step 6: Commit**

```bash
git add templates/layout.html templates/section_plain.html templates/section_hero.html templates/section_timeline.html templates/section_grouped.html
git commit -m "dispatch sections from layout.html, add plain partial"
```

---

### Task 6: Implement `section_hero` partial + hero CSS

**Files:**
- Modify: `templates/section_hero.html`
- Modify: `static/style.css`

- [ ] **Step 1: Replace `templates/section_hero.html`**

```html
{{ define "section_hero" }}
{{ with .Hero }}
<section class="section section-hero">
    <div class="hero-img-wrapper">
        <img src="{{ .Image }}" alt="{{ .Name }}">
    </div>
    <h1 class="hero-name">{{ .Name }}</h1>
    <p class="hero-meta">{{ .Role }} · {{ .Location }}</p>
    {{ if .Tagline }}<p class="hero-tagline">{{ .Tagline }}</p>{{ end }}
    {{ if .CTAs }}
    <p class="hero-ctas">
        {{ range $i, $c := .CTAs }}{{ if $i }} {{ end }}<a href="{{ $c.URL }}">{{ $c.Label }}</a>{{ end }}
    </p>
    {{ end }}
</section>
{{ end }}
{{ end }}
```

- [ ] **Step 2: Append hero styles to `static/style.css`**

Add at the end of the file:

```css
/* === Section shell === */
.section {
    background: var(--card-bg);
    border: 1px solid var(--border);
    padding: 48px 40px;
    width: 100%;
    box-sizing: border-box;
    transition: background 0.2s ease, border-color 0.2s ease;
}
.profile-main {
    flex-direction: column;
    gap: 24px;
    max-width: 560px;
    margin: 0 auto;
}

/* === Hero === */
.section-hero {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 60px 40px;
}
.hero-img-wrapper {
    width: 100px;
    height: 100px;
    margin-bottom: 28px;
    border: 1px solid var(--border);
    overflow: hidden;
}
.hero-img-wrapper img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
}
.hero-name {
    font-size: 1.8rem;
    font-weight: 600;
    letter-spacing: -0.02em;
    margin: 0 0 0.5em;
    color: var(--text-primary);
}
.hero-meta {
    font-size: 0.95rem;
    color: var(--text-secondary);
    margin: 0 0 1em;
}
.hero-tagline {
    font-size: 0.95rem;
    line-height: 1.6;
    color: var(--text-primary);
    margin: 0 0 1.5em;
}
.hero-ctas {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 16px;
    margin: 0;
}
.hero-ctas a {
    color: var(--text-primary);
    text-decoration: none;
    border-bottom: 1px solid var(--link-border);
    transition: opacity 0.2s;
    font-size: 0.95rem;
}
.hero-ctas a:hover { opacity: 0.6; }

@media (max-width: 600px) {
    .section {
        padding: 36px 24px;
        border-left: none;
        border-right: none;
    }
    .section-hero { padding: 40px 24px; }
    .hero-img-wrapper { width: 80px; height: 80px; margin-bottom: 24px; }
    .hero-name { font-size: 1.5rem; }
}
```

Note: this section coexists with the existing `.profile-main`, `.profile-card`, `.profile-img-wrapper`, `.profile-content` styles. Those older selectors are no longer referenced from any template after Task 5 and can be removed in Task 9. Leaving them in place during this task keeps the diff focused.

- [ ] **Step 3: Smoke-test**

Create a minimal hero file to verify:

```bash
mkdir -p content
cat > content/01-hero.md <<'EOF'
---
layout: hero
name: goku
role: Software Engineer
location: Ho Chi Minh, Viet Nam
image: /static/goku.png
tagline: Backend developer with 7 years building high-traffic distributed systems.
ctas:
  - label: hi.im@goku.dev
    url: mailto:hi.im@goku.dev
  - label: GitHub
    url: https://github.com/goku-devv
  - label: X
    url: https://x.com/goku_dev
---
EOF
```

Start the server and open http://localhost:8080 in a browser. Confirm:
- Photo, name, "Software Engineer · Ho Chi Minh, Viet Nam", tagline, and three underlined CTA links visible (the GitHub and X labels swap to icon + text via the existing inline script)
- Light/dark toggle still works

Stop the server.

- [ ] **Step 4: Commit**

```bash
git add templates/section_hero.html static/style.css content/01-hero.md
git commit -m "implement hero section and styles"
```

---

### Task 7: Implement `section_timeline` partial + timeline CSS

**Files:**
- Modify: `templates/section_timeline.html`
- Modify: `static/style.css`

- [ ] **Step 1: Replace `templates/section_timeline.html`**

```html
{{ define "section_timeline" }}
<section class="section section-timeline section-{{ .Slug }}">
    <h2 class="section-title">{{ .Title }}</h2>
    <div class="timeline">
        {{ range .Entries }}
        <div class="timeline-row">
            <div class="timeline-dates">{{ .Dates }}</div>
            <div class="timeline-entry">
                <div class="timeline-company">
                    {{ if .URL }}<a href="{{ .URL }}">{{ .Company }}</a>{{ else }}{{ .Company }}{{ end }}
                </div>
                {{ if .Role }}<div class="timeline-role">{{ .Role }}</div>{{ end }}
                {{ if .Bullets }}
                <ul class="timeline-bullets">
                    {{ range .Bullets }}<li>{{ . }}</li>{{ end }}
                </ul>
                {{ end }}
            </div>
        </div>
        {{ end }}
    </div>
</section>
{{ end }}
```

- [ ] **Step 2: Append timeline styles to `static/style.css`**

```css
/* === Section title (shared by timeline + grouped) === */
.section-title {
    font-size: 0.85rem;
    color: var(--text-primary);
    margin: 0 0 24px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

/* === Timeline === */
.timeline {
    display: grid;
    grid-template-columns: 110px 1fr;
    column-gap: 20px;
    row-gap: 28px;
}
.timeline-row {
    display: contents;
}
.timeline-dates {
    font-size: 0.85rem;
    color: var(--text-secondary);
    padding-top: 2px;
}
.timeline-entry {
    position: relative;
    padding-left: 20px;
    border-left: 1px solid var(--border);
}
.timeline-entry::before {
    content: "";
    position: absolute;
    left: -4px;
    top: 8px;
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: var(--text-primary);
}
.timeline-company {
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--text-primary);
}
.timeline-company a {
    color: var(--text-primary);
    text-decoration: none;
    border-bottom: 1px solid var(--link-border);
}
.timeline-company a:hover { opacity: 0.6; }
.timeline-role {
    font-size: 0.9rem;
    color: var(--text-secondary);
    font-style: italic;
    margin-top: 2px;
}
.timeline-bullets {
    margin: 10px 0 0;
    padding-left: 18px;
}
.timeline-bullets li {
    color: var(--text-secondary);
    font-size: 0.95rem;
    line-height: 1.7;
    margin-bottom: 4px;
}

@media (max-width: 600px) {
    .timeline { grid-template-columns: 1fr; row-gap: 24px; }
    .timeline-entry { padding-left: 0; border-left: none; }
    .timeline-entry::before { display: none; }
    .timeline-dates { padding-top: 0; margin-bottom: 4px; }
}
```

- [ ] **Step 3: Smoke-test**

Create the experience file (use this minimal version; full content lands in Task 9):

```bash
cat > content/03-experience.md <<'EOF'
---
title: Experience
layout: timeline
entries:
  - company: Autonomous Inc.
    dates: Dec 2020 – Present
    role: Senior Software Engineer
    url: https://autonomous.ai
    bullets:
      - Maintain a high-traffic monolithic eCommerce system.
      - Led migration to microservices architecture.
EOF
```

Start the server and visit http://localhost:8080. Confirm:
- "EXPERIENCE" heading
- Dates on the left, company linked to autonomous.ai on the right, italic role, two bullets
- On mobile width (Chrome devtools, 375px), single-column stack with no rail/dot

Stop the server.

- [ ] **Step 4: Commit**

```bash
git add templates/section_timeline.html static/style.css content/03-experience.md
git commit -m "implement timeline section and styles"
```

---

### Task 8: Implement `section_grouped` partial + grouped CSS

**Files:**
- Modify: `templates/section_grouped.html`
- Modify: `static/style.css`

- [ ] **Step 1: Replace `templates/section_grouped.html`**

```html
{{ define "section_grouped" }}
<section class="section section-grouped section-{{ .Slug }}">
    <h2 class="section-title">{{ .Title }}</h2>
    <div class="grouped">
        {{ range .Groups }}
        <div class="grouped-row">
            <div class="grouped-name">{{ .Name }}</div>
            <div class="grouped-items">
                {{ range $i, $item := .Items }}{{ if $i }} · {{ end }}{{ $item }}{{ end }}
            </div>
        </div>
        {{ end }}
    </div>
</section>
{{ end }}
```

- [ ] **Step 2: Append grouped styles to `static/style.css`**

```css
/* === Grouped (Skills) === */
.grouped {
    display: grid;
    grid-template-columns: 110px 1fr;
    column-gap: 20px;
    row-gap: 12px;
}
.grouped-row {
    display: contents;
}
.grouped-name {
    font-size: 0.75rem;
    color: var(--text-secondary);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    padding-top: 3px;
}
.grouped-items {
    font-size: 0.95rem;
    color: var(--text-primary);
    line-height: 1.7;
}

@media (max-width: 600px) {
    .grouped { grid-template-columns: 1fr; row-gap: 16px; }
    .grouped-name { padding-top: 0; }
}
```

- [ ] **Step 3: Smoke-test**

```bash
cat > content/05-skills.md <<'EOF'
---
title: Skills
layout: grouped
groups:
  - name: Languages
    items: [Golang]
  - name: Databases
    items: [MySQL, Redis, MongoDB, Elasticsearch]
EOF
```

Start the server. Confirm "SKILLS" heading, two rows with group name on left and dot-separated items on right. Mobile width stacks correctly. Stop the server.

- [ ] **Step 4: Commit**

```bash
git add templates/section_grouped.html static/style.css content/05-skills.md
git commit -m "implement grouped section and styles"
```

---

### Task 9: Style plain sections + remove legacy `.profile-*` selectors

**Files:**
- Modify: `static/style.css`

The plain partial already exists and renders. This task makes it look right inside the new section card, and removes the now-unused `.profile-*` selectors that were leftover from the single-card layout.

- [ ] **Step 1: Remove the legacy `.profile-card`, `.profile-img-wrapper`, `.profile-card img`, and all `.profile-content*` selectors**

Open `static/style.css` and delete every rule whose selector contains `.profile-card`, `.profile-img-wrapper`, or `.profile-content`. Keep `.profile-main` — it's still used as the outer wrapper. Keep the new `.section`, `.section-*`, `.hero-*`, `.timeline-*`, `.grouped-*` rules.

Also remove the old `@media (max-width: 600px)` block that targets `.profile-card`, `.profile-img-wrapper`, `.profile-content h1`, and `.profile-main` padding — the new mobile rules per section already cover the breakpoint.

- [ ] **Step 2: Add `.section-plain` typography**

Append to `static/style.css`:

```css
/* === Plain (About, Projects, Recommendations, Gallery, Contact) === */
.section-plain h1,
.section-plain h2 {
    font-size: 0.85rem;
    color: var(--text-primary);
    margin: 0 0 16px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}
.section-plain h3 {
    font-size: 1rem;
    color: var(--text-primary);
    margin: 1.5em 0 0.5em;
    font-weight: 600;
}
.section-plain p,
.section-plain li {
    color: var(--text-secondary);
    font-size: 0.95rem;
    line-height: 1.8;
    font-weight: 400;
}
.section-plain ul {
    padding-left: 18px;
    margin: 0 0 1em;
}
.section-plain a {
    color: var(--text-primary);
    text-decoration: none;
    border-bottom: 1px solid var(--link-border);
    transition: opacity 0.2s;
}
.section-plain a:hover { opacity: 0.6; }
.section-plain > :first-child { margin-top: 0; }
.section-plain > :last-child { margin-bottom: 0; }
```

- [ ] **Step 3: Verify build**

Run:
```bash
go build ./...
```

Expected: builds with no output.

- [ ] **Step 4: Commit**

```bash
git add static/style.css
git commit -m "style plain sections and drop legacy profile-card CSS"
```

---

### Task 10: Create remaining content files

**Files:**
- Create: `content/02-about.md`, `content/04-education.md`, `content/06-projects.md`, `content/07-recommendations.md`, `content/08-gallery.md`, `content/09-contact.md`
- Modify: `content/03-experience.md`, `content/05-skills.md` (expand from Task 7/8 smoke-test versions to full CV content)

- [ ] **Step 1: Expand `content/03-experience.md` to all four roles from the CV**

Overwrite with:

```yaml
---
title: Experience
layout: timeline
entries:
  - company: Autonomous Inc.
    dates: Dec 2020 – Present
    role: Senior Software Engineer
    url: https://autonomous.ai
    bullets:
      - Maintain and optimize a high-traffic monolithic eCommerce system (autonomous.ai) handling 1M+ orders, 1M+ shipments, and 150K daily sessions.
      - Led migration from monolithic to microservices architecture, improving scalability and fault tolerance.
      - Implemented OpenTelemetry for distributed tracing and Datadog for real-time monitoring.
      - Developed an AI-powered chatbot using Llama, FLUX, and Hermès, supporting expansion into AI-driven services.
      - Integrated multi-chain cryptocurrency payments (BTC, ETH, SOL, BNB).
      - Stack — Golang, MySQL, MongoDB, Redis, Google Pub/Sub, gRPC, Elasticsearch.
  - company: Bestarion Software Company Ltd.
    dates: Mar 2020 – Dec 2020
    role: Senior Software Engineer
    bullets:
      - Developed SC-Innovate, a Kanban workflow management system for Standard Chartered Singapore.
      - Built backend services with Golang and MongoDB, collaborating with a VueJS frontend team.
  - company: WeVenture Pte Ltd
    dates: Feb 2018 – Feb 2020
    role: Software Engineer
    bullets:
      - Developed a mobile payment integration platform with Centili, txtNation, IMImobile, Comviva, Onebip, APIGate, and Stripe.
      - Designed a microservices architecture using Golang, MySQL, RabbitMQ, gRPC, and Redis.
  - company: Gumi Vietnam
    dates: Jul 2016 – Jan 2018
    role: Software Engineer
    bullets:
      - Built a RESTful API for a hospital drug testing system using Ruby on Rails, MySQL, and AngularJS.
      - Optimized complex queries in a large-scale EAV (Entity-Attribute-Value) model.
---
```

- [ ] **Step 2: Expand `content/05-skills.md` to full CV skill list**

Overwrite with:

```yaml
---
title: Skills
layout: grouped
groups:
  - name: Languages
    items: [Golang]
  - name: Databases
    items: [MySQL, Redis, MongoDB, Elasticsearch]
  - name: Messaging
    items: [RabbitMQ, Google Pub/Sub]
  - name: Architecture
    items: [Monolithic, Microservices]
  - name: Observability
    items: [OpenTelemetry, Datadog, NewRelic, Prometheus, ELK Stack]
---
```

- [ ] **Step 3: Create `content/02-about.md`**

```markdown
## About

I'm a backend developer with 7 years of experience specializing in Go and microservices. I design and optimize large-scale distributed systems that handle high-traffic workloads, with a focus on backend performance, database optimization, and real-time observability. Always interested in complex technical challenges and improving system scalability.
```

- [ ] **Step 4: Create `content/04-education.md`**

```yaml
---
title: Education
layout: timeline
entries:
  - company: University of Technology, Ho Chi Minh City
    dates: Oct 2012 – Jul 2016
    role: Software Development
---
```

- [ ] **Step 5: Create `content/06-projects.md`**

```markdown
## Projects

_To be filled in._
```

- [ ] **Step 6: Create `content/07-recommendations.md`**

```markdown
## Recommendations

_To be filled in._
```

- [ ] **Step 7: Create `content/08-gallery.md`**

```markdown
## Gallery

_To be filled in._
```

- [ ] **Step 8: Create `content/09-contact.md`**

```markdown
## Contact

[hi.im@goku.dev](mailto:hi.im@goku.dev)
```

- [ ] **Step 9: Commit**

```bash
git add content/
git commit -m "add full portfolio content"
```

---

### Task 11: Migrate from `profile.md` and final verification

**Files:**
- Delete: `profile.md`

- [ ] **Step 1: Delete the legacy markdown file**

```bash
git rm profile.md
```

- [ ] **Step 2: Run all tests**

```bash
go test ./...
```

Expected: all `TestLoadSections_*` tests pass.

- [ ] **Step 3: Manual browser verification**

Start the server:
```bash
go run main.go
```

Open http://localhost:8080. Verify in order:

1. **Hero**: photo, "goku", "Software Engineer · Ho Chi Minh, Viet Nam", tagline, three CTA links (email + GitHub + X, with the GitHub and X labels swapped to icon + text by the inline script).
2. **About**: "ABOUT" heading and narrative paragraph.
3. **Experience**: "EXPERIENCE" heading; four entries with dates left, company/role/bullets right, vertical rail with dots, "Autonomous Inc." linked to autonomous.ai.
4. **Education**: single entry, dates left, university right.
5. **Skills**: five rows of grouped items separated by middle dots.
6. **Projects / Recommendations / Gallery**: each shows only its heading and the "_To be filled in._" italic line.
7. **Contact**: "CONTACT" heading and email link.
8. **Theme toggle**: top-right button switches light ⇄ dark; preference persists across reload (test with browser localStorage).
9. **Mobile**: resize to <600px (or use devtools 375px) — section cards go edge-to-edge, timeline/skills collapse to single column.

Stop the server.

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "remove legacy profile.md after migration"
```

---

## Self-review

- **Spec coverage:** Hero (Task 6), About (Tasks 5+10), Experience (Task 7+10), Education (Task 10), Skills (Task 8+10), Projects/Recs/Gallery placeholders (Task 10), Contact (Task 10). Server changes (Tasks 1–4). Templates (Tasks 5–8). CSS (Tasks 6–9). Migration & social-icon scope (Tasks 5, 11). Dark mode preservation: covered — Tasks 5–9 only append CSS and don't touch the `--bg`/`--card-bg`/theme-toggle code in `layout.html`.
- **Placeholder scan:** "_To be filled in._" appears only in user-facing markdown for sections the user explicitly chose to leave empty (Task 10 steps 5–7), not in plan instructions. No TBD/TODO in implementation steps.
- **Type consistency:** `Section`, `HeroData`, `CTA`, `TimelineEntry`, `SkillGroup` defined once in Task 1, used unchanged in Tasks 3, 5, 6, 7, 8. Template field names (`.Hero.Name`, `.Entries`, `.Groups`) match Go struct fields. Template names referenced in `layout.html` (Task 5) all exist as `{{ define }}` blocks (stubs in Task 5, real in Tasks 6/7/8, and `section_plain` in Task 5).
