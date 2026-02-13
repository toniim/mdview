# System Architecture

## High-Level Overview

mdview is a **single-binary HTTP server** that renders markdown and code files with a calm, book-like reading experience. The system comprises three main layers:

1. **CLI Layer** - Command-line interface for starting/stopping the server
2. **Go Backend** - HTTP server, file rendering, plan parsing, and navigation
3. **Web Frontend** - Vanilla JavaScript for theme management, interactions, and Mermaid rendering

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Browser                              │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │         Web Frontend (reader.js + CSS)                   │  │
│  │  Theme Management │ Sidebar Toggle │ Mermaid Rendering  │  │
│  └──────────────────────────────────────────────────────────┘  │
└──────────────────────────┬──────────────────────────────────────┘
                           │ HTTP
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                    mdview HTTP Server                            │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ Routes (5 endpoints) + Security Layer                      │ │
│  │ GET / | /view | /browse | /assets/* | /file/*             │ │
│  └────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ Rendering Pipeline                                         │ │
│  │ Frontmatter → Mermaid Preprocess → Image Resolution      │ │
│  │ Goldmark Convert → TOC Generation → Plan Detection        │ │
│  │ Navigation Generation → Template Rendering                │ │
│  └────────────────────────────────────────────────────────────┘ │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ Navigation Engine                                          │ │
│  │ Plan Detection → 6-Format Parser → Sidebar + Footer Nav   │ │
│  └────────────────────────────────────────────────────────────┘ │
└──────────────────────┬───────────────────────────────────────────┘
                       │ File System Access
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│              Allowed Directories (Whitelist)                     │
│ markdown files, code files, images, plan.md                     │
└─────────────────────────────────────────────────────────────────┘
```

## Request Flow

### Markdown View Request

```
GET /view?file=/path/to/file.md
  ↓
Security Check (isPathSafe)
  ├─ In allowed directories?
  ├─ No ".." traversal?
  └─ No null bytes?
  ↓
Detect File Type
  ├─ Markdown (.md)?
  └─ Code file?
  ↓
Render Markdown
  ├─ Parse YAML/TOML frontmatter
  ├─ Preprocess Mermaid blocks
  │  (```mermaid → <pre class="mermaid">)
  ├─ Resolve relative image paths
  │  (../images/foo.png → /file/?path=...)
  ├─ Convert with Goldmark
  │  (with syntax highlighting extension)
  ├─ Add Phase Heading IDs
  │  (h2 "Phase X" → id="phase-x")
  └─ Generate Table of Contents
      (extract h1-h3, create slug links)
  ↓
Detect Plan
  ├─ Look for plan.md in file directory
  └─ If found, parse navigation structure
  ↓
Generate Navigation
  ├─ Parse plan.md (6 formats supported)
  ├─ Build sidebar with status badges
  └─ Create previous/next links
  ↓
Render Page Template
  ├─ Pass Title, Content, TOC, Nav to template
  ├─ Template adds reader.js + CSS
  └─ Return HTML
  ↓
HTTP Response (200 OK, text/html)
```

### Code File Rendering

```
GET /view?file=/path/to/file.go
  ↓
Security & MIME Type Check
  ├─ Safe path?
  └─ Viewable code extension?
  ↓
Render Code File
  ├─ Read file contents
  ├─ Line-by-line syntax highlighting (Chroma)
  ├─ Build line-numbered table
  └─ Generate HTML
  ↓
Generate Basic Navigation
  ├─ No plan parsing for code files
  └─ Simple previous/next links (if applicable)
  ↓
Render Page Template
  └─ Return HTML
```

### Directory Browse Request

```
GET /browse?dir=/path/to/directory
  ↓
Security Check
  ├─ Safe path?
  └─ Is directory?
  ↓
List Directory Contents
  ├─ Get file/folder list
  ├─ Map file types to icons
  └─ Generate HTML listing
  ↓
Render Directory Template
  ├─ Navigation breadcrumbs
  ├─ Grouped file listing
  └─ Links to view/browse each item
```

### Static Asset Request

```
GET /assets/reader.js
GET /assets/styles/novel-theme.css
  ↓
Check Embedded FS
  └─ Found in go:embed
  ↓
HTTP Response (200 OK)
  ├─ Cache-Control: max-age=3600
  └─ Content (JS/CSS)
```

## Component Architecture

### Backend Components

#### cmd/ - CLI Layer
- **root.go** - Cobra command setup
- **serve.go** - Server bootstrap
  - Path resolution (handles `.`, relative paths, absolute paths)
  - Port allocation (3456-3500 range)
  - PID file management
  - Auto-browser open (with `-o` flag)
  - Signal handling (SIGINT/SIGTERM for graceful shutdown)
  - Foreground mode for Claude Code integration (`--foreground` JSON output)
- **stop.go** - Stop all instances
- **version.go** - Version display

#### internal/server/ - HTTP Server
- **server.go** - Server factory
  ```go
  type Config struct {
    Host       string
    Port       int
    AllowedDirs []string
  }
  ```
  - Creates net.http.Server with timeouts
  - Sets up graceful shutdown

- **routes.go** (463 LOC) - Main routing logic
  - Route handlers for 5 endpoints
  - MIME type detection (20+ types)
  - Viewable code extensions list (25+ types)
  - File icon mapping (30+ types)
  - Binary file detection (40+ types)
  - Full rendering pipeline orchestration

- **security.go** - Path validation
  - `isPathSafe()` function
  - Whitelist-based directory validation
  - ".." traversal blocking
  - Null byte detection

#### internal/navigation/ - Plan Parsing & Navigation (679 LOC)
- **detect.go**
  - Scans directory for plan.md
  - Returns plan contents or nil

- **parser.go** (311 LOC) - 6-format plan parser
  ```
  Format 1: Markdown table
  | Phase | Status |
  |-------|--------|
  | Planning | completed |

  Format 2: Heading-based
  ## Phase Name (completed)

  Format 3: Bullet list
  - Phase Name (in-progress)

  Format 4: Checkboxes
  - [x] Phase Name (auto-completed)

  Format 5-6: Mixed structures
  ```
  - `ParsePlanTable()` returns `[]Phase` with name and status
  - `normalizeStatus()` converts to standard format

- **sidebar.go** (176 LOC)
  - `GenerateNavSidebar()` creates accordion HTML
  - Groups phases (max 10 per group)
  - Adds status badges (✓ completed, ⧗ in-progress, ○ pending)
  - Inline sub-sections with links
  - Stores accordion state in localStorage

- **footer.go** (128 LOC)
  - `GetNavigationContext()` finds current phase
  - `GenerateNavFooter()` creates prev/next links
  - Breadcrumb navigation

#### internal/renderer/ - Content Rendering (480 LOC)
- **markdown.go** (165 LOC) - Full pipeline
  ```go
  type RenderResult struct {
    HTML       string
    TOC        []TOCItem
    Frontmatter map[string]interface{}
    Title      string
  }
  ```
  - `RenderMarkdownFile()` orchestrates:
    1. Parse frontmatter (YAML/TOML)
    2. Preprocess Mermaid
    3. Resolve images
    4. Convert with Goldmark
    5. Generate TOC
    6. Add phase IDs

- **code.go** (104 LOC)
  - `RenderCodeFile()` for syntax highlighting
  - Line-numbered table rendering
  - Per-line Chroma highlighting

- **toc.go** (139 LOC)
  - `generateTOC()` from h1-h3 headings
  - `Slugify()` for URL-safe IDs
  - Phase heading ID generation
  - Table anchor links

- **mermaid.go** (22 LOC)
  - `preprocessMermaid()` before Goldmark
  - Converts ```mermaid to `<pre class="mermaid">` tags
  - Prevents code block conflicts

- **images.go** (50 LOC)
  - `resolveImages()` for relative paths
  - Converts to `/file/?path=X` URLs
  - Maintains proper escaping

#### internal/templates/ - Asset Embedding (66 LOC)
- **embed.go**
  ```go
  //go:embed all:assets
  var Assets embed.FS
  ```
  - Global Assets FS
  - `SetAssets()` for testing/override

- **page.go**
  ```go
  type PageData struct {
    Title      string
    Content    string
    TOC        []TOCItem
    NavSidebar string
    NavFooter  string
    HasPlan    bool
    Frontmatter map[string]interface{}
    BackButton bool
    HeaderNav  string
  }
  ```
  - `RenderPage()` renders template
  - Lazy-loaded template with sync.Once

#### internal/process/ - Process Management (122 LOC)
- **pid.go** (93 LOC)
  - PID files: `/tmp/md-novel-viewer-{port}.pid`
  - `WritePidFile()` creates file
  - `FindRunningInstances()` scans /tmp
  - `StopAllServers()` sends SIGTERM

- **ports.go** (29 LOC)
  - Port range: 3456-3500
  - `IsPortAvailable()` checks TCP binding
  - `FindAvailablePort()` auto-allocates

### Frontend Components

#### reader.js (743 LOC)
**State Management:**
- localStorage for theme, font size, sidebar width, accordion states
- CSS data-attribute for theme switching

**Core Features:**

| Feature | Implementation |
|---------|---|
| Theme Toggle | CSS custom properties + data-attribute |
| Font Control | localStorage persistence, 3 sizes (S/M/L) |
| Sidebar | Toggle visibility, drag-to-resize (200-480px) |
| Keyboard Shortcuts | T/S/←→/?/Esc |
| Progress Bar | Scroll tracking with rAF throttling |
| Mermaid Rendering | v11 CDN, re-render on theme change |
| Section Tracking | IntersectionObserver for active state |
| Accordion | Per-phase collapse/expand state |
| Mobile | FAB group, bottom sheet, responsive |

**Performance Optimizations:**
- requestAnimationFrame throttling for scroll
- Event delegation for click handlers
- Cached DOM element references
- IntersectionObserver for lazy tracking

#### CSS System (1,739 LOC across 9 files)

**Architecture:**

```
novel-theme.css (imports all below)
├── novel-theme-variables.css (56 LOC)
│   └─ CSS custom properties, light/dark themes
├── novel-theme-base.css (54 LOC)
│   └─ Reset, scrollbars, transitions
├── novel-theme-header.css (217 LOC)
│   └─ Header, font controls, progress bar
├── novel-theme-sidebar.css (359 LOC)
│   └─ Sidebar, TOC, accordion, resize
├── novel-theme-content.css (175 LOC)
│   └─ Typography, headings, lists
├── novel-theme-components.css (250 LOC)
│   └─ Code, quotes, tables, images
├── novel-theme-mermaid.css (141 LOC)
│   └─ Diagram styling
├── novel-theme-overlays.css (202 LOC)
│   └─ Modals, toasts, focus trap
└── novel-theme-responsive.css (285 LOC)
    └─ Breakpoints, mobile, print
```

**Theming Strategy:**
- CSS custom properties in `:root` (light)
- Override in `[data-theme="dark"]`
- No class-based themes
- Smooth transitions (200-300ms)

**Breakpoints:**
- 900px: Desktop → tablet
- 768px: Tablet → mobile
- 600px: Small mobile
- Print: Special print stylesheet

#### template.html (149 LOC)
```html
<!DOCTYPE html>
<html>
  <head>
    {{.MetaTags}}
    <link rel="stylesheet" href="/assets/styles/novel-theme.css">
    <script>
      /* FOUC prevention: apply theme before CSS loads */
      applyTheme();
    </script>
  </head>
  <body>
    <header>{{.HeaderNav}}</header>
    <aside id="sidebar">{{.NavSidebar}}</aside>
    <main>{{.Content}}</main>
    <footer>{{.NavFooter}}</footer>
    <div id="overlays"></div>
    <script src="/assets/reader.js"></script>
  </body>
</html>
```

## Data Flow

### Rendering Pipeline
```
File Request
  ↓
Security Validation
  ↓
Frontmatter Parsing ──→ Page Metadata
  ↓
Content Processing
  ├─ Mermaid Preprocessing
  ├─ Image Resolution
  ├─ Markdown→HTML Conversion
  ├─ Syntax Highlighting
  └─ Phase ID Injection
  ↓
Navigation Processing
  ├─ Plan Detection
  ├─ 6-Format Parsing
  ├─ Sidebar Generation
  └─ Footer Navigation
  ↓
TOC Generation
  ├─ Extract h1-h3
  ├─ Create Slugs
  └─ Generate Links
  ↓
Template Rendering
  ├─ Insert Content
  ├─ Insert Navigation
  ├─ Inject Reader JS
  └─ Inject CSS
  ↓
HTTP Response
```

### Theme Application
```
Browser Load
  ↓
FOUC Prevention Script (inline)
  ├─ Read localStorage[theme]
  └─ Apply data-attribute before CSS
  ↓
CSS Loads
  ├─ Apply custom properties from :root or [data-theme]
  └─ Render with correct colors
  ↓
Reader JS Loads
  ├─ Listen for theme toggle
  ├─ Update localStorage
  └─ Re-render Mermaid diagrams
```

## Security Model

### Path Validation
- **Whitelist-based:** Only serve from allowedDirs
- **Traversal prevention:** Block ".." in paths
- **Null byte protection:** Reject paths with \x00
- **Absolute path conversion:** Resolve symlinks safely

### Access Control
- **Default localhost:** Only accessible from machine (`127.0.0.1:PORT`)
- **Custom host flag:** `-H 0.0.0.0` for network access (use with caution)
- **No authentication:** Assumes trusted network

### Binary Protection
- **Single binary:** All assets compiled in, no external file dependencies
- **No code execution:** No script evaluation, template safe

## Performance Characteristics

### Memory
- Lazy template loading (sync.Once)
- Streaming HTTP responses
- No full file buffering for rendering

### Disk
- No temp files
- All rendering in-memory
- Chroma CSS cached in /tmp

### Network
- 15s read timeout
- 30s write timeout
- 1hr asset cache headers
- Streaming large files

## Deployment Model

### Single Binary Distribution
```bash
mdview serve /path/to/docs --port 3456 --open
  ↓
Binary embeds:
  ├─ reader.js
  ├─ template.html
  ├─ novel-theme.css (+ 8 imports)
  ├─ chroma-light.css, chroma-dark.css
  └─ directory-browser.css
  ↓
Zero external dependencies
```

### Multi-Instance Support
- PID files track running servers
- Port auto-allocation (3456-3500)
- `mdview stop` terminates all instances
- Graceful shutdown on SIGINT/SIGTERM

## Extension Points

### Adding New Renderers
1. Implement renderer in `internal/renderer/`
2. Add file type check in `routes.go`
3. Register MIME type and viewable extension

### Custom Plan Formats
1. Add parser case in `parser.go`
2. Update `normalizeStatus()` if needed
3. Regenerate sidebar in `sidebar.go`

### Theme Customization
1. Modify CSS custom properties in `novel-theme-variables.css`
2. Update dark theme overrides
3. All other files inherit changes

### Frontend Enhancements
1. Add features to `reader.js`
2. Add localStorage keys with `novel-viewer-` prefix
3. Add corresponding CSS to overlays/components
