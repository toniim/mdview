package templates

import (
	"bytes"
	"html/template"
	"sync"
)

type PageData struct {
	Title       string
	Content     template.HTML
	TOC         template.HTML
	NavSidebar  template.HTML
	NavFooter   template.HTML
	HasPlan     string
	Frontmatter string
	BackButton  template.HTML
	HeaderNav   template.HTML
}

var (
	pageTmpl *template.Template
	tmplOnce sync.Once
)

func getTemplate() (*template.Template, error) {
	var err error
	tmplOnce.Do(func() {
		data, readErr := Assets.ReadFile("assets/template.html")
		if readErr != nil {
			err = readErr
			return
		}
		pageTmpl, err = template.New("page").Parse(string(data))
	})
	return pageTmpl, err
}

// ToHTML converts a raw string to template.HTML (trusted, unescaped).
func ToHTML(s string) template.HTML {
	return template.HTML(s)
}

func RenderPage(data *PageData) (string, error) {
	tmpl, err := getTemplate()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
