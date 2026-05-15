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
