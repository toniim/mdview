# mdview Codebase Summary

## Project Overview

**mdview** (Markdown Novel Viewer) is a lightweight, single-binary HTTP server that serves markdown files and code with a calm, book-like reading experience. It features a responsive sidebar navigation, plan-based documentation support, syntax highlighting, Mermaid diagram rendering, and a customizable theme system.

- **Module:** `github.com/bilabl/mdview`
- **Language:** Go 1.22
- **Version:** 0.1.0
- **Total LOC:** ~5,255

## Directory Structure

```
mdview/
├── cmd/                          # CLI commands
│   ├── root.go                   # Cobra root command setup
│   ├── serve.go                  # HTTP server startup logic
│   ├── stop.go                   # Stop running instances
│   └── version.go                # Version display
├── internal/
│   ├── server/                   # HTTP server & routes
│   │   ├── server.go             # Server factory, config
│   │   ├── routes.go             # 5 HTTP routes + MIME handling
│   │   └── security.go           # Path validation
│   ├── navigation/               # Plan-based navigation
│   │   ├── detect.go             # Plan.md detection
│   │   ├── parser.go             # 6-format plan parser
│   │   ├── sidebar.go            # Sidebar generation
│   │   └── footer.go             # Navigation footer
│   ├── renderer/                 # Content rendering
│   │   ├── markdown.go           # Markdown→HTML pipeline
│   │   ├── code.go               # Code file rendering
│   │   ├── toc.go                # Table of contents
│   │   ├── mermaid.go            # Mermaid preprocessing
│   │   └── images.go             # Image path resolution
│   ├── templates/                # Asset embedding & templating
│   │   ├── embed.go              # Global embed.FS
│   │   └── page.go               # Page rendering
│   └── process/                  # Process management
│       ├── pid.go                # PID file management
│       └── ports.go              # Port allocation (3456-3500)
├── assets/
│   ├── reader.js                 # Frontend logic (743 LOC)
│   ├── template.html             # HTML template
│   ├── styles/                   # CSS theme system (9 files, 1,739 LOC)
│   ├── chroma-light.css          # Syntax highlighting (light)
│   ├── chroma-dark.css           # Syntax highlighting (dark)
│   └── directory-browser.css     # Directory listing styles
├── tools/
│   └── chroma-generator.go       # CSS generator for syntax themes
├── main.go                       # Entry point (embeds assets)
├── go.mod                        # Go dependencies
├── Makefile                      # Build targets
└── docs/                         # Documentation (this directory)
```

## Core Packages

### cmd/
**Cobra CLI framework** with three commands:

| Command | File | Purpose |
|---------|------|---------|
| `serve` | serve.go (240 LOC) | Start HTTP server, resolve paths, allocate port, open browser |
| `stop` | stop.go | Stop all running instances via SIGTERM |
| `version` | version.go | Display version (set via ldflags) |

### internal/server/
**HTTP server and routing:**

- `server.go` - Server factory with Config (Host, Port, AllowedDirs), timeouts (R:15s, W:30s, I:60s)
- `routes.go` (463 LOC) - 5 routes:
  - `GET /` - Home page
  - `GET /view?file=X` - File rendering (markdown/code)
  - `GET /browse?dir=X` - Directory browser
  - `GET /assets/*` - Static assets (1hr cache)
  - `GET /file/*` - Direct file serving
  - Maps: mimeTypes (20+), viewableCodeExts (25+), fileIcons (30+), binaryExts (40+)
- `security.go` - `isPathSafe()` validates paths against allowedDirs, blocks ".." and null bytes

### internal/navigation/
**Plan-based sidebar and navigation:**

- `detect.go` - Detects plan.md in file directory
- `parser.go` (311 LOC) - 6 format parsers:
  - Markdown tables
  - Heading-based phases
  - Bullet lists
  - Checkboxes
  - Normalizes status to "completed"/"in-progress"/"pending"
- `sidebar.go` (176 LOC) - Accordion groups of up to 10 phases, status badges, inline sections
- `footer.go` (128 LOC) - Previous/next navigation links and context

### internal/renderer/
**Content rendering pipeline:**

- `markdown.go` (165 LOC) - Full pipeline:
  - YAML/TOML frontmatter parsing
  - Mermaid preprocessing
  - Image path resolution
  - Goldmark conversion to HTML
  - Phase heading ID generation
  - Table of contents generation
- `code.go` (104 LOC) - Line-numbered tables with Chroma syntax highlighting per line
- `toc.go` (139 LOC) - TOC generation from h1-h3, slug generation, phase IDs, table anchors
- `mermaid.go` (22 LOC) - Converts ```mermaid blocks to `<pre class="mermaid">` before Goldmark
- `images.go` (50 LOC) - Converts relative image paths to /file/ URLs

### internal/templates/
**Asset embedding and page rendering:**

- `embed.go` - Global embed.FS containing all assets
- `page.go` - PageData struct and RenderPage() template function

**PageData fields:**
- Title, Content, TOC, NavSidebar, NavFooter
- HasPlan, Frontmatter (YAML metadata)
- BackButton, HeaderNav flags

### internal/process/
**Process and port management:**

- `pid.go` (93 LOC) - PID files in `/tmp/md-novel-viewer-{port}.pid`
  - FindRunningInstances(), StopAllServers()
- `ports.go` (29 LOC) - Port range 3456-3500, IsPortAvailable(), FindAvailablePort()

## Frontend Architecture

### reader.js (743 LOC)
Vanilla JavaScript (no frameworks) implementing:

- **Theme Management:** Light/dark toggle via CSS custom properties + data-attribute
- **Font Control:** S/M/L sizes, localStorage persistence
- **Sidebar:** Toggle, drag-to-resize (200-480px)
- **Keyboard Shortcuts:**
  - T: Toggle theme
  - S: Toggle sidebar
  - ←→: Previous/next navigation
  - ?: Shortcuts cheatsheet
  - Esc: Close modals
- **Progress Bar:** Scroll tracking with requestAnimationFrame throttling
- **Mermaid v11:** Auto-render on load, re-render on theme change
- **Section Tracking:** IntersectionObserver for active sidebar state
- **Mobile:** FAB group, bottom sheet with swipe-to-close, responsive breakpoints
- **Accordion:** Plan phase groups, expand/collapse state in localStorage
- **Expandable Content:** Code blocks and diagrams scale to viewport
- **Toasts + Modal:** Notifications and shortcuts help

### HTML Template
**template.html** (149 LOC):
- FOUC prevention (inline script applies theme pre-CSS)
- Google Fonts: Inter (body), Libre Baskerville (headings), JetBrains Mono (code)
- Mermaid v11 CDN (ESM module)
- Semantic layout: header → sidebar → content → footer → overlays

### CSS Architecture
**Modular theme system** (9 files, 1,739 LOC imported by novel-theme.css):

| File | LOC | Purpose |
|------|-----|---------|
| novel-theme-variables.css | 56 | CSS custom properties, light/dark themes |
| novel-theme-base.css | 54 | Reset, scrollbars, transitions |
| novel-theme-header.css | 217 | Fixed header, font controls, progress bar |
| novel-theme-sidebar.css | 359 | Sidebar, plan nav, TOC, accordion, resize |
| novel-theme-content.css | 175 | Typography, headings, lists, nav footer |
| novel-theme-components.css | 250 | Code, blockquotes, tables, images, expand |
| novel-theme-mermaid.css | 141 | Diagram rendering, error display |
| novel-theme-overlays.css | 202 | Toast, shortcuts modal, focus trap |
| novel-theme-responsive.css | 285 | Breakpoints (900/768/600px), FAB, print |

**Theme Colors:**
- Accent: #8b4513 (saddle brown) light / #d4a574 (golden brown) dark
- Backgrounds: #faf8f3 (cream) light / #1a1a1a (black) dark
- Syntax highlighting: Chroma GitHub light/dark

### Syntax Highlighting
- `chroma-light.css` (76 LOC) - GitHub light style
- `chroma-dark.css` (73 LOC) - GitHub dark style
- Dynamically toggled by reader.js

### Directory Browser
- `directory-browser.css` (215 LOC) - Standalone styles for /browse route

## localStorage Keys

| Key | Default | Values |
|-----|---------|--------|
| `theme` | OS preference | 'light' / 'dark' |
| `novel-viewer-font` | 'M' | 'S' / 'M' / 'L' |
| `novel-viewer-sidebar` | 'visible' | 'hidden' / 'visible' |
| `novel-viewer-sidebar-width` | 280 | 200-480 (px) |
| `reader:accordion:{planId}:{phaseId}` | expanded | 'true' / 'false' |
| `reader:shortcuts-toast-shown` | false | 'true' |

## Dependencies

```
github.com/spf13/cobra v1.8.1              # CLI framework
github.com/yuin/goldmark v1.7.8            # Markdown parser
github.com/alecthomas/chroma/v2 v2.14.0    # Syntax highlighting
github.com/adrg/frontmatter v0.2.0         # YAML/TOML metadata
github.com/yuin/goldmark-highlighting/v2   # Code block highlighting
```

## Request Flow

```
mdview serve <path>
  → resolvePath() + FindAvailablePort()
  → WritePidFile()
  → Server.ListenAndServe()

GET /view?file=X
  → isPathSafe() check
  → RenderMarkdownFile() or RenderCodeFile()
    → parseFrontmatter()
    → preprocessMermaid()
    → resolveImages()
    → goldmark.Convert()
    → addPhaseHeadingIDs()
    → generateTOC()
    → DetectPlan() → ParsePlanTable()
    → GetNavigationContext()
    → GenerateNavSidebar() + GenerateNavFooter()
  → RenderPage() with template
  → HTML response

GET /browse?dir=X
  → isPathSafe()
  → renderDirectoryBrowser()
  → HTML with icons and file listing

GET /assets/*
  → Embedded FS
  → 1hr cache headers

GET /file/*
  → Direct file serve
  → MIME type detection
```

## Key Design Patterns

1. **Go embed** - All assets compiled into binary for single-file distribution
2. **Lazy template loading** - sync.Once prevents repeated file I/O
3. **PID file management** - Enables multi-instance tracking and graceful shutdown
4. **Path security** - Whitelist-based allowedDirs, prevents directory traversal
5. **6-format plan parser** - Flexible markdown structure support
6. **Mermaid preprocessing** - Handles before Goldmark to avoid code block conflicts
7. **CSS custom properties** - Theming without class manipulation
8. **localStorage persistence** - User preferences across sessions
9. **IntersectionObserver** - Efficient sidebar active state tracking
10. **requestAnimationFrame throttling** - Smooth scroll performance

## Build Configuration

**Makefile targets:**
- `build` - Compile with version ldflags
- `install` - Build + copy to ~/.local/bin/mdview
- `css` - Generate chroma CSS if missing
- `generate` - Force regenerate chroma CSS
- `test` - Run tests (`go test ./...`)
- `clean` - Remove binary + CSS
- `build-all` - Cross-compile (linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)

## File Statistics

| Component | LOC | % |
|-----------|-----|---|
| assets/styles/ | 1,739 | 33% |
| assets/ (root) | 1,272 | 24% |
| internal/navigation/ | 679 | 13% |
| internal/server/ | 538 | 10% |
| internal/renderer/ | 480 | 9% |
| cmd/ | 301 | 8% |
| internal/process/ | 122 | 2% |
| internal/templates/ | 66 | 1% |
| tools/ | 39 | <1% |

## Token Distribution (from repomix)

Top files by token count:
1. assets/reader.js (5,210 tokens) - Frontend logic
2. internal/server/routes.go (4,099 tokens) - Routing and MIME handling
3. internal/navigation/parser.go (2,919 tokens) - Plan parsing
4. assets/styles/novel-theme-sidebar.css (1,953 tokens) - Sidebar styling
5. assets/styles/novel-theme-responsive.css (1,700+ tokens) - Responsive design
