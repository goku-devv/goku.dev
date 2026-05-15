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
	Name      string
	Role      string
	Location  string
	Image     string
	ImageDark string
	Tagline   template.HTML
	CTAs      []CTA
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
	Layout    string          `yaml:"layout"`
	Title     string          `yaml:"title"`
	Name      string          `yaml:"name"`
	Role      string          `yaml:"role"`
	Location  string          `yaml:"location"`
	Image     string          `yaml:"image"`
	ImageDark string          `yaml:"image_dark"`
	Tagline   string          `yaml:"tagline"`
	CTAs      []CTA           `yaml:"ctas"`
	Entries   []TimelineEntry `yaml:"entries"`
	Groups    []SkillGroup    `yaml:"groups"`
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
			Name:      meta.Name,
			Role:      meta.Role,
			Location:  meta.Location,
			Image:     meta.Image,
			ImageDark: meta.ImageDark,
			Tagline:   renderInline(meta.Tagline),
			CTAs:      meta.CTAs,
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

// renderInline parses a short string as markdown and strips the wrapping <p>
// tag, so the result can sit inside another element (e.g. <p class="hero-tagline">).
// Empty input returns empty.
func renderInline(s string) template.HTML {
	if s == "" {
		return ""
	}
	out := strings.TrimSpace(string(markdownToHTML([]byte(s))))
	out = strings.TrimPrefix(out, "<p>")
	out = strings.TrimSuffix(out, "</p>")
	return template.HTML(out)
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
	nl := bytes.IndexByte(rest, '\n')
	if nl < 0 {
		return nil, raw
	}
	rest = rest[nl+1:]
	end := bytes.Index(rest, []byte("\n"+delim))
	if end < 0 {
		return nil, raw
	}
	fm := rest[:end]
	body := rest[end+1+len(delim):]
	body = bytes.TrimLeft(body, "\r\n")
	return fm, body
}
