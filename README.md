# goku.dev

Personal portfolio site. Multi-section markdown-driven layout served by a small Go HTTP server. Light/dark theme, Geist Mono typography, no build pipeline, no JS framework.

## Stack

- **Go 1.24** — `net/html/template`, `net/http`
- **[gomarkdown/markdown](https://github.com/gomarkdown/markdown)** — markdown rendering
- **[gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)** — frontmatter parsing
- Vanilla CSS + ~50 lines of inline JS (theme toggle, gallery slider, social-icon enhance)

## Project layout

```
goku.dev/
├── main.go                 HTTP server, handler, markdown renderer
├── content.go              LoadSections() — globs content/*.md, parses
│                           frontmatter, dispatches per-section types
├── content_test.go
├── content/                Each file = one section, ordered by NN- prefix
│   ├── 01-hero.md          layout: hero
│   ├── 02-about.md         plain markdown
│   ├── 03-experience.md    layout: timeline
│   ├── 04-education.md     layout: timeline
│   ├── 05-skills.md        layout: grouped
│   ├── 08-gallery.md       plain markdown + raw HTML slider
│   └── 09-contact.md       plain markdown
├── templates/
│   ├── layout.html         outer shell, theme toggle, dispatcher
│   ├── section_hero.html
│   ├── section_timeline.html
│   ├── section_grouped.html
│   └── section_plain.html
├── static/
│   ├── style.css
│   ├── favicon.png
│   └── photos/             profile + gallery images
├── scripts/
│   └── deploy.sh           gitignored — see Deploy
├── docs/superpowers/       design spec + implementation plan
├── Makefile
├── go.mod
└── go.sum
```

## Section model

Every file in `content/` becomes one `<section>` on the page. Files without YAML frontmatter render as plain markdown. Files with frontmatter dispatch on `layout:`:

| layout     | Used for                | Frontmatter shape                                                      |
| ---------- | ----------------------- | ---------------------------------------------------------------------- |
| `hero`     | Top of page             | `name`, `role`, `location`, `image`, `tagline`, `ctas: [{label, url}]` |
| `timeline` | Experience, Education   | `title`, `entries: [{company, dates, role, url, bullets}]`             |
| `grouped`  | Skills                  | `title`, `groups: [{name, items}]`                                     |
| _omitted_  | About, Gallery, Contact | none — body rendered as markdown                                       |

Filenames use a numeric `NN-` prefix to control render order. Gaps are allowed — drop a file to skip a section.

Example `03-experience.md`:

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
      - Maintain a high-traffic monolithic eCommerce system.
      - Led migration to microservices architecture.
---
```

A bad YAML block in a single file is logged and that section is skipped — the rest of the page still renders.

## Quick start

```bash
git clone https://github.com/goku-devv/goku.dev.git
cd goku.dev
make run
```

Visit http://localhost:8080. Templates and content are reloaded per request, so edit `content/*.md`, `templates/*.html`, or `static/style.css` and reload — no restart needed.

## Make targets

```bash
make build         # local binary -> ./profilepage
make build-linux   # cross-compile -> ./profilepage (linux/amd64, static, stripped)
make run           # go run .
make test          # go test ./...
make clean         # remove binaries
```

## Theme

Light/dark theme is driven by a `data-theme` attribute on `<html>`. An inline `<script>` in `<head>` reads `localStorage.theme` (falling back to `prefers-color-scheme`) and sets it before paint, preventing FOUC. The toggle button in the top-right flips and persists.

CSS uses six variables (`--bg`, `--card-bg`, `--border`, `--text-primary`, `--text-secondary`, `--link-border`) overridden under `[data-theme="dark"]`.

## Author

- Email: hi.im@goku.dev
- X: [@goku_dev](https://x.com/goku_dev)
- GitHub: [@goku-devv](https://github.com/goku-devv)

## License

MIT.
