package renderer

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// RenderResult holds the output of rendering a file.
type RenderResult struct {
	HTML        string
	TOC         []TOCItem
	Frontmatter map[string]interface{}
	Title       string
}

// newGoldmark creates a configured Goldmark instance.
func newGoldmark() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true),
					chromahtml.WithLineNumbers(false),
				),
				highlighting.WithGuessLanguage(true),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldhtml.WithUnsafe(),
			goldhtml.WithHardWraps(),
		),
	)
}

// RenderMarkdownFile renders a markdown file to HTML with TOC and frontmatter.
func RenderMarkdownFile(filePath string) (*RenderResult, error) {
	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	basePath := filepath.Dir(filePath)

	// Parse frontmatter
	var fm map[string]interface{}
	content, err := parseFrontmatter(rawContent, &fm)
	if err != nil {
		// If frontmatter parsing fails, use raw content
		content = rawContent
		fm = make(map[string]interface{})
	}

	// Pre-process mermaid blocks into raw HTML before goldmark rendering
	contentStr := preprocessMermaid(string(content))

	// Resolve image paths
	contentStr = resolveImages(contentStr, basePath)

	// Render markdown
	md := newGoldmark()
	var buf bytes.Buffer
	if err := md.Convert([]byte(contentStr), &buf); err != nil {
		return nil, err
	}

	htmlStr := buf.String()

	// Add heading IDs for phases
	htmlStr = addPhaseHeadingIDs(htmlStr)

	// Add phase table anchors
	htmlStr = addPhaseTableAnchors(htmlStr)

	// Generate TOC
	toc := generateTOC(htmlStr)

	// Extract title
	title := extractTitle(fm, htmlStr, filePath)

	return &RenderResult{
		HTML:        htmlStr,
		TOC:         toc,
		Frontmatter: fm,
		Title:       title,
	}, nil
}

func parseFrontmatter(content []byte, fm *map[string]interface{}) ([]byte, error) {
	rest, err := frontmatter.Parse(bytes.NewReader(content), fm)
	if err != nil {
		return content, err
	}
	return rest, nil
}

func extractTitle(fm map[string]interface{}, htmlContent, filePath string) string {
	if t, ok := fm["title"]; ok {
		if s, ok := t.(string); ok && s != "" {
			return s
		}
	}

	// Find first h1
	if idx := strings.Index(htmlContent, "<h1"); idx >= 0 {
		start := strings.Index(htmlContent[idx:], ">")
		if start >= 0 {
			end := strings.Index(htmlContent[idx+start:], "</h1>")
			if end >= 0 {
				title := htmlContent[idx+start+1 : idx+start+end]
				title = stripHTMLTags(title)
				return strings.TrimSpace(title)
			}
		}
	}

	return strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
}

func stripHTMLTags(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// GenerateChromaCSS generates CSS for a given Chroma style.
func GenerateChromaCSS(styleName string) string {
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}
	formatter := chromahtml.New(chromahtml.WithClasses(true))
	var buf bytes.Buffer
	formatter.WriteCSS(&buf, style)
	return buf.String()
}
