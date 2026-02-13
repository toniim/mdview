# Design Guidelines

## Philosophy

mdview follows a **calm, book-like aesthetic** that prioritizes reading experience over flashy features. The design is inspired by physical books and minimalist reading applications.

**Core Principles:**
1. **Clarity** - Clear typography, generous whitespace
2. **Calmness** - Muted colors, smooth animations, no distractions
3. **Efficiency** - One-click access to features, predictable interactions
4. **Accessibility** - Works for all users regardless of ability or device

## Color Palette

### Light Theme

```
Primary Text:     #2c2c2c (near-black for maximum readability)
Secondary Text:   #666666 (muted gray for metadata)
Background:       #faf8f3 (warm cream, reduces eye strain)
Accent:           #8b4513 (saddle brown, warm and book-like)
Border:           #e8e6df (light warm gray)
Code Background:  #f5f3ed (slightly darker cream)
Success:          #4a9d6f (muted green)
Warning:          #c97f4a (warm orange)
Error:            #c94c4c (muted red)
```

### Dark Theme

```
Primary Text:     #f0ede8 (off-white for contrast)
Secondary Text:   #a8a8a8 (muted gray)
Background:       #1a1a1a (true black for OLED efficiency)
Accent:           #d4a574 (golden brown, warm in darkness)
Border:           #3a3a3a (dark gray)
Code Background:  #242424 (dark gray)
Success:          #6bb992 (brightened green)
Warning:          #e5a166 (brightened orange)
Error:            #e86868 (brightened red)
```

### CSS Implementation

**Light theme** (default, `:root`):
```css
:root {
  --color-text-primary: #2c2c2c;
  --color-text-secondary: #666666;
  --color-background: #faf8f3;
  --color-accent: #8b4513;
  --color-border: #e8e6df;
  --color-code-bg: #f5f3ed;
  /* ... more variables */
}
```

**Dark theme** (`[data-theme="dark"]`):
```css
[data-theme="dark"] {
  --color-text-primary: #f0ede8;
  --color-text-secondary: #a8a8a8;
  --color-background: #1a1a1a;
  --color-accent: #d4a574;
  --color-border: #3a3a3a;
  --color-code-bg: #242424;
  /* ... more variables */
}
```

**Transitions:**
All color changes use `transition: color 200ms, background-color 200ms ease-out;` for smooth theme switching without page reload.

## Typography

### Font Stack

| Element | Font | Size | Weight | Usage |
|---------|------|------|--------|-------|
| Body | Inter | 16px (14px mobile) | 400 | Main text |
| Headings | Libre Baskerville | 1.5em-2.2em | 600 | Section headers |
| Code | JetBrains Mono | 14px | 400 | Code blocks, inline code |
| UI | Inter | 13px-14px | 500 | Buttons, labels |

### Font Loading

All fonts from **Google Fonts** (no self-hosted fonts for simplicity):

```html
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Libre+Baskerville:wght@400;600&family=JetBrains+Mono:wght@400&display=swap" rel="stylesheet">
```

### Line Heights & Spacing

| Element | Line Height | Margin Bottom |
|---------|------------|---------------|
| Paragraph | 1.6 | 1.2em |
| Heading (h1) | 1.3 | 0.8em |
| Heading (h2) | 1.4 | 0.7em |
| Heading (h3) | 1.5 | 0.6em |
| List items | 1.6 | 0.5em |
| Code block | 1.5 | 1.2em |

**Rationale:** Generous line-height (1.6) for long-form content reduces cognitive load and improves readability.

## Layout & Spacing

### Spacing Scale

```
0.25rem (4px)
0.5rem (8px)
1rem (16px)
1.5rem (24px)
2rem (32px)
3rem (48px)
4rem (64px)
```

**Rule:** Use multiples of 8px for consistent spacing.

### Grid System

- **Sidebar width:** 280px (default), 200-480px resizable
- **Content max-width:** 800px
- **Content margin:** 2rem horizontal
- **Content padding:** 2rem top/bottom

**Desktop Layout:**
```
┌──────────────────────────────────────────────────┐
│ Header (fixed, 60px)                              │
├─────────────────────────────────────────────────┤
│  Sidebar    │  Content Area (max 800px)  │       │
│  (280px)    │                            │ Footer│
│             │                            │ Nav   │
└──────────────────────────────────────────────────┘
```

**Tablet Layout (900px breakpoint):**
- Sidebar collapses to FAB
- Content takes full width
- Header adjusts

**Mobile Layout (600px breakpoint):**
- Full mobile stack
- Bottom sheet navigation
- Single column content

## Component Design

### Header

**Height:** 60px (fixed)
**Background:** Matches page background with subtle top border
**Contents:**
- Title (left-aligned)
- Font size control (S/M/L)
- Theme toggle
- Sidebar toggle (mobile)

**Interactions:**
- Font size persists in localStorage
- Theme toggle updates all colors instantly
- No page reload on toggle

### Sidebar

**Default State:** Visible, 280px width
**Resizable:** Drag handle on right edge (200-480px range)
**Behavior:**
- Toggle visibility with S key or button
- State persists in localStorage
- Smooth slide-out animation (200ms)
- Accordion groups of phases
- Status badges color-coded

**Accordion Groups:**
- Max 10 phases per group
- Current phase highlighted
- Expand/collapse per group
- State saved in localStorage with key: `reader:accordion:{planId}:{phaseId}`

### Content Area

**Typography:**
- Generous margins around text (2rem)
- Max width 800px for optimal line length
- Proper heading hierarchy (h1-h3)
- Strong contrast (WCAG AA minimum, AAA targeted)

**Code Blocks:**
- Language-specific syntax highlighting
- Line numbers (left-aligned, subtle)
- Copy button on hover
- Scroll horizontally if exceeds width
- Monospace font (JetBrains Mono)

**Images:**
- Responsive sizing (max 100% width)
- Center-aligned in content
- Light border in dark theme
- Lazy loading (native)

**Blockquotes:**
- Left border accent color (4px)
- Italic text
- Subtle background color
- Indented

**Tables:**
- Full width (up to content max-width)
- Horizontal scroll on mobile
- Alternating row colors (subtle)
- Clear header styling (background, bold)

### Mermaid Diagrams

**Styling:**
- Light theme: GitHub light colors
- Dark theme: GitHub dark colors
- Smooth transitions (200ms)
- Error message displayed inline

**Expandable:**
- Click diagram to expand to full viewport
- Esc key closes expanded view
- Maintains aspect ratio

### Navigation Footer

**Layout:** Sticky to bottom
**Contents:**
- Previous page link (left)
- Current phase/section (center)
- Next page link (right)
**Styling:**
- Subtle border top
- Muted text colors
- Disabled state for first/last pages

### Overlays

**Toast Notifications:**
- Bottom-right position
- Auto-dismiss after 4 seconds
- Stackable (multiple messages)
- Color-coded (success/warning/error)

**Shortcuts Modal:**
- Darkened backdrop
- Centered content
- Keyboard list with icons
- Close on Esc or click outside

**Expandable Content:**
- Full viewport width and height
- Fixed positioning
- Close button (top-right)
- Smooth fade-in/out

## Responsive Design

### Breakpoints

| Breakpoint | Device | Changes |
|-----------|--------|---------|
| 900px | Desktop → Tablet | Sidebar collapses to FAB |
| 768px | Tablet → Large Mobile | Adjustments for smaller screen |
| 600px | Mobile | Full mobile layout |

### Mobile-First Approach

**Base (mobile-first):**
- Full width content
- Single column
- Bottom sheet navigation
- FAB controls

**900px and up (desktop):**
- Add sidebar
- Multi-column layout
- Fixed header with controls
- Desktop-optimized spacing

### Touch-Friendly Design

**Mobile Interactions:**
- FAB size: 56px x 56px (material standard)
- Touch targets: min 44px x 44px
- Bottom sheet swipe-to-close
- Swipe-left for next page
- Swipe-right for previous page

## Animation & Transitions

### Smooth, Subtle Motion

**Timing:**
- Fast interactions: 150ms (hover, focus)
- Transitions: 200ms (theme, sidebar)
- Animations: 300ms (page load, expand)

**Easing:**
- ease-out for focus/hover
- ease-in-out for transitions
- cubic-bezier(0.23, 1, 0.320, 1) for smooth morphing

**No Animations on:**
- `prefers-reduced-motion: reduce` (respect user preference)
- Text selection interactions
- Scrolling

### Specific Animations

| Interaction | Duration | Easing |
|------------|----------|--------|
| Theme toggle | 200ms | ease-out |
| Sidebar open/close | 250ms | ease-out |
| Button hover | 150ms | ease-out |
| Progress bar | 100ms | linear |
| Mermaid render | 300ms | ease-in-out |
| Toast appear | 200ms | ease-out |
| Focus ring | 200ms | ease-out |

## Accessibility Design

### Color Contrast

**WCAG AA (minimum):**
- Normal text (14px+): 4.5:1
- Large text (18px+): 3:1

**Target (AAA):**
- All text: 7:1 or higher

**Testing:**
- Use WebAIM contrast checker
- Test in both light and dark themes
- Verify with color blindness simulations

### Focus Management

**Visible focus indicators:**
- 2px solid accent color border
- 2px offset from element
- Works on all interactive elements
- Visible in both light and dark themes

**Focus trapping:**
- Modal overlays trap focus
- Tab order follows visual hierarchy
- Esc closes overlays

### Semantic HTML

- Use `<nav>` for navigation
- Use `<main>` for primary content
- Use `<header>`, `<footer>` for sections
- Use `<button>` for buttons (not `<div>`)
- Use `<a>` for links
- ARIA labels for icon-only buttons

### Keyboard Navigation

**All features accessible via keyboard:**
- Tab: Navigate to next element
- Shift+Tab: Previous element
- Enter/Space: Activate button
- Arrow keys: Navigate lists
- Esc: Close modals
- Custom: T, S, ←, →, ?

## Print Styles

**Print Media:**
```css
@media print {
  body { font-size: 12pt; line-height: 1.5; }
  header, nav, .sidebar { display: none; }
  a { color: #000; text-decoration: underline; }
  code { background: #eee; padding: 2px 4px; }
  page-break-inside: avoid; /* Prevent breaking content */
}
```

**Considerations:**
- Remove navigation and header
- Optimize for black and white
- Increase line height for printed text
- Avoid page breaks in important sections

## Theming System

### How to Customize Colors

**Option 1: Modify CSS variables**
Edit `novel-theme-variables.css`:
```css
:root {
  --color-accent: #your-color;
  --color-background: #your-color;
  /* ... */
}
```

**Option 2: Configuration file (future v1.0)**
`mdview.config.json`:
```json
{
  "theme": {
    "light": { "accent": "#..." },
    "dark": { "accent": "#..." }
  }
}
```

### Adding New Theme Variants

1. Add CSS variable definitions
2. Create separate CSS file
3. Link in template.html
4. Select via configuration

**Example: Sepia theme**
```css
[data-theme="sepia"] {
  --color-text-primary: #3a2817;
  --color-background: #f4e8d8;
  --color-accent: #8b5a3c;
  /* ... */
}
```

## Dark Mode Detection

**Automatic detection:**
```javascript
const darkModeEnabled = window.matchMedia('(prefers-color-scheme: dark)').matches;
```

**User override:**
- Theme button toggles preference
- Stored in localStorage
- Persists across sessions

**CSS implementation:**
```css
/* Light theme (default) */
:root { --color-background: #faf8f3; }

/* Dark theme (preference or override) */
@media (prefers-color-scheme: dark) {
  :root { --color-background: #1a1a1a; }
}

/* Explicit override via data attribute */
[data-theme="dark"] { --color-background: #1a1a1a; }
```

## Icon Design

**Style:** Minimal, outline-based
**Size:** 24px (UI elements), 32px (buttons)
**Stroke width:** 2px
**Consistency:** All icons follow same visual language

**Icon Uses:**
- File type indicators (directory listing)
- Navigation arrows (prev/next)
- Status badges (completed/pending)
- Control buttons (theme, sidebar, expand)

**Accessibility:**
- Use with text labels (no icon-only except buttons with aria-label)
- High contrast with background
- Scale smoothly on different devices

## Micro-interactions

### Button Hover State
```css
button:hover {
  background-color: var(--color-accent);
  color: var(--color-background);
  transition: all 150ms ease-out;
}
```

### Link Hover State
```css
a:hover {
  text-decoration: underline;
  color: var(--color-accent);
}
```

### Sidebar Drag Handle
```css
.resize-handle:hover {
  background-color: var(--color-accent);
  opacity: 0.8;
  cursor: col-resize;
}
```

### Progress Bar
- Animates smoothly with scroll
- Uses accent color
- Subtle shadow for depth

## Design QA Checklist

Before every release, verify:

- [ ] Light and dark themes render correctly
- [ ] All text meets WCAG AA contrast (AAA preferred)
- [ ] Focus indicators visible on all elements
- [ ] Mobile layouts at 600px and 768px breakpoints
- [ ] Mermaid diagrams render and theme switch correctly
- [ ] Keyboard navigation complete (no mouse required)
- [ ] Print styles don't break layouts
- [ ] Animation performance (60 FPS on scroll)
- [ ] Emoji and special characters render correctly
- [ ] Images are responsive and accessible
- [ ] Code blocks copy correctly
- [ ] No horizontal scrolling on content (except code)
- [ ] Toast notifications don't overlap
- [ ] Sidebar resize doesn't break layout
- [ ] Time-based features (toast dismiss) work consistently

## Related Documentation

- [Code Standards - CSS Architecture](./code-standards.md#css-architecture)
- [System Architecture - Frontend Components](./system-architecture.md#frontend-components)
- [Codebase Summary - CSS System](./codebase-summary.md#css-architecture)
