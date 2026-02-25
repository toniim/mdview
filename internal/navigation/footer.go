package navigation

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
)

// NavPhase represents a phase in the navigation context.
type NavPhase struct {
	Phase  int
	Name   string
	Status string
	File   string
	Anchor string
}

// NavContext holds navigation state for a file.
type NavContext struct {
	PlanInfo     *PlanInfo
	CurrentIndex int
	Prev         *NavPhase
	Next         *NavPhase
	AllPhases    []NavPhase
}

// GetNavigationContext returns navigation context for a file.
func GetNavigationContext(filePath string) *NavContext {
	planInfo := DetectPlan(filePath)

	if !planInfo.IsPlan {
		return &NavContext{PlanInfo: planInfo, CurrentIndex: -1}
	}

	phaseMeta := ParsePlanTable(planInfo.PlanFile)

	// Build all phases list including plan.md
	allPhases := []NavPhase{
		{Phase: 0, Name: "Plan Overview", Status: "overview", File: planInfo.PlanFile},
	}
	for _, pm := range phaseMeta {
		allPhases = append(allPhases, NavPhase{
			Phase:  pm.Phase,
			Name:   pm.Name,
			Status: pm.Status,
			File:   pm.File,
			Anchor: pm.Anchor,
		})
	}

	// Find current file
	normalizedPath := filepath.Clean(filePath)
	currentIndex := -1
	for i, p := range allPhases {
		if filepath.Clean(p.File) == normalizedPath {
			currentIndex = i
			break
		}
	}

	var prev, next *NavPhase
	if currentIndex > 0 {
		p := allPhases[currentIndex-1]
		prev = &p
	}
	if currentIndex >= 0 && currentIndex < len(allPhases)-1 {
		n := allPhases[currentIndex+1]
		next = &n
	}

	return &NavContext{
		PlanInfo:     planInfo,
		CurrentIndex: currentIndex,
		Prev:         prev,
		Next:         next,
		AllPhases:    allPhases,
	}
}

// GenerateNavFooter generates prev/next navigation footer HTML.
func GenerateNavFooter(filePath, rootDir string) string {
	navCtx := GetNavigationContext(filePath)

	if navCtx.Prev == nil && navCtx.Next == nil {
		return ""
	}

	prevHTML := "<span></span>"
	if navCtx.Prev != nil {
		_, err := os.Stat(navCtx.Prev.File)
		if err == nil {
			prevHTML = fmt.Sprintf(`<a href="%s" class="nav-prev">
      <span class="nav-arrow">&larr;</span>
      <span class="nav-label">%s</span>
    </a>`, ViewURL(navCtx.Prev.File, rootDir), html.EscapeString(navCtx.Prev.Name))
		} else {
			prevHTML = fmt.Sprintf(`<span class="nav-prev nav-unavailable" title="Phase planned but not yet implemented">
      <span class="nav-arrow">&larr;</span>
      <span class="nav-label">%s</span>
      <span class="nav-badge">Planned</span>
    </span>`, html.EscapeString(navCtx.Prev.Name))
		}
	}

	nextHTML := "<span></span>"
	if navCtx.Next != nil {
		_, err := os.Stat(navCtx.Next.File)
		if err == nil {
			nextHTML = fmt.Sprintf(`<a href="%s" class="nav-next">
      <span class="nav-label">%s</span>
      <span class="nav-arrow">&rarr;</span>
    </a>`, ViewURL(navCtx.Next.File, rootDir), html.EscapeString(navCtx.Next.Name))
		} else {
			nextHTML = fmt.Sprintf(`<span class="nav-next nav-unavailable" title="Phase planned but not yet implemented">
      <span class="nav-label">%s</span>
      <span class="nav-badge">Planned</span>
      <span class="nav-arrow">&rarr;</span>
    </span>`, html.EscapeString(navCtx.Next.Name))
		}
	}

	return fmt.Sprintf(`<footer class="nav-footer">
      %s
      %s
    </footer>`, prevHTML, nextHTML)
}
