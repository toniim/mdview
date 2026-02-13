---
name: mdview
description: View markdown files with calm, book-like reading experience via HTTP server. Use for long-form content, documentation preview, novel reading, report viewing, distraction-free reading.
---

# mdview

Single-binary Go HTTP server rendering markdown files with a calm, book-like reading experience.

## Behavior

When `/mdview` is invoked:
- **No arguments**: Run `mdview serve . --host 0.0.0.0 --no-open` to browse the current working directory
- **With path argument**: Run `mdview serve <path> --host 0.0.0.0 --no-open` to view the specified file or directory
- **`stop`**: Run `mdview stop` to stop all running servers
- Always use `--host 0.0.0.0 --no-open` flags
- Display both Local and Network URLs from the output

## Quick Start

```bash
# View a markdown file
mdview serve ./plans/my-plan/plan.md --host 0.0.0.0 --no-open

# Browse a directory
mdview serve ./plans/ --host 0.0.0.0 --no-open

# Stop all running servers
mdview stop
```

## Features

### Rendering
- CommonMark markdown via Goldmark with GFM extensions
- YAML/TOML frontmatter parsing
- Syntax highlighting (25+ languages) via Chroma
- Mermaid v11 diagram rendering with theme-aware re-rendering
- Code file viewing with line numbers
- Automatic table of contents (h1-h3)
- Relative image path resolution

### Novel Theme
- Warm cream background (#faf8f3 light / #1a1a1a dark)
- Saddle brown accent (#8b4513 light / #d4a574 dark)
- Libre Baskerville serif headings, Inter body, JetBrains Mono code
- CSS custom properties for easy theming

### Plan Navigation
- Auto-detects `plan.md` in file directory
- Parses 6 plan formats (tables, headings, bullets, checkboxes)
- Accordion sidebar with status badges (completed/in-progress/pending)
- Previous/Next navigation footer and header buttons
- Phase groups (chunks of 10) with progress indicators

### Keyboard Shortcuts
- `T` - Toggle theme (light/dark)
- `S` - Toggle sidebar
- `←` / `→` - Navigate previous/next page
- `?` - Show shortcuts cheatsheet
- `Esc` - Close modal/overlay

### Mobile
- FAB (floating action buttons) for navigation
- Bottom sheet sidebar with swipe-to-close
- Responsive breakpoints: 900px (tablet), 768px (mobile), 600px (small)

### Other
- Progress bar tracking scroll position
- Expandable code blocks and diagrams (full viewport width)
- Font size control (S/M/L) persisted to localStorage
- Resizable sidebar (200-480px, drag handle)
- Print styles (content only)

## CLI Reference

```
mdview serve <path> [flags]
  -p, --port int       Port (default 3456, auto-increments to 3500)
  -H, --host string    Host (default localhost)
  -o, --open           Auto-open browser (default true)
  --no-open            Disable auto-open
  --foreground         JSON output for programmatic use

mdview stop            Stop all running instances
mdview version         Print version
```

## HTTP Routes

| Route | Description |
|-------|-------------|
| `/` | Server info page |
| `/view?file=<path>` | Render markdown/code file |
| `/browse?dir=<path>` | Directory browser |
| `/assets/*` | Embedded static assets (1hr cache) |
| `/file/*` | Direct file serving (images, binaries) |

## Building

```bash
cd ~/.claude/skills/mdview
make build      # Build binary
make install    # Install to ~/.local/bin/mdview
make test       # Run tests
```

Requires Go 1.22+.

## Security

- Path validation with allowedDirs whitelist
- Directory traversal (`..`) and null byte injection blocked
- Default localhost binding
- Read-only file serving, no code execution

## Troubleshooting

**Port in use**: Auto-increments to next available port (3456-3500)
**Images not loading**: Ensure image paths are relative to markdown file
**Server won't stop**: Check `/tmp/md-novel-viewer-*.pid` for stale PID files
**Remote access**: Use `--host 0.0.0.0` to bind to all interfaces
