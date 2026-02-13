# Code Standards & Patterns

## Go Code Standards

### File Organization

**Go files** follow standard naming conventions:
- `main.go` - Entry point with asset embedding
- `{feature}.go` - Single responsibility per file
- `{feature}_test.go` - Test files alongside implementation

**Package structure:**
- `cmd/` - CLI commands (Cobra-based)
- `internal/` - Unexported packages (server, navigation, renderer, templates, process)
- `tools/` - Utility scripts (CSS generators)

### Code Style

- **Interfaces:** Use composition over inheritance
- **Error Handling:** Return errors as last argument, use `if err != nil` checks
- **Naming:** Follow Go conventions (CamelCase, avoid underscores)
- **Comments:** Exported functions have doc comments (`// FunctionName does X`)
- **Constants:** Define in logical groups at package level
- **Functions:** Keep under 50 lines; extract complex logic to helpers

### Security Practices

**Path Validation:**
```go
// Always validate user input paths before file operations
func isPathSafe(requestedPath string, allowedDirs []string) bool {
    // Check against allowed directories
    // Block ".." and null bytes
    // Return false for unsafe paths
}
```

**Access Control:**
- HTTP server runs on localhost by default (`-H` flag for custom)
- PID files in `/tmp` with restrictive permissions
- No arbitrary code execution or shell commands

### HTTP Server Patterns

**Server Configuration:**
```go
type Config struct {
    Host       string
    Port       int
    AllowedDirs []string
}
```

**Timeouts:**
- Read timeout: 15 seconds
- Write timeout: 30 seconds
- Idle timeout: 60 seconds

**MIME Type Mapping:**
- Maintain `mimeTypes` map for 20+ file types
- Viewable code extensions list (25+ types)
- Binary extensions list to prevent rendering

### Markdown Rendering Pipeline

**Standard flow:**
1. Parse YAML/TOML frontmatter with `adrg/frontmatter`
2. Preprocess Mermaid blocks (convert to HTML before Goldmark)
3. Resolve relative image paths to `/file/` URLs
4. Convert with `yuin/goldmark` + highlighting extension
5. Add phase heading IDs for navigation
6. Generate table of contents from h1-h3
7. Detect and parse plan.md if present
8. Generate sidebar and footer navigation
9. Render with Go template

**Extension Points:**
- Custom renderers for code blocks, images, links
- Table of contents generation with slug support
- Frontmatter extraction for metadata

### Plan Parser Implementation

**Support 6 formats:**
1. Markdown tables with | Phase | Status | column
2. Heading-based: `## Phase Name (status)`
3. Bullet lists: `- Phase Name (status)`
4. Checkboxes: `- [x] Phase Name`
5. Mixed structures (prioritize consistent format)
6. Status normalization: "completed" / "in-progress" / "pending"

**Example:**
```markdown
| Phase | Status |
|-------|--------|
| Planning | completed |
| Development | in-progress |
| Testing | pending |
```

### Dependency Management

**Required packages:**
- `github.com/spf13/cobra` - CLI framework
- `github.com/yuin/goldmark` - Markdown parsing
- `github.com/alecthomas/chroma/v2` - Syntax highlighting
- `github.com/adrg/frontmatter` - YAML/TOML metadata
- `github.com/yuin/goldmark-highlighting/v2` - Code block highlighting

**No external UI frameworks** - vanilla Go templates and JavaScript

## JavaScript/Frontend Standards

### Code Organization

**reader.js (743 LOC)** - Main frontend logic:
- Module-style functions (no frameworks)
- Event delegation for performance
- localStorage for persistence
- IntersectionObserver for efficiency
- requestAnimationFrame for scroll performance

### UI State Management

**localStorage Keys:**
- All preference keys prefixed with `novel-viewer-` or `reader:`
- Values stored as strings; parse as needed
- Defaults applied on first run

**State Preservation:**
- Theme: reads OS preference, allows override
- Font size: persists across sessions
- Sidebar width: restored on page load
- Accordion states: per-phase collapsible groups
- Toast visibility: once per session

### Theme System

**CSS Custom Properties:**
- Light theme variables in `:root`
- Dark theme variables in `[data-theme="dark"]`
- No class-based styling; use data attributes
- Transitions smooth between theme switches

**Design tokens:**
- Colors: accent, text, background, borders
- Spacing: 0.5rem, 1rem, 1.5rem, 2rem increments
- Typography: Inter (body), Libre Baskerville (headings), JetBrains Mono (code)
- Breakpoints: 900px, 768px, 600px

### Keyboard Shortcuts

**Standard bindings:**
- T: Toggle theme
- S: Toggle sidebar
- ←: Previous page
- →: Next page
- ?: Show shortcuts
- Esc: Close modal

**No conflicts** with browser/OS shortcuts

### Performance Guidelines

1. **Scroll Handlers:** Use requestAnimationFrame throttling
2. **DOM Queries:** Cache frequently accessed elements
3. **Event Listeners:** Use event delegation where possible
4. **IntersectionObserver:** Track active sections without constant polling
5. **Mermaid Rendering:** Re-render only on theme change, not every scroll

### Accessibility

- Semantic HTML (nav, main, footer)
- ARIA labels for interactive elements
- Focus visible on all interactive controls
- Color contrast meets WCAG AA
- Keyboard-navigable throughout

### Mermaid Integration

- CDN load from esm.sh (v11 ESM module)
- Auto-initialize with default config
- Error handling for invalid diagrams (show error message)
- Re-render on theme change via `mermaid.contentLoaded()`
- Expandable diagrams with viewport-width toggle

## CSS Architecture

### File Organization

**9 modular files imported by novel-theme.css:**

| File | Purpose |
|------|---------|
| novel-theme-variables.css | Design tokens (colors, fonts, spacing) |
| novel-theme-base.css | Reset, scrollbars, baseline styles |
| novel-theme-header.css | Header, controls, progress bar |
| novel-theme-sidebar.css | Sidebar, navigation, accordion |
| novel-theme-content.css | Body text, headings, lists, links |
| novel-theme-components.css | Code, quotes, tables, images |
| novel-theme-mermaid.css | Diagram styling |
| novel-theme-overlays.css | Modals, toasts, overlays |
| novel-theme-responsive.css | Breakpoints, mobile layouts |

### Styling Patterns

**No utility classes** - all styles use semantic selectors
**Modular approach:**
- Each file handles one responsibility
- Minimal selectors to prevent cascading conflicts
- CSS custom properties for theming

**Responsive design:**
- Mobile-first approach
- Three breakpoints: 900px, 768px, 600px
- Print styles included

**Animation:**
- Smooth transitions (200-300ms)
- No animations on reduced-motion preference
- GPU-accelerated transforms where possible

### Color Palette

**Light theme:**
- Text: #2c2c2c (near-black)
- Accent: #8b4513 (saddle brown)
- Background: #faf8f3 (warm cream)
- Border: #e8e6df (light gray)

**Dark theme:**
- Text: #f0ede8 (off-white)
- Accent: #d4a574 (golden brown)
- Background: #1a1a1a (black)
- Border: #3a3a3a (dark gray)

**Syntax highlighting:** Chroma GitHub light/dark themes

### Typography

- **Body:** Inter, 16px (mobile 14px), line-height 1.6
- **Headings:** Libre Baskerville, scaled (h1: 2.2em, h2: 1.8em, h3: 1.5em)
- **Code:** JetBrains Mono, 14px, line-height 1.5
- **Print:** Serif font, black text, larger line-height

## HTML Template Standards

**template.html patterns:**
- Semantic HTML5 elements (header, nav, main, footer, aside)
- FOUC prevention with inline style in head
- Single point of templating (Go templates)
- Minimal inline scripts (only theme application)

**Template variables (PageData):**
- `.Title` - Page title
- `.Content` - Rendered HTML content
- `.TOC` - Table of contents items
- `.NavSidebar` - Plan navigation HTML
- `.NavFooter` - Previous/next navigation
- `.HasPlan` - Boolean flag
- `.Frontmatter` - YAML metadata map
- `.BackButton`, `.HeaderNav` - Navigation flags

## Build & Deployment

### Asset Embedding

**Go embed directive:**
```go
//go:embed all:assets
var Assets embed.FS
```

**All static files compiled into binary:**
- reader.js
- template.html
- CSS files (9 files)
- Syntax highlighting CSS

**Benefits:**
- Single binary distribution
- No external file dependencies
- Faster load times
- No PATH issues

### Cross-compilation

**Supported platforms:**
- linux/arm64 (Raspberry Pi, ARM servers)
- darwin/amd64 (Intel Mac)
- darwin/arm64 (Apple Silicon Mac)
- windows/amd64 (Windows)

**Build with `make build-all`**

### Version Management

- Version injected via ldflags at compile time
- Displayed with `mdview version`
- Semantic versioning (MAJOR.MINOR.PATCH)

## Testing Guidelines

**Unit tests:**
- Test utility functions (path validation, plan parsing, rendering)
- Mock HTTP handlers
- Test error cases

**Integration tests:**
- Full request/response cycles
- File serving and rendering
- Navigation generation

**No external services** - all tests self-contained

## Documentation Standards

**Code comments:**
- Exported functions have doc comments
- Complex logic has inline explanations
- No over-commenting (code should be self-documenting)

**README sections:**
- Project overview
- Quick start
- Features
- Build instructions
- Configuration

**Architecture documentation:**
- Request flow diagrams
- Package responsibilities
- Data flow between components
