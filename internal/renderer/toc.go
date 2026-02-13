package renderer

import (
	"fmt"
	"regexp"
	"strings"
)

// TOCItem represents a heading in the table of contents.
type TOCItem struct {
	Level int
	ID    string
	Text  string
}

var (
	headingRegex      = regexp.MustCompile(`<h([1-3])[^>]*id="([^"]+)"[^>]*>([^<]+)</h[1-3]>`)
	rawHeadingRegex   = regexp.MustCompile(`(?i)<h([1-6])>([^<]+)</h[1-6]>`)
	phaseHeadingRegex = regexp.MustCompile(`(?i)^Phase\s*(\d+)[:\s]+(.+)`)
	phaseTableRegex   = regexp.MustCompile(`(?i)<tr>\s*<td>(\d{2})</td>\s*<td>([^<]+)</td>`)
)

func generateTOC(htmlContent string) []TOCItem {
	matches := headingRegex.FindAllStringSubmatch(htmlContent, -1)
	var items []TOCItem
	for _, m := range matches {
		level := 1
		if m[1] == "2" {
			level = 2
		} else if m[1] == "3" {
			level = 3
		}
		items = append(items, TOCItem{
			Level: level,
			ID:    m[2],
			Text:  strings.TrimSpace(m[3]),
		})
	}
	return items
}

// RenderTOCHTML generates HTML for the TOC sidebar.
func RenderTOCHTML(toc []TOCItem) string {
	if len(toc) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<ul class="toc-list">`)
	for _, item := range toc {
		indent := (item.Level - 1) * 12
		fmt.Fprintf(&b, `<li style="padding-left: %dpx"><a href="#%s">%s</a></li>`,
			indent, item.ID, item.Text)
	}
	b.WriteString("</ul>")
	return b.String()
}

// Slugify converts text to a URL-safe slug.
func Slugify(text string) string {
	s := strings.ToLower(text)
	// Replace non-alphanumeric with hyphens
	var result []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			result = append(result, c)
		} else {
			if len(result) > 0 && result[len(result)-1] != '-' {
				result = append(result, '-')
			}
		}
	}
	// Trim leading/trailing hyphens
	return strings.Trim(string(result), "-")
}

// addPhaseHeadingIDs adds phase-specific IDs to headings that don't already have IDs.
func addPhaseHeadingIDs(htmlContent string) string {
	usedIDs := make(map[string]bool)

	return rawHeadingRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := rawHeadingRegex.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		level := sub[1]
		text := sub[2]

		var id string
		phaseMatch := phaseHeadingRegex.FindStringSubmatch(text)
		if phaseMatch != nil {
			phaseNum := phaseMatch[1]
			phaseName := strings.TrimSpace(phaseMatch[2])
			// Pad phase number
			if len(phaseNum) == 1 {
				phaseNum = "0" + phaseNum
			}
			id = fmt.Sprintf("phase-%s-%s", phaseNum, Slugify(phaseName))
		} else {
			id = Slugify(text)
		}

		// Deduplicate
		uniqueID := id
		counter := 1
		for usedIDs[uniqueID] {
			uniqueID = fmt.Sprintf("%s-%d", id, counter)
			counter++
		}
		usedIDs[uniqueID] = true

		return fmt.Sprintf(`<h%s id="%s">%s</h%s>`, level, uniqueID, text, level)
	})
}

// addPhaseTableAnchors adds anchor IDs to phase table rows.
func addPhaseTableAnchors(htmlContent string) string {
	usedIDs := make(map[string]bool)

	return phaseTableRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		sub := phaseTableRegex.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		phaseNum := sub[1]
		description := strings.TrimSpace(sub[2])
		id := fmt.Sprintf("phase-%s-%s", phaseNum, Slugify(description))

		uniqueID := id
		counter := 1
		for usedIDs[uniqueID] {
			uniqueID = fmt.Sprintf("%s-%d", id, counter)
			counter++
		}
		usedIDs[uniqueID] = true

		return fmt.Sprintf(`<tr id="%s"><td>%s</td><td>%s</td>`, uniqueID, phaseNum, description)
	})
}
