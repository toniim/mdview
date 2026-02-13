# mdview Project Roadmap

## Current Status: v0.1.0 (Alpha)

Initial release with core markdown viewing capabilities, plan-based navigation, and theme system.

## Version Timeline

### v0.1.0 - Alpha Release (Current)

**Status:** Released

**Features:**
- Core markdown rendering with Goldmark
- YAML/TOML frontmatter support
- Syntax highlighting with Chroma
- Mermaid diagram rendering
- Table of contents generation
- Plan detection (6 format parser)
- Sidebar navigation with accordion
- Theme system (light/dark)
- Font size control
- Keyboard shortcuts
- Mobile responsive design
- Code file rendering
- Directory browsing
- Single-binary distribution

**Known Limitations:**
- No search functionality
- No PDF export
- Limited mobile optimization
- Plan parser may not handle edge cases
- Mermaid v11 CDN dependency

**Testing:**
- Basic unit tests for core functions
- Manual integration testing
- Cross-browser testing (Chrome, Firefox, Safari)
- Platform testing (Linux, macOS, Windows)

---

### v0.2.0 - Bug Fixes & Performance (Q1 2026)

**Timeline:** 4-6 weeks after v0.1.0 release

**Goals:**
- Fix reported bugs
- Improve performance
- Enhanced mobile UX
- Better error messages

**Features:**

1. **Bug Fixes**
   - Address user-reported issues
   - Fix edge cases in plan parser
   - Correct CSS rendering issues
   - Fix keyboard shortcut conflicts

2. **Performance Optimizations**
   - Reduce bundle size
   - Optimize CSS loading
   - Improve scroll performance
   - Cache improvements

3. **Mobile Enhancements**
   - Improved touch interactions
   - Better FAB positioning
   - Bottom sheet refinements
   - Viewport meta tag optimization

4. **Error Handling**
   - User-friendly error messages
   - Better logging
   - Recovery from render errors
   - Graceful Mermaid failures

**Quality Targets:**
- Test coverage: > 75%
- Mobile performance: 60 FPS maintained
- Load time: < 300ms median

**Acceptance Criteria:**
- Zero critical bugs from v0.1.0 feedback
- Performance benchmarks met
- All mobile breakpoints tested
- Community feedback incorporated

---

### v1.0.0 - Stable Release (Q2 2026)

**Timeline:** 10-12 weeks after v0.2.0

**Goals:**
- Stable API
- Search functionality
- PDF export
- Advanced features
- Full documentation

**Features:**

1. **Full-Text Search**
   - Index all markdown content
   - Real-time search as-you-type
   - Search result highlighting
   - Filter by file type
   - Search UI in sidebar

   **Implementation:**
   - Build search index on startup
   - Incremental indexing for new files
   - Client-side search optimization
   - Keyboard shortcut: `/` to search

2. **PDF Export**
   - Export single page to PDF
   - Export entire plan as PDF
   - Custom styling for PDF
   - Embed images and diagrams
   - Print-optimized layout

   **Implementation:**
   - Use Go's PDF libraries
   - Generate PDF server-side
   - Download via /export endpoint
   - Preserve formatting

3. **Advanced Navigation**
   - Breadcrumb refinements
   - Full-text TOC search
   - Section bookmarks
   - Jump to phase
   - Reading progress tracking

4. **Extended Markdown Support**
   - Callouts/alerts
   - Footnotes
   - Subscript/superscript
   - Definition lists
   - Custom containers

5. **Improved Accessibility**
   - WCAG AAA compliance
   - Enhanced screen reader support
   - Better color contrast
   - Focus management
   - Skip links

6. **Configuration File**
   - `mdview.config.json` support
   - Custom theme colors
   - Default port setting
   - Markdown extensions
   - Navigation customization

**Quality Targets:**
- Test coverage: > 85%
- WCAG AAA compliance verified
- Zero security issues
- < 5 open bugs

**Acceptance Criteria:**
- Search works on all file types
- PDF export preserves formatting
- Config file customization works
- All accessibility tests pass
- Performance targets maintained

---

### v1.1.0 - Plugin System (Q3 2026)

**Timeline:** 8-10 weeks after v1.0.0

**Goals:**
- Extensibility
- Community contributions
- Custom renderers
- Plugin marketplace

**Features:**

1. **Plugin Architecture**
   - Plugin interface definition
   - Plugin loading mechanism
   - Plugin configuration
   - Plugin security model

2. **Custom Renderers**
   - Plugin renderer interface
   - Example: SVG diagrams renderer
   - Example: Math renderer
   - Example: Custom table styles

3. **Plugin Development Kit**
   - Plugin template
   - Documentation
   - Example plugins
   - Testing utilities

4. **Plugin Registry**
   - Central plugin listing
   - Version management
   - Dependency tracking
   - Security review process

**Acceptance Criteria:**
- At least 3 example plugins
- Plugin documentation complete
- Community plugin submission process defined
- Security review completed

---

### v2.0.0 - Advanced Features (Q4 2026+)

**Timeline:** Flexible, based on community needs

**Potential Features:**

1. **Multi-Document Projects**
   - Linking between documents
   - Cross-document references
   - Project-wide search
   - Dependency visualization

2. **Collaboration Features**
   - Live document sharing
   - Comments and annotations
   - Change tracking
   - Multi-user navigation (future)

3. **Advanced Plan Features**
   - Gantt chart visualization
   - Dependency graphs
   - Timeline views
   - Milestone tracking

4. **Version Control Integration**
   - Git history viewing
   - Diff visualization
   - Commit annotations
   - Branch browsing

5. **Mobile Apps**
   - iOS app (future)
   - Android app (future)
   - Offline-first sync
   - Native integrations

---

## Feature Backlog (Not Prioritized)

### High Priority
- [ ] Search functionality (v1.0)
- [ ] PDF export (v1.0)
- [ ] Configuration file support (v1.0)
- [ ] Performance benchmarking
- [ ] WCAG AAA compliance
- [ ] Extended markdown syntax
- [ ] Improved error messages

### Medium Priority
- [ ] Custom theme builder UI
- [ ] Diagram rendering alternatives (SVG support)
- [ ] LaTeX/Math rendering
- [ ] Code block copy button
- [ ] Syntax highlighting options
- [ ] Plugin system (v1.1)
- [ ] CLI argument validation

### Low Priority
- [ ] Web app version (no build required)
- [ ] Docker image
- [ ] Reverse proxy support
- [ ] Analytics (optional, privacy-preserving)
- [ ] Internationalization
- [ ] Mobile apps (v2.0)
- [ ] Database backend support

---

## Known Issues & Workarounds

### Current Issues

| Issue | Severity | Workaround | Target Fix |
|-------|----------|-----------|-----------|
| Mermaid renders on CDN delay | Medium | Include diagram alt text | v1.0 |
| Plan parser edge cases | Low | Use standard formats | v0.2 |
| Mobile sidebar scroll | Low | Use touch with momentum | v0.2 |
| Large file rendering | Medium | Split documents | v0.2 |
| Windows path handling | Low | Use forward slashes | v0.2 |

---

## Development Milestones

### Completed
- [x] Core markdown rendering
- [x] Plan parsing (6 formats)
- [x] Theme system
- [x] Mobile responsiveness
- [x] Keyboard shortcuts
- [x] Syntax highlighting
- [x] Mermaid integration

### In Progress
- [ ] Community feedback incorporation
- [ ] Performance optimization
- [ ] Bug fixing and stability

### Upcoming
- [ ] Full-text search implementation
- [ ] PDF export feature
- [ ] Configuration system
- [ ] Accessibility audit
- [ ] Plugin architecture
- [ ] Extended markdown support

---

## Dependencies & Constraints

### Technical Constraints
- Go 1.22+ required
- Single binary distribution (no external assets)
- No database required (read-only file serving)
- Stateless HTTP server

### Resource Constraints
- Small binary size (< 20MB preferred)
- Low memory footprint (< 100MB typical)
- Fast startup (< 1s target)
- No external dependencies for core features

### Timeline Constraints
- v1.0.0 required before plugin system
- Search must precede collaborative features
- Accessibility audit required before major release
- Security review before production recommendation

---

## Success Criteria by Version

### v0.1.0
- GitHub repo created and documented
- Initial users can view markdown
- Theme system functional
- Mobile layouts work
- Status: ✓ Complete

### v0.2.0
- Zero critical bugs reported
- Performance improved > 20%
- Mobile UX refined
- Community feedback incorporated
- Target: 4-6 weeks

### v1.0.0
- Search fully functional
- PDF export working
- Accessibility AAA compliant
- > 100 GitHub stars
- Production-ready
- Target: 10-12 weeks from v0.2.0

### v1.1.0
- At least 3 example plugins
- Plugin documentation complete
- First community plugin
- Community engagement active
- Target: 8-10 weeks from v1.0.0

---

## Community & Contribution Path

### Version 0.1.0 - v0.2.0
- Bug reports and feedback collection
- Issue triage and labeling
- Good first issue identification
- Documentation improvements accepted
- Small feature PRs welcome

### Version 0.2.0 - v1.0.0
- Feature proposal discussions
- Architecture feedback welcomed
- Larger feature contributions considered
- Maintenance contributor program starts
- Reviewer recruitment begins

### Version 1.0.0 - v1.1.0
- Plugin development community
- Maintainer team expansion
- Sponsorship/funding consideration
- Official plugin registry
- Community governance model

---

## Metrics & Tracking

### Adoption Metrics
- GitHub stars
- Monthly downloads
- Social media mentions
- Community discussions
- Enterprise usage

### Quality Metrics
- Test coverage percentage
- Security issue count
- Bug fix turnaround time
- Feature request response time
- Performance benchmarks

### Community Metrics
- Open issues/PRs
- Community contributors
- Plugin submissions
- Documentation contributions
- Stack Overflow mentions

---

## Related Documentation

- [Project Overview & PDR](./project-overview-pdr.md)
- [System Architecture](./system-architecture.md)
- [Codebase Summary](./codebase-summary.md)
- [Code Standards](./code-standards.md)
