package navigation

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bilabl/mdview/internal/renderer"
)

// PhaseEntry represents a parsed phase from a plan table.
type PhaseEntry struct {
	Phase    int
	Name     string
	Status   string
	File     string
	LinkText string
	Anchor   string
}

func normalizeStatus(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if strings.Contains(s, "complete") || strings.Contains(s, "done") || strings.Contains(s, "✓") || strings.Contains(s, "✅") {
		return "completed"
	}
	if strings.Contains(s, "progress") || strings.Contains(s, "active") || strings.Contains(s, "wip") || strings.Contains(s, "🔄") {
		return "in-progress"
	}
	return "pending"
}

var (
	// Format 1: Standard table | Phase | Name | Status | [Link](path) |
	standardRegex = regexp.MustCompile(`\|\s*(\d+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|\s*\[([^\]]+)\]\(([^)]+)\)`)
	// Format 2: Link-first | [Phase X](path) | Description | Status |
	linkFirstRegex = regexp.MustCompile(`\|\s*\[(?:Phase\s*)?(\d+)\]\(([^)]+)\)\s*\|\s*([^|]+)\s*\|\s*([^|]+)`)
	// Format 2b: Number-first with link | 1 | [Name](path) | Status |
	numLinkRegex = regexp.MustCompile(`\|\s*(\d+)\s*\|\s*\[([^\]]+)\]\(([^)]+)\)\s*\|\s*([^|]+)`)
	// Format 2c: Simple table | Phase | Description | Status |
	simpleTblRegex = regexp.MustCompile(`\|\s*0?(\d+)\s*\|\s*([^|]+)\s*\|\s*([^|]+)\s*\|`)
	// Format 3: Heading-based ### Phase X: Name
	headingPhaseRegex = regexp.MustCompile(`(?i)###\s*Phase\s*(\d+)[:\s]+(.+)`)
	statusLineRegex   = regexp.MustCompile(`(?i)-\s*Status:\s*(.+)`)
	// Format 4: Bullet-list phases
	bulletPhaseRegex = regexp.MustCompile(`(?i)^-\s*Phase\s*0?(\d+)[:\s]+([^✅✓\n]+)`)
	bulletFileRegex  = regexp.MustCompile(`(?i)^\s+-\s*File:\s*` + "`?" + `([^` + "`" + `\n]+)` + "`?")
	bulletStatusLine = regexp.MustCompile(`(?i)^\s+-\s*(Completed|Status):\s*(.+)`)
	// Format 5: Numbered list
	numberedPhaseRegex = regexp.MustCompile(`(?m)^(\d+)[).]\s*\*\*([^*]+)\*\*`)
	checkboxRegex      = regexp.MustCompile(`(?mi)^-\s*\[(x| )\]\s*([^:]+)`)
	// Format 6: Checkbox list with bold links
	checkboxLinkRegex = regexp.MustCompile(`(?mi)^-\s*\[(x| )\]\s*\*\*\[(?:Phase\s*)?(\d+)[:\s]*([^\]]*)\]\(([^)]+)\)\*\*`)
	// Phase files section
	phaseFilesRegex = regexp.MustCompile(`(?is)##\s*Phase\s*Files[\s\S]*?(?:##|$)`)
	phaseFileLink   = regexp.MustCompile(`\d+\.\s*\[([^\]]+)\]\(([^)]+\.md)\)`)
	phaseNumInName  = regexp.MustCompile(`(?i)phase-0?(\d+)`)
	// Bullet phase pattern check
	bulletPhaseCheck = regexp.MustCompile(`(?m)^-\s*Phase\s*\d+[:\s]`)
)

// ParsePlanTable parses plan.md to extract phase metadata from various table formats.
func ParsePlanTable(planFilePath string) []PhaseEntry {
	data, err := os.ReadFile(planFilePath)
	if err != nil {
		return nil
	}
	content := string(data)
	dir := filepath.Dir(planFilePath)

	var phases []PhaseEntry

	// Format 1: Standard table
	for _, m := range standardRegex.FindAllStringSubmatch(content, -1) {
		num, _ := strconv.Atoi(m[1])
		phases = append(phases, PhaseEntry{
			Phase:    num,
			Name:     strings.TrimSpace(m[2]),
			Status:   normalizeStatus(m[3]),
			File:     filepath.Join(dir, m[5]),
			LinkText: strings.TrimSpace(m[4]),
		})
	}

	// Format 2: Link-first
	if len(phases) == 0 {
		for _, m := range linkFirstRegex.FindAllStringSubmatch(content, -1) {
			num, _ := strconv.Atoi(m[1])
			phases = append(phases, PhaseEntry{
				Phase:    num,
				Name:     strings.TrimSpace(m[3]),
				Status:   normalizeStatus(m[4]),
				File:     filepath.Join(dir, m[2]),
				LinkText: "Phase " + m[1],
			})
		}
	}

	// Format 2b: Number-first with link
	if len(phases) == 0 {
		for _, m := range numLinkRegex.FindAllStringSubmatch(content, -1) {
			num, _ := strconv.Atoi(m[1])
			phases = append(phases, PhaseEntry{
				Phase:    num,
				Name:     strings.TrimSpace(m[2]),
				Status:   normalizeStatus(m[4]),
				File:     filepath.Join(dir, m[3]),
				LinkText: strings.TrimSpace(m[2]),
			})
		}
	}

	// Format 2c: Simple table without links
	if len(phases) == 0 {
		for _, m := range simpleTblRegex.FindAllStringSubmatch(content, -1) {
			name := strings.TrimSpace(m[2])
			nameLower := strings.ToLower(name)
			if nameLower == "description" || nameLower == "name" ||
				strings.Contains(name, "---") || strings.Contains(name, "===") {
				continue
			}
			num, _ := strconv.Atoi(m[1])
			phases = append(phases, PhaseEntry{
				Phase:    num,
				Name:     name,
				Status:   normalizeStatus(m[3]),
				File:     planFilePath,
				LinkText: name,
				Anchor:   padPhaseAnchor(num, name),
			})
		}
	}

	// Format 3: Heading-based
	if len(phases) == 0 {
		lines := strings.Split(content, "\n")
		var current *PhaseEntry
		for _, line := range lines {
			if hm := headingPhaseRegex.FindStringSubmatch(line); hm != nil {
				if current != nil {
					phases = append(phases, *current)
				}
				num, _ := strconv.Atoi(hm[1])
				name := strings.TrimSpace(hm[2])
				current = &PhaseEntry{
					Phase:    num,
					Name:     name,
					Status:   "pending",
					File:     planFilePath,
					LinkText: "Phase " + hm[1],
					Anchor:   padPhaseAnchor(num, name),
				}
			}
			if current != nil {
				if sm := statusLineRegex.FindStringSubmatch(line); sm != nil {
					current.Status = normalizeStatus(sm[1])
				}
			}
		}
		if current != nil {
			phases = append(phases, *current)
		}
	}

	// Format 4: Bullet-list phases
	if len(phases) == 0 && bulletPhaseCheck.MatchString(content) {
		lines := strings.Split(content, "\n")
		var current *PhaseEntry
		for _, line := range lines {
			if bm := bulletPhaseRegex.FindStringSubmatch(line); bm != nil {
				if current != nil {
					phases = append(phases, *current)
				}
				num, _ := strconv.Atoi(bm[1])
				name := strings.TrimSpace(bm[2])
				name = regexp.MustCompile(`\s*\([^)]*\)\s*$`).ReplaceAllString(name, "")
				hasCheck := strings.Contains(line, "✅") || strings.Contains(line, "✓")
				status := "pending"
				if hasCheck {
					status = "completed"
				}
				current = &PhaseEntry{
					Phase:    num,
					Name:     name,
					Status:   status,
					File:     planFilePath,
					LinkText: name,
					Anchor:   padPhaseAnchor(num, name),
				}
				continue
			}
			if current != nil {
				if fm := bulletFileRegex.FindStringSubmatch(line); fm != nil {
					current.File = filepath.Join(dir, strings.TrimSpace(fm[1]))
					current.Anchor = ""
				}
				if sm := bulletStatusLine.FindStringSubmatch(line); sm != nil {
					current.Status = normalizeStatus(sm[2])
				}
				if strings.HasPrefix(line, "##") ||
					(strings.HasPrefix(line, "- ") && !strings.HasPrefix(strings.TrimSpace(line), "- Phase") && !strings.HasPrefix(line, "  ")) {
					phases = append(phases, *current)
					current = nil
				}
			}
		}
		if current != nil {
			phases = append(phases, *current)
		}
	}

	// Format 5: Numbered list with checkbox status
	if len(phases) == 0 {
		phaseMap := make(map[string]*PhaseEntry)
		var order []string
		for _, m := range numberedPhaseRegex.FindAllStringSubmatch(content, -1) {
			num, _ := strconv.Atoi(m[1])
			name := strings.TrimSpace(m[2])
			key := strings.ToLower(name)
			entry := &PhaseEntry{
				Phase:    num,
				Name:     name,
				Status:   "pending",
				File:     planFilePath,
				LinkText: name,
				Anchor:   padPhaseAnchor(num, name),
			}
			phaseMap[key] = entry
			order = append(order, key)
		}
		for _, m := range checkboxRegex.FindAllStringSubmatch(content, -1) {
			key := strings.ToLower(strings.TrimSpace(m[2]))
			if e, ok := phaseMap[key]; ok {
				if strings.ToLower(m[1]) == "x" {
					e.Status = "completed"
				}
			}
		}
		if len(phaseMap) > 0 {
			for _, key := range order {
				phases = append(phases, *phaseMap[key])
			}
			sort.Slice(phases, func(i, j int) bool { return phases[i].Phase < phases[j].Phase })
		}
	}

	// Format 6: Checkbox list with bold links
	if len(phases) == 0 {
		for _, m := range checkboxLinkRegex.FindAllStringSubmatch(content, -1) {
			num, _ := strconv.Atoi(m[2])
			name := strings.TrimSpace(m[3])
			if name == "" {
				name = "Phase " + m[2]
			}
			status := "pending"
			if strings.ToLower(m[1]) == "x" {
				status = "completed"
			}
			phases = append(phases, PhaseEntry{
				Phase:    num,
				Name:     name,
				Status:   status,
				File:     filepath.Join(dir, m[4]),
				LinkText: name,
			})
		}
	}

	// Enhancement: extract file paths from "Phase Files" section
	if len(phases) > 0 {
		if section := phaseFilesRegex.FindString(content); section != "" {
			for _, m := range phaseFileLink.FindAllStringSubmatch(section, -1) {
				linkName := m[1]
				linkPath := m[2]
				if nm := phaseNumInName.FindStringSubmatch(linkName); nm != nil {
					num, _ := strconv.Atoi(nm[1])
					for i := range phases {
						if phases[i].Phase == num && phases[i].File == planFilePath {
							phases[i].File = filepath.Join(dir, linkPath)
						}
					}
				}
			}
		}
	}

	// Filter out phases pointing to plan.md itself
	var filtered []PhaseEntry
	for _, p := range phases {
		abs1, _ := filepath.Abs(p.File)
		abs2, _ := filepath.Abs(planFilePath)
		if abs1 != abs2 {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func padPhaseAnchor(num int, name string) string {
	return "phase-" + padNum(num) + "-" + renderer.Slugify(name)
}

func padNum(n int) string {
	s := strconv.Itoa(n)
	if len(s) < 2 {
		return "0" + s
	}
	return s
}
