# mdview - Project Overview & PDR

## Executive Summary

**mdview** is a lightweight, single-binary HTTP server that transforms markdown files and code into a calm, book-like reading experience. Designed for developers and documentation writers, it provides distraction-free viewing with intuitive navigation, theme customization, and plan-based documentation support.

- **Status:** Version 0.1.0 (initial release)
- **Language:** Go 1.22
- **Distribution:** Single embedded binary
- **Use Cases:** Documentation, long-form content, plan-based project tracking, code viewing

## Vision

Create a minimal, elegant markdown viewer that prioritizes reading experience over feature complexity. mdview should feel like reading a physical book—calm, focused, and responsive to user preferences.

## Core Values

1. **Simplicity** - Single binary, no external dependencies, minimal configuration
2. **Distraction-Free** - Clean UI, smooth scrolling, no ads or tracking
3. **Accessibility** - Keyboard shortcuts, theme support, responsive design
4. **Extensibility** - Modular architecture for adding new renderers and plan formats

## Product Requirements

### Functional Requirements

#### FR-1: File Serving
- **Serve markdown files** with full HTML rendering including:
  - YAML/TOML frontmatter parsing
  - Syntax highlighting in code blocks
  - Mermaid diagram rendering
  - Table of contents generation
  - Relative image path resolution
- **Serve code files** with:
  - Line numbering
  - Syntax highlighting (Chroma)
  - MIME type detection
- **Serve static assets** from embedded filesystem with caching

**Acceptance Criteria:**
- GET /view?file=X renders markdown with all features
- GET /view?file=X (code file) renders with syntax highlighting
- GET /assets/* returns with proper cache headers
- All paths validated against allowedDirs whitelist

#### FR-2: Navigation System
- **Detect plan.md** in file directories and parse structure
- **Support 6 plan formats:**
  - Markdown tables
  - Heading-based phases
  - Bullet lists
  - Checkboxes
  - Mixed structures
- **Generate sidebar navigation** with:
  - Status badges (completed/in-progress/pending)
  - Accordion groups of phases
  - Clickable links to sections
  - Collapse/expand state persistence
- **Generate footer navigation** with:
  - Previous/next page links
  - Breadcrumb trails
  - Current phase highlighting

**Acceptance Criteria:**
- Plan detection works for all 6 formats
- Sidebar renders with proper styling
- Accordion states persist across sessions
- Navigation links correctly update on page change

#### FR-3: User Interface
- **Theme system:**
  - Light/dark mode toggle
  - CSS custom property-based theming
  - OS preference detection
  - Persistence across sessions
- **Font size control:** S/M/L sizes with localStorage persistence
- **Sidebar management:**
  - Toggle visibility
  - Drag-to-resize (200-480px range)
  - Smooth transitions
- **Keyboard shortcuts:**
  - T: Toggle theme
  - S: Toggle sidebar
  - ←→: Previous/next navigation
  - ?: Show shortcuts
  - Esc: Close modals
- **Mobile responsiveness:**
  - Responsive breakpoints (900/768/600px)
  - FAB group for mobile controls
  - Bottom sheet for navigation
  - Touch-friendly interactions

**Acceptance Criteria:**
- Theme toggles smoothly without page reload
- Font sizes apply consistently
- Sidebar resizes smoothly within bounds
- All keyboard shortcuts function correctly
- Mobile layouts render properly at all breakpoints

#### FR-4: Reading Features
- **Progress bar** showing scroll position
- **Table of contents** with links to sections
- **Code block expand** for full viewport viewing
- **Diagram expand** for Mermaid diagrams
- **Toast notifications** for user feedback
- **Shortcuts modal** accessible via ? key

**Acceptance Criteria:**
- Progress bar updates smoothly
- TOC links jump to correct sections
- Expandable elements fill viewport
- Toasts appear and auto-dismiss
- Shortcuts modal displays keyboard map

#### FR-5: Server Management
- **Start server:** `mdview serve <path>` with flags:
  - `-p/--port` for custom port
  - `-H/--host` for binding address
  - `-o/--open` to open browser
  - `--foreground` for JSON output (Claude Code)
- **Stop servers:** `mdview stop` terminates all instances
- **Multi-instance support** via PID file tracking
- **Graceful shutdown** on SIGINT/SIGTERM
- **Automatic port allocation** (3456-3500 range)

**Acceptance Criteria:**
- Server starts and listens on specified port
- Browser opens automatically with `-o` flag
- `mdview stop` terminates server gracefully
- PID files correctly track instances
- Port auto-allocation works for 45+ concurrent instances

#### FR-6: Security
- **Path validation:**
  - Whitelist-based directory access
  - ".." traversal prevention
  - Null byte detection
- **Access control:**
  - Default localhost binding
  - Optional network exposure via flag
- **Asset embedding:**
  - No external file dependencies
  - Compiled assets cannot be modified

**Acceptance Criteria:**
- Attempts to access files outside allowedDirs fail
- Directory traversal attacks fail
- All paths properly validated before access
- No security warnings in code review

### Non-Functional Requirements

#### NFR-1: Performance
- **Page load:** < 500ms for markdown rendering
- **Scroll performance:** 60 FPS with smooth animations
- **Memory footprint:** < 50MB for typical document
- **Startup time:** < 1s from CLI to listening port
- **HTTP timeouts:**
  - Read: 15s
  - Write: 30s
  - Idle: 60s

**Acceptance Criteria:**
- Rendering benchmarks meet targets
- No jank on scroll (rAF throttling)
- Memory profiling shows clean heap
- Server responds to SIGTERM immediately

#### NFR-2: Reliability
- **Single point of failure:** Zero (stateless HTTP server)
- **Data loss:** Impossible (read-only file serving)
- **Error handling:** Graceful degradation on errors
- **Availability:** 99.9% (single server uptime)

**Acceptance Criteria:**
- Server restart doesn't lose state
- Render errors show user-friendly messages
- No panic() in production paths
- SIGTERM forces clean shutdown

#### NFR-3: Compatibility
- **Browsers:** Chrome, Firefox, Safari (last 2 versions)
- **OS:** Linux, macOS (Intel/ARM), Windows
- **Go:** 1.22+
- **Markdown:** CommonMark + extensions (tables, strikethrough, tasklists)

**Acceptance Criteria:**
- Cross-browser testing passes
- Binary compiles for all platforms
- Syntax highlighting works on all platforms
- No platform-specific bugs

#### NFR-4: Maintainability
- **Code quality:**
  - Max 200 LOC per file
  - Single responsibility per package
  - Clear error messages
  - Documentation comments
- **Test coverage:** > 70% unit tests
- **Architecture:** Layered (CLI → Server → Renderer)

**Acceptance Criteria:**
- Code review passes without style issues
- Tests are comprehensive and readable
- Architecture allows feature additions without refactoring

#### NFR-5: Accessibility
- **WCAG AA compliance:**
  - Color contrast > 4.5:1
  - Keyboard navigation throughout
  - Semantic HTML
  - ARIA labels where needed
- **Font:** Readable serif/sans-serif with good line-height

**Acceptance Criteria:**
- Accessibility audit passes
- All features work with keyboard only
- Screen reader testing succeeds

## Scope

### In Scope
- Single-file markdown serving
- Directory browsing
- Code file rendering with syntax highlighting
- Plan-based navigation
- Theme customization
- Responsive UI
- Mermaid diagram rendering
- YAML/TOML frontmatter support
- Table of contents generation

### Out of Scope
- User authentication/authorization
- Database integration
- Multi-user collaboration
- Version control integration
- Search functionality (planned for v1.0)
- PDF export (planned for v1.0)
- Plugin system (planned for v1.1+)

## Architecture Decisions

### Design Pattern: Single Binary Distribution
- **Rationale:** Simplicity, portability, easy deployment
- **Trade-offs:** All assets compiled in, larger binary size
- **Alternative:** Separate asset files (rejected)

### Framework Choice: Vanilla JavaScript + CSS
- **Rationale:** No dependencies, lightweight, full control
- **Trade-offs:** More manual state management
- **Alternative:** React/Vue (rejected due to complexity)

### Markdown Parser: Goldmark
- **Rationale:** Fast, extensible, CommonMark compliant
- **Trade-offs:** Pure Go implementation (no C bindings)
- **Alternative:** Pandoc (rejected due to external binary)

### Routing: Cobra CLI + Go net/http
- **Rationale:** Idiomatic Go, minimal dependencies
- **Trade-offs:** Manual request handling, no middleware framework
- **Alternative:** Echo/Gin (rejected for simplicity)

### Theming: CSS Custom Properties
- **Rationale:** No runtime JS overhead, works without JS
- **Trade-offs:** Limited to CSS capabilities
- **Alternative:** CSS-in-JS (rejected for size/complexity)

## Success Metrics

### Adoption
- GitHub stars: > 100 by v1.0
- Monthly downloads: > 500 by v1.0
- Community PRs: > 5 by v1.0

### Quality
- Test coverage: > 70%
- Zero critical security issues
- < 5 open bugs at release

### Performance
- Page load: < 500ms (median)
- Scroll FPS: > 55 (95th percentile)
- Memory: < 50MB typical usage

### User Satisfaction
- Issue resolution time: < 1 week
- Feature request response: < 3 days
- Community feedback: Positive

## Technical Stack

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Language | Go 1.22 | Performance, single binary, cross-platform |
| CLI | Cobra | Idiomatic Go, mature ecosystem |
| Markdown | Goldmark | Fast, extensible, CommonMark-compliant |
| Highlighting | Chroma v2 | Comprehensive language support |
| Frontend | Vanilla JS | Lightweight, no dependencies |
| Styling | Modular CSS | No dependencies, full control |
| Build | Make | Standard, portable, simple |
| Embed | Go embed | Built-in, no external tools |

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/spf13/cobra | v1.8.1 | CLI framework |
| github.com/yuin/goldmark | v1.7.8 | Markdown parsing |
| github.com/alecthomas/chroma/v2 | v2.14.0 | Syntax highlighting |
| github.com/adrg/frontmatter | v0.2.0 | YAML/TOML metadata |
| github.com/yuin/goldmark-highlighting/v2 | latest | Code highlighting |

**Zero external JS dependencies** - all frontend code vanilla JavaScript.

## Release Plan

### Version 0.1.0 (Current)
- Core markdown rendering
- Basic navigation
- Theme system
- Code highlighting
- Plan detection (6 formats)

### Version 0.2.0
- Bug fixes from user feedback
- Performance optimizations
- Improved mobile UX
- Extended syntax support

### Version 1.0.0
- Stable API
- Search functionality
- PDF export
- Improved accessibility
- Extended documentation

### Version 1.1.0+
- Plugin system
- Custom renderers
- Advanced plan formats
- Community contributions

## Risk Assessment

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|-----------|
| Security vulnerability | High | Low | Regular audits, dependabot |
| Breaking API change | Medium | Low | Semantic versioning |
| Performance regression | Medium | Low | Benchmarking in CI |
| Go version incompatibility | Low | Low | Vendor testing with new versions |
| Browser compatibility | Medium | Medium | Automated cross-browser testing |

## Next Steps

1. **Immediate:** Release v0.1.0, gather community feedback
2. **Short-term:** Fix reported bugs, optimize performance
3. **Medium-term:** Add search, improve navigation
4. **Long-term:** Build plugin ecosystem, expand feature set

## Related Documentation

- [Codebase Summary](./codebase-summary.md)
- [Code Standards](./code-standards.md)
- [System Architecture](./system-architecture.md)
- [Design Guidelines](./design-guidelines.md)
- [Project Roadmap](./project-roadmap.md)
