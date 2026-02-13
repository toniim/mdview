package server

import (
	"encoding/json"
	"fmt"
	"html"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bilabl/mdview/internal/navigation"
	"github.com/bilabl/mdview/internal/renderer"
	"github.com/bilabl/mdview/internal/templates"
)

var mimeTypes = map[string]string{
	".html": "text/html",
	".css":  "text/css",
	".js":   "application/javascript",
	".json": "application/json",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".webp": "image/webp",
	".ico":  "image/x-icon",
	".md":   "text/markdown",
	".txt":  "text/plain",
	".pdf":  "application/pdf",
}

var viewableCodeExts = map[string]bool{
	".ts": true, ".tsx": true, ".js": true, ".jsx": true,
	".css": true, ".html": true, ".json": true, ".py": true, ".go": true,
	".rs": true, ".sh": true, ".yaml": true, ".yml": true, ".toml": true,
	".sql": true, ".java": true, ".kt": true, ".swift": true,
	".rb": true, ".php": true, ".vue": true, ".svelte": true,
	".cjs": true, ".mjs": true, ".env": true,
}

func getMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if mt, ok := mimeTypes[ext]; ok {
		return mt
	}
	return "application/octet-stream"
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/view", s.handleView)
	mux.HandleFunc("/browse", s.handleBrowse)
	mux.HandleFunc("/assets/", s.handleAssets)
	mux.HandleFunc("/file/", s.handleFile)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		sendError(w, 404, "Not found")
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head>
  <title>Markdown Novel Viewer</title>
  <style>
    body { font-family: system-ui; max-width: 600px; margin: 2rem auto; padding: 1rem; }
    h1 { color: #8b4513; }
    code { background: #f5f5f5; padding: 0.2rem 0.4rem; border-radius: 3px; }
    .routes { background: #faf8f3; padding: 1rem; border-radius: 8px; margin: 1rem 0; }
  </style>
</head>
<body>
  <h1>Markdown Novel Viewer</h1>
  <p>A calm, book-like viewer for markdown files.</p>
  <div class="routes">
    <h3>Routes</h3>
    <ul>
      <li><code>/view?file=/path/to/file.md</code> - View markdown</li>
      <li><code>/browse?dir=/path/to/dir</code> - Browse directory</li>
    </ul>
  </div>
  <p>Use <code>mdview serve &lt;path&gt;</code> to start viewing files.</p>
</body>
</html>`)
}

func (s *Server) handleView(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		sendError(w, 400, "Missing ?file= parameter")
		return
	}

	if !isPathSafe(filePath, s.config.AllowedDirs) {
		sendError(w, 403, "Access denied")
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sendError(w, 404, "File not found")
		return
	}

	pageHTML, err := s.generateFullPage(filePath)
	if err != nil {
		sendError(w, 500, "Error rendering file")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, pageHTML)
}

func (s *Server) handleBrowse(w http.ResponseWriter, r *http.Request) {
	dirPath := r.URL.Query().Get("dir")
	if dirPath == "" {
		sendError(w, 400, "Missing ?dir= parameter")
		return
	}

	if !isPathSafe(dirPath, s.config.AllowedDirs) {
		sendError(w, 403, "Access denied")
		return
	}

	info, err := os.Stat(dirPath)
	if err != nil || !info.IsDir() {
		sendError(w, 404, "Directory not found")
		return
	}

	htmlContent, err := renderDirectoryBrowser(dirPath)
	if err != nil {
		sendError(w, 500, "Error listing directory")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, htmlContent)
}

func (s *Server) handleAssets(w http.ResponseWriter, r *http.Request) {
	// Strip leading /assets/ to get relative path within embedded FS
	relPath := strings.TrimPrefix(r.URL.Path, "/")
	if strings.Contains(relPath, "..") {
		sendError(w, 403, "Access denied")
		return
	}

	data, err := templates.Assets.ReadFile(relPath)
	if err != nil {
		sendError(w, 404, "Asset not found")
		return
	}

	mt := getMimeType(relPath)
	w.Header().Set("Content-Type", mt)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(data)
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimPrefix(r.URL.Path, "/file/")
	// URL decode
	filePath, err := url.PathUnescape(filePath)
	if err != nil {
		sendError(w, 400, "Invalid path")
		return
	}

	if !isPathSafe(filePath, s.config.AllowedDirs) {
		sendError(w, 403, "Access denied")
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sendError(w, 404, "File not found")
		return
	}

	mt := getMimeType(filePath)
	w.Header().Set("Content-Type", mt)
	http.ServeFile(w, r, filePath)
}

func (s *Server) generateFullPage(filePath string) (string, error) {
	isCode := renderer.IsCodeFile(filePath)

	var result *renderer.RenderResult
	var err error
	if isCode {
		result, err = renderer.RenderCodeFile(filePath)
	} else {
		result, err = renderer.RenderMarkdownFile(filePath)
	}
	if err != nil {
		return "", err
	}

	tocHTML := renderer.RenderTOCHTML(result.TOC)

	var navSidebar, navFooter string
	var planInfo *navigation.PlanInfo
	var navCtx *navigation.NavContext

	if !isCode {
		navSidebar = navigation.GenerateNavSidebar(filePath)
		navFooter = navigation.GenerateNavFooter(filePath)
		planInfo = navigation.DetectPlan(filePath)
		navCtx = navigation.GetNavigationContext(filePath)
	} else {
		planInfo = &navigation.PlanInfo{}
		navCtx = &navigation.NavContext{}
	}

	// Generate back button
	parentDir := filepath.Dir(filePath)
	backButton := fmt.Sprintf(`<a href="/browse?dir=%s" class="icon-btn back-btn" title="Back to folder">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M19 12H5M12 19l-7-7 7-7"/>
      </svg>
    </a>`, url.QueryEscape(parentDir))

	// Generate header nav
	var headerNav string
	if navCtx.Prev != nil || navCtx.Next != nil {
		var prevBtn, nextBtn string
		if navCtx.Prev != nil {
			if _, err := os.Stat(navCtx.Prev.File); err == nil {
				prevBtn = fmt.Sprintf(`<a href="/view?file=%s" class="header-nav-btn prev" title="%s">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 18l-6-6 6-6"/></svg>
          <span>Prev</span>
        </a>`, url.QueryEscape(navCtx.Prev.File), html.EscapeString(navCtx.Prev.Name))
			}
		}
		if navCtx.Next != nil {
			if _, err := os.Stat(navCtx.Next.File); err == nil {
				nextBtn = fmt.Sprintf(`<a href="/view?file=%s" class="header-nav-btn next" title="%s">
          <span>Next</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 18l6-6-6-6"/></svg>
        </a>`, url.QueryEscape(navCtx.Next.File), html.EscapeString(navCtx.Next.Name))
			}
		}
		headerNav = fmt.Sprintf(`<div class="header-nav">%s%s</div>`, prevBtn, nextBtn)
	}

	hasPlan := ""
	if planInfo.IsPlan {
		hasPlan = "has-plan"
	}

	fm, _ := json.Marshal(result.Frontmatter)

	data := &templates.PageData{
		Title:       result.Title,
		Content:     templates.ToHTML(result.HTML),
		TOC:         templates.ToHTML(tocHTML),
		NavSidebar:  templates.ToHTML(navSidebar),
		NavFooter:   templates.ToHTML(navFooter),
		HasPlan:     hasPlan,
		Frontmatter: string(fm),
		BackButton:  templates.ToHTML(backButton),
		HeaderNav:   templates.ToHTML(headerNav),
	}

	return templates.RenderPage(data)
}

func sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><title>Error %d</title></head>
<body style="font-family: system-ui; padding: 2rem;">
  <h1>Error %d</h1>
  <p>%s</p>
</body>
</html>`, code, code, html.EscapeString(message))
}

// File icon mapping
var fileIcons = map[string]string{
	".md": "\U0001F4C4", ".txt": "\U0001F4DD", ".json": "\U0001F4CB",
	".js": "\U0001F4DC", ".cjs": "\U0001F4DC", ".mjs": "\U0001F4DC",
	".ts": "\U0001F4D8", ".tsx": "\U0001F4D8", ".jsx": "\U0001F4DC",
	".css": "\U0001F3A8", ".html": "\U0001F310",
	".py": "\U0001F40D", ".go": "\U0001F535", ".rs": "\U0001F980",
	".vue": "\U0001F49A", ".svelte": "\U0001F9E1",
	".png": "\U0001F5BC\uFE0F", ".jpg": "\U0001F5BC\uFE0F", ".jpeg": "\U0001F5BC\uFE0F",
	".gif": "\U0001F5BC\uFE0F", ".svg": "\U0001F5BC\uFE0F",
	".pdf": "\U0001F4D5",
	".yaml": "\u2699\uFE0F", ".yml": "\u2699\uFE0F", ".toml": "\u2699\uFE0F",
	".env": "\U0001F510", ".sh": "\U0001F4BB", ".bash": "\U0001F4BB",
}

func getFileIcon(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if icon, ok := fileIcons[ext]; ok {
		return icon
	}
	return "\U0001F4C4"
}

func renderDirectoryBrowser(dirPath string) (string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	baseName := filepath.Base(dirPath)
	displayPath := dirPath
	if len(displayPath) > 50 {
		displayPath = "..." + displayPath[len(displayPath)-47:]
	}

	type dirEntry struct {
		Name  string
		IsDir bool
	}

	var dirs, files []dirEntry
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") || name == "deprecated" {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.IsDir() {
			dirs = append(dirs, dirEntry{Name: name, IsDir: true})
		} else {
			files = append(files, dirEntry{Name: name})
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	var listHTML strings.Builder

	// Parent directory
	parentDir := filepath.Dir(dirPath)
	if parentDir != dirPath {
		fmt.Fprintf(&listHTML, `<li class="dir-item parent">
      <a href="/browse?dir=%s">
        <span class="icon">📁</span>
        <span class="name">..</span>
      </a>
    </li>`, url.QueryEscape(parentDir))
	}

	// Directories
	for _, d := range dirs {
		fullPath := filepath.Join(dirPath, d.Name)
		fmt.Fprintf(&listHTML, `<li class="dir-item folder">
      <a href="/browse?dir=%s">
        <span class="icon">📁</span>
        <span class="name">%s/</span>
      </a>
    </li>`, url.QueryEscape(fullPath), html.EscapeString(d.Name))
	}

	// Files - all text files are viewable via /view, binary files via /file
	binaryExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".webp": true,
		".ico": true, ".svg": true, ".pdf": true, ".zip": true, ".tar": true,
		".gz": true, ".bz2": true, ".xz": true, ".7z": true, ".rar": true,
		".exe": true, ".dll": true, ".so": true, ".dylib": true, ".bin": true,
		".woff": true, ".woff2": true, ".ttf": true, ".otf": true, ".eot": true,
		".mp3": true, ".mp4": true, ".wav": true, ".avi": true, ".mov": true,
		".webm": true, ".ogg": true, ".flac": true,
	}
	for _, f := range files {
		fullPath := filepath.Join(dirPath, f.Name)
		icon := getFileIcon(f.Name)
		ext := strings.ToLower(filepath.Ext(f.Name))
		isMarkdown := ext == ".md"

		if binaryExts[ext] {
			// Binary files: serve directly
			fmt.Fprintf(&listHTML, `<li class="dir-item file">
        <a href="/file/%s" target="_blank">
          <span class="icon">%s</span>
          <span class="name">%s</span>
        </a>
      </li>`, url.PathEscape(fullPath), icon, html.EscapeString(f.Name))
		} else {
			// All text files: view in reader
			cls := "file code"
			if isMarkdown {
				cls = "file markdown"
			}
			fmt.Fprintf(&listHTML, `<li class="dir-item %s">
        <a href="/view?file=%s">
          <span class="icon">%s</span>
          <span class="name">%s</span>
        </a>
      </li>`, cls, url.QueryEscape(fullPath), icon, html.EscapeString(f.Name))
		}
	}

	if len(dirs) == 0 && len(files) == 0 {
		listHTML.WriteString(`<li class="empty">This directory is empty</li>`)
	}

	// Read CSS from embedded assets
	css := ""
	if data, err := fs.ReadFile(templates.Assets, "assets/directory-browser.css"); err == nil {
		css = string(data)
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>📁 %s</title>
  <style>%s</style>
</head>
<body>
  <div class="container">
    <header>
      <h1>📁 %s</h1>
      <p class="path">%s</p>
    </header>
    <ul class="file-list">%s</ul>
    <footer>
      <p>%d folder%s, %d file%s</p>
    </footer>
  </div>
</body>
</html>`,
		html.EscapeString(baseName),
		css,
		html.EscapeString(baseName),
		html.EscapeString(displayPath),
		listHTML.String(),
		len(dirs), plural(len(dirs)),
		len(files), plural(len(files)),
	), nil
}

func plural(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}

