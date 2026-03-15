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
	mux.HandleFunc("/assets/", s.handleAssets)
	mux.HandleFunc("/file/", s.handleFile)
	mux.HandleFunc("/", s.handlePath)
}

// handlePath is the unified handler: stat the path, directory → browse, file → view.
func (s *Server) handlePath(w http.ResponseWriter, r *http.Request) {
	relPath := strings.TrimPrefix(r.URL.Path, "/")

	// Root path → browse root directory
	if relPath == "" {
		s.serveBrowse(w, s.config.RootDir)
		return
	}

	decoded, err := url.PathUnescape(relPath)
	if err != nil {
		sendError(w, 400, "Invalid path")
		return
	}

	fullPath := filepath.Join(s.config.RootDir, decoded)

	if !isPathSafe(fullPath, s.config.AllowedDirs) {
		sendError(w, 403, "Access denied")
		return
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		sendError(w, 404, "Not found")
		return
	}

	if info.IsDir() {
		s.serveBrowse(w, fullPath)
		return
	}

	pageHTML, err := s.generateFullPage(fullPath)
	if err != nil {
		sendError(w, 500, "Error rendering file")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, pageHTML)
}

func (s *Server) serveBrowse(w http.ResponseWriter, dirPath string) {
	if !isPathSafe(dirPath, s.config.AllowedDirs) {
		sendError(w, 403, "Access denied")
		return
	}

	info, err := os.Stat(dirPath)
	if err != nil || !info.IsDir() {
		sendError(w, 404, "Directory not found")
		return
	}

	htmlContent, err := s.renderDirectoryBrowser(dirPath)
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
	rawPath := strings.TrimPrefix(r.URL.Path, "/file/")
	filePath, err := url.PathUnescape(rawPath)
	if err != nil {
		sendError(w, 400, "Invalid path")
		return
	}

	// If path is not absolute, resolve relative to root dir
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(s.config.RootDir, filePath)
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

	rootDir := s.config.RootDir

	var navSidebar, navFooter string
	var planInfo *navigation.PlanInfo
	var navCtx *navigation.NavContext

	if !isCode {
		navSidebar = navigation.GenerateNavSidebar(filePath, rootDir)
		navFooter = navigation.GenerateNavFooter(filePath, rootDir)
		planInfo = navigation.DetectPlan(filePath)
		navCtx = navigation.GetNavigationContext(filePath)
	} else {
		planInfo = &navigation.PlanInfo{}
		navCtx = &navigation.NavContext{}
	}

	// Generate home + back buttons
	homeButton := `<a href="/" class="icon-btn home-btn" title="Home">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
        <polyline points="9 22 9 12 15 12 15 22"/>
      </svg>
    </a>`
	parentDir := filepath.Dir(filePath)
	backURL := navigation.BrowseURL(parentDir, rootDir)
	backButton := fmt.Sprintf(`%s<a href="%s" class="icon-btn back-btn" title="Back to folder">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M19 12H5M12 19l-7-7 7-7"/>
      </svg>
    </a>`, homeButton, backURL)

	// Generate header nav
	var headerNav string
	if navCtx.Prev != nil || navCtx.Next != nil {
		var prevBtn, nextBtn string
		if navCtx.Prev != nil {
			if _, err := os.Stat(navCtx.Prev.File); err == nil {
				prevBtn = fmt.Sprintf(`<a href="%s" class="header-nav-btn prev" title="%s">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 18l-6-6 6-6"/></svg>
          <span>Prev</span>
        </a>`, navigation.ViewURL(navCtx.Prev.File, rootDir), html.EscapeString(navCtx.Prev.Name))
			}
		}
		if navCtx.Next != nil {
			if _, err := os.Stat(navCtx.Next.File); err == nil {
				nextBtn = fmt.Sprintf(`<a href="%s" class="header-nav-btn next" title="%s">
          <span>Next</span>
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 18l6-6-6-6"/></svg>
        </a>`, navigation.ViewURL(navCtx.Next.File, rootDir), html.EscapeString(navCtx.Next.Name))
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

func (s *Server) renderDirectoryBrowser(dirPath string) (string, error) {
	rootDir := s.config.RootDir
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	baseName := filepath.Base(dirPath)
	// Show relative path from root
	displayPath, _ := filepath.Rel(rootDir, dirPath)
	if displayPath == "." {
		displayPath = "/"
	} else {
		displayPath = "/" + displayPath
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

	// Parent directory (only if not already at root)
	parentDir := filepath.Dir(dirPath)
	cleanDir := filepath.Clean(dirPath)
	cleanRoot := filepath.Clean(rootDir)
	if cleanDir != cleanRoot && parentDir != dirPath {
		parentURL := navigation.BrowseURL(parentDir, rootDir)
		fmt.Fprintf(&listHTML, `<li class="dir-item parent">
      <a href="%s">
        <span class="icon">📁</span>
        <span class="name">..</span>
      </a>
    </li>`, parentURL)
	}

	// Directories
	for _, d := range dirs {
		fullPath := filepath.Join(dirPath, d.Name)
		fmt.Fprintf(&listHTML, `<li class="dir-item folder">
      <a href="%s">
        <span class="icon">📁</span>
        <span class="name">%s/</span>
      </a>
    </li>`, navigation.BrowseURL(fullPath, rootDir), html.EscapeString(d.Name))
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
			// Binary files: serve via /file/ with relative path
			relFile, _ := filepath.Rel(rootDir, fullPath)
			fmt.Fprintf(&listHTML, `<li class="dir-item file">
        <a href="/file/%s" target="_blank">
          <span class="icon">%s</span>
          <span class="name">%s</span>
        </a>
      </li>`, navigation.EncodePath(relFile), icon, html.EscapeString(f.Name))
		} else {
			// All text files: view in reader
			cls := "file code"
			if isMarkdown {
				cls = "file markdown"
			}
			fmt.Fprintf(&listHTML, `<li class="dir-item %s">
        <a href="%s">
          <span class="icon">%s</span>
          <span class="name">%s</span>
        </a>
      </li>`, cls, navigation.ViewURL(fullPath, rootDir), icon, html.EscapeString(f.Name))
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

	homeLink := ""
	if cleanDir != cleanRoot {
		homeLink = `<a href="/" class="home-link" title="Home">🏠 Home</a>`
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
      %s
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
		homeLink,
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

