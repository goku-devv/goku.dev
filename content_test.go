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
image: /static/photos/goku.png
tagline: Backend dev.
summary: Go backends.
ctas:
  - label: Email
    url: mailto:test@example.com
    primary: true
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
	if h.Image != "/static/photos/goku.png" || h.Tagline != "Backend dev." {
		t.Errorf("hero image/tagline: %+v", h)
	}
	if h.Summary != "Go backends." {
		t.Errorf("hero summary: %q", h.Summary)
	}
	if len(h.CTAs) != 2 || h.CTAs[0].Label != "Email" || h.CTAs[1].URL != "https://github.com/x" {
		t.Errorf("hero CTAs: %+v", h.CTAs)
	}
	if !h.CTAs[0].Primary || h.CTAs[1].Primary {
		t.Errorf("hero CTA primary flags: %+v", h.CTAs)
	}
}

func TestLoadSections_Timeline(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "03-experience.md", `---
title: Experience
layout: timeline
entries:
  - company: Autonomous Inc
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
	if s.Entries[0].Company != "Autonomous Inc" || s.Entries[0].URL != "https://autonomous.ai" {
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

// macOS tar/cp can leave AppleDouble files ("._01-hero.md") next to real
// content; they must never render as sections.
func TestLoadSections_SkipsHiddenFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "._01-hero.md", "\x00\x05\x16\x07Mac OS X ATTR")
	writeFile(t, dir, "01-hero.md", "---\nlayout: hero\nname: goku\n---\n")
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Slug != "hero" {
		t.Fatalf("want only the hero section, got %+v", sections)
	}
}

// A plain section with `collapse: true` renders each ### case as a closed
// <details>: the heading and its kicker line become the <summary>.
func TestLoadSections_CollapsibleCases(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "03-work.md", "---\ncollapse: true\n---\n## Work\n\nIntro.\n\n### One\n\n*kicker one*\n\nBody **one**.\n\n### Two\n\nBody two.\n")
	sections, err := LoadSections(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 || sections[0].Layout != "plain" {
		t.Fatalf("want one plain section, got %+v", sections)
	}
	got := string(sections[0].HTML)
	for _, want := range []string{
		"<h2>Work</h2>",
		"<p>Intro.</p>",
		"<details class=\"case\">\n<summary><h3>One</h3>\n\n<p><em>kicker one</em></p></summary>",
		"<p>Body <strong>one</strong>.</p>",
		"<details class=\"case\">\n<summary><h3>Two</h3></summary>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in:\n%s", want, got)
		}
	}
	if n := strings.Count(got, "</details>"); n != 2 {
		t.Errorf("want 2 closed details, got %d in:\n%s", n, got)
	}
	if strings.Index(got, "<p>Intro.</p>") > strings.Index(got, "<details") {
		t.Errorf("intro must stay outside the cases:\n%s", got)
	}
}
