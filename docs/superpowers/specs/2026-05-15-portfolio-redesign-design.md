# Portfolio Redesign — Design

Bryllim.com-inspired multi-section portfolio for goku.dev, evolving the existing minimalist single-card layout. Content is split across multiple markdown files, with light YAML frontmatter on sections that benefit from typed rendering (hero, experience, education, skills).

## Goals

- Replace the single profile card with a multi-section vertical layout: Hero, About, Experience, Education, Skills, Projects, Recommendations, Gallery, Contact.
- Preserve the current minimalist aesthetic — thin 1px borders, narrow centered column, ample whitespace, Inter font, existing light/dark theme toggle.
- Use real CV-derived content for Experience, Education, Skills. Projects, Recommendations, Gallery render as empty placeholder sections to be filled later.
- Keep the existing pseudonym "goku" on the public site; do not surface the CV's real name.

## Non-goals

- No nav bar / scroll-spy / anchor menu.
- No animations, scroll-snap, or scroll-triggered effects.
- No client-side framework or build pipeline. Plain Go templates + CSS + the existing single inline JS block.
- No structured-data alternative to markdown (e.g. YAML-only or hard-coded Go struct). Already decided: multi-file markdown with optional frontmatter.

## Content layout

New `content/` directory at the repo root:

```
content/
  01-hero.md             typed: hero
  02-about.md            plain markdown
  03-experience.md       typed: timeline
  04-education.md        typed: timeline
  05-skills.md           typed: grouped
  06-projects.md         plain markdown (placeholder)
  07-recommendations.md  plain markdown (placeholder)
  08-gallery.md          plain markdown (placeholder)
  09-contact.md          plain markdown
```

- `NN-` prefix determines render order. Gaps in numbering are allowed.
- Missing files are silently skipped — placeholder sections may be absent or empty.
- The existing `profile.md` at repo root is deleted after content migration.

### Frontmatter schemas

**Hero (`01-hero.md`)**

```yaml
---
layout: hero
name: goku
role: Software Engineer
location: Ho Chi Minh, Viet Nam
image: /static/photos/goku.png
tagline: Backend developer with 7 years building high-traffic distributed systems.
ctas:
  - label: hi.im@goku.dev
    url: mailto:hi.im@goku.dev
  - label: GitHub
    url: https://github.com/goku-devv
  - label: X
    url: https://x.com/goku_dev
---
```

Body is ignored for `layout: hero`.

**Timeline (`03-experience.md`, `04-education.md`)**

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
      - Maintain and optimize a high-traffic monolithic eCommerce system (autonomous.ai)…
      - Led migration from monolithic to microservices architecture…
      - Implemented OpenTelemetry for distributed tracing and Datadog for monitoring.
      - Developed an AI-powered chatbot using Llama, FLUX, and Hermès.
      - Integrated multi-chain crypto payments (BTC, ETH, SOL, BNB).
  - company: Bestarion Software Company Ltd.
    dates: Mar 2020 – Dec 2020
    role: Senior Software Engineer
    bullets:
      - Developed SC-Innovate, a Kanban workflow management system for Standard Chartered Singapore.
      - Built backend services with Golang and MongoDB, collaborating with a VueJS frontend team.
  # … WeVenture, Gumi VietNam
---
```

`url` is optional. Body is ignored.

**Grouped (`05-skills.md`)**

```yaml
---
title: Skills
layout: grouped
groups:
  - name: Languages
    items: [Golang]
  - name: Databases
    items: [MySQL, Redis, MongoDB, Elasticsearch]
  - name: Messaging & Streaming
    items: [RabbitMQ, Google Pub/Sub]
  - name: Architecture
    items: [Monolithic, Microservices]
  - name: Observability
    items: [OpenTelemetry, NewRelic, Datadog, Prometheus, ELK Stack]
---
```

Body is ignored.

**Plain (`02-about.md`, `06-projects.md`, `07-recommendations.md`, `08-gallery.md`, `09-contact.md`)**

No frontmatter required. The entire file is rendered as markdown and emitted verbatim inside `<section class="section-plain">`. There is no separate `Title` field for plain sections — authors put their own `# Heading` or `## Heading` at the top of the file as they like. If frontmatter is present, `layout: plain` is the implicit default; any unknown frontmatter keys are ignored.

## Server changes

### New file: `content.go`

```go
type Section struct {
    Slug    string        // derived from filename without NN- prefix and .md suffix
    Layout  string        // "hero" | "timeline" | "grouped" | "plain"
    Title   string        // section heading; set for timeline/grouped, unused for hero/plain
    HTML    template.HTML // rendered body, used only for layout: plain
    Hero    *HeroData
    Entries []TimelineEntry
    Groups  []SkillGroup
}

type HeroData struct {
    Name, Role, Location, Image, Tagline string
    CTAs []CTA
}
type CTA struct{ Label, URL string }

type TimelineEntry struct {
    Company, Dates, Role, URL string
    Bullets []string
}

type SkillGroup struct {
    Name  string
    Items []string
}

func LoadSections(dir string) ([]Section, error)
```

- Globs `dir + "/*.md"`, sorts by filename (lexicographic — the `NN-` prefix gives the order).
- For each file: split frontmatter (between leading `---` lines) from body. Parse frontmatter with `gopkg.in/yaml.v3`; render body with the existing `markdownToHTML`.
- Dispatch on `layout:` to populate the corresponding typed field. Unknown or missing `layout:` → `plain` with `HTML` set.
- A parse error on a single section logs a warning and skips that section; the page still renders the rest.
- A missing `content/` directory returns an empty slice without erroring.

### `main.go` changes

- Add `gopkg.in/yaml.v3` dependency (`go get`).
- `PageData` becomes `{ Sections []Section }`.
- `homeHandler` calls `LoadSections("content")` per request (site is tiny; no caching needed yet) and passes the result to the template.
- `tmpl.ParseFiles` extended to include the new partial templates.

## Templates

```
templates/
  layout.html             outer shell, theme toggle, ranges over .Sections
  section_hero.html       photo + name + role/location + tagline + CTA row
  section_timeline.html   date column, divider, role+company+bullets
  section_grouped.html    group name + inline-dot-separated items per row
  section_plain.html      raw {{ .HTML }} inside <section>
```

`layout.html` dispatches via:

```
{{ range .Sections }}
  {{ if eq .Layout "hero" }}      {{ template "section_hero" . }}
  {{ else if eq .Layout "timeline" }} {{ template "section_timeline" . }}
  {{ else if eq .Layout "grouped" }}  {{ template "section_grouped" . }}
  {{ else }}                       {{ template "section_plain" . }}
  {{ end }}
{{ end }}
```

The existing inline `<script>` that swaps social-link text for SVG icons currently scopes itself to `.profile-content a`. It is updated to scope to anchors anywhere inside `<main>` (i.e. any rendered section), since social links now live in `01-hero.md` CTAs and `09-contact.md`. The theme-toggle script and pre-paint `data-theme` script stay unchanged.

## CSS

Extends the existing `static/style.css`. Keeps the `--bg`, `--card-bg`, `--border`, `--text-primary`, `--text-secondary`, `--link-border` variables and their `[data-theme="dark"]` overrides from the previous change. No new fonts.

### Page rhythm

- `.profile-main`: existing centered flex column, gap `24px` between sections, `max-width: 560px` (was 480 — widened to fit two-column timeline).
- `.section`: 1px border, `--card-bg`, padding `40px 32px`, same border-radius as current card (none).

### Per-layout styles

- **Hero** — photo wrapper unchanged (100px, 1px border). Name `h1` 1.8rem weight 600 centered. Role · location line 0.95rem secondary color. Tagline 0.95rem primary color. CTA row: inline text links separated by no special separator, each underlined like the existing `.profile-content a`.

- **Timeline** — `display: grid; grid-template-columns: 110px 1fr; column-gap: 20px; row-gap: 28px;`. Left column: dates 0.85rem secondary. Right column: company name 0.95rem weight 600 (linked if `url` set, with the existing underline style); role italic 0.9rem secondary; `<ul>` bullets reusing existing body type. A 1px `border-left` on the right column acts as the timeline rail, with a small dot (`::before`, 6px circle, `--text-primary`) at the start of each entry.

- **Grouped** — `display: grid; grid-template-columns: 110px 1fr; row-gap: 12px; column-gap: 20px;`. Left: group name 0.75rem uppercase letter-spacing 0.05em secondary. Right: items joined by `·` (middle dot) at 0.95rem.

- **Plain** — inherits existing `.profile-content` typography rules. Empty body = section renders only its `<h2>` heading (acts as a labeled placeholder).

### Mobile (≤600px)

- Section cards lose left/right borders edge-to-edge (matches current behavior).
- Timeline and Grouped both collapse to single column: date row above company row, group name above items row. Vertical rail and dot are hidden on mobile to avoid visual clutter.
- Padding reduced to `32px 20px`.

## Migration

- Create `content/` with the nine files listed above. Real content for hero, about, experience, education, skills, contact. Empty (only `# Heading`) files for projects, recommendations, gallery.
- Delete `profile.md` from repo root.
- Run `go run main.go`; visit http://localhost:8080; verify each section renders in both light and dark themes; verify mobile breakpoint at 600px.

## Error handling

- Missing `content/` directory: empty page rendered (just `<body>` + theme toggle). Logged warning.
- Section parse error (bad YAML, bad markdown): log warning identifying the file, skip that section, render the rest.
- A `Section` with `Layout: hero` but missing `Hero` payload: section is skipped (treated as parse error).

## Testing

- No automated tests added in this round. The repo currently has no tests; this redesign isn't a good place to introduce a test harness.
- Manual verification per "Migration" above. UI feature-completeness, light/dark switching, and mobile breakpoint checked in a browser.
