package navigation

import (
	"fmt"
	"html"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type indexedPhase struct {
	Phase NavPhase
	Index int
}

type phaseGroup struct {
	Start  int
	End    int
	Phases []indexedPhase
}

// GenerateNavSidebar generates the plan navigation sidebar HTML.
func GenerateNavSidebar(filePath string) string {
	navCtx := GetNavigationContext(filePath)
	if !navCtx.PlanInfo.IsPlan {
		return ""
	}

	planName := filepath.Base(navCtx.PlanInfo.PlanDir)
	normalizedCurrent := filepath.Clean(filePath)

	var groups []phaseGroup
	var currentGroup []indexedPhase
	groupStart := 0

	for i, phase := range navCtx.AllPhases {
		if len(currentGroup) == 0 {
			groupStart = phase.Phase
		}
		currentGroup = append(currentGroup, indexedPhase{Phase: phase, Index: i})

		if len(currentGroup) == 10 || i == len(navCtx.AllPhases)-1 ||
			(phase.Phase%10 == 0 && phase.Phase != groupStart) {
			groups = append(groups, phaseGroup{
				Start:  groupStart,
				End:    phase.Phase,
				Phases: append([]indexedPhase{}, currentGroup...),
			})
			currentGroup = nil
		}
	}

	var groupsHTML strings.Builder
	for _, group := range groups {
		groupID := fmt.Sprintf("phase-group-%d-%d", group.Start, group.End)
		var groupLabel string
		if group.Start == 0 {
			groupLabel = "Overview"
		} else if group.Start == group.End {
			groupLabel = fmt.Sprintf("Phase %d", group.Start)
		} else {
			groupLabel = fmt.Sprintf("Phases %d-%d", group.Start, group.End)
		}

		groupBadge := getGroupBadge(group.Phases)

		var itemsHTML strings.Builder
		for _, ip := range group.Phases {
			phase := ip.Phase
			isActive := ip.Index == navCtx.CurrentIndex
			statusClass := strings.ReplaceAll(phase.Status, " ", "-")
			normalizedPhase := filepath.Clean(phase.File)
			isSameFile := normalizedPhase == normalizedCurrent

			// Check file existence
			_, fileErr := os.Stat(phase.File)
			fileExists := fileErr == nil

			if !fileExists {
				fmt.Fprintf(&itemsHTML, `
        <li class="phase-item unavailable" data-status="%s" title="Phase planned but not yet implemented">
          <span class="phase-link-disabled">
            <span class="status-dot %s"></span>
            <span class="phase-name">%s</span>
            <span class="unavailable-badge">Planned</span>
          </span>
        </li>`, statusClass, statusClass, html.EscapeString(phase.Name))
				continue
			}

			// Build href
			var href string
			isInlineSection := false
			if isSameFile && phase.Anchor != "" {
				href = "#" + phase.Anchor
				isInlineSection = true
			} else if phase.Anchor != "" {
				href = fmt.Sprintf("/view?file=%s#%s", url.QueryEscape(phase.File), phase.Anchor)
			} else {
				href = fmt.Sprintf("/view?file=%s", url.QueryEscape(phase.File))
			}

			dataAnchor := ""
			if phase.Anchor != "" {
				dataAnchor = fmt.Sprintf(` data-anchor="%s"`, phase.Anchor)
			}

			inlineSectionClass := ""
			if isInlineSection {
				inlineSectionClass = " inline-section"
			}

			activeClass := ""
			if isActive {
				activeClass = " active"
			}

			// Type icon
			var typeIcon string
			if isInlineSection {
				typeIcon = `<svg class="phase-type-icon" viewBox="0 0 16 16" fill="currentColor"><path d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-.5 9.45a.75.75 0 01-1.06-1.06l-1.25 1.25a2 2 0 01-2.83-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25z"/></svg>`
			} else {
				typeIcon = `<svg class="phase-type-icon" viewBox="0 0 16 16" fill="currentColor"><path d="M3.75 1.5a.25.25 0 00-.25.25v12.5c0 .138.112.25.25.25h8.5a.25.25 0 00.25-.25V4.664a.25.25 0 00-.073-.177l-2.914-2.914a.25.25 0 00-.177-.073H3.75zM2 1.75C2 .784 2.784 0 3.75 0h5.586c.464 0 .909.184 1.237.513l2.914 2.914c.329.328.513.773.513 1.237v9.586A1.75 1.75 0 0112.25 16h-8.5A1.75 1.75 0 012 14.25V1.75z"/></svg>`
			}

			fmt.Fprintf(&itemsHTML, `
        <li class="phase-item%s%s" data-status="%s"%s>
          <a href="%s">
            %s
            <span class="status-dot %s"></span>
            <span class="phase-name">%s</span>
          </a>
        </li>`, activeClass, inlineSectionClass, statusClass, dataAnchor, href, typeIcon, statusClass, html.EscapeString(phase.Name))
		}

		fmt.Fprintf(&groupsHTML, `
      <div class="phase-group" data-phase-id="%s">
        <button class="phase-header" tabindex="0" aria-expanded="true" aria-controls="%s-items">
          <span class="phase-chevron">▼</span>
          <span class="phase-name">%s</span>
          %s
        </button>
        <ul class="phase-items" id="%s-items">%s
        </ul>
      </div>`, groupID, groupID, groupLabel, groupBadge, groupID, itemsHTML.String())
	}

	return fmt.Sprintf(`
    <nav class="plan-nav" id="plan-nav">
      <div class="plan-title">
        <span class="plan-icon">&#128214;</span>
        <span>%s</span>
      </div>%s
    </nav>`, html.EscapeString(planName), groupsHTML.String())
}

func getGroupBadge(phases []indexedPhase) string {
	completed := 0
	inProgress := 0
	for _, ip := range phases {
		s := ip.Phase.Status
		if s == "completed" || s == "done" {
			completed++
		} else if s == "in-progress" {
			inProgress++
		}
	}
	if completed == len(phases) {
		return `<span class="phase-badge badge-done">✓</span>`
	}
	if inProgress > 0 {
		return `<span class="phase-badge badge-progress">●</span>`
	}
	return `<span class="phase-badge badge-pending">○</span>`
}
