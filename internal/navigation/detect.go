package navigation

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// PlanInfo holds plan detection results.
type PlanInfo struct {
	IsPlan  bool
	PlanDir string
	PlanFile string
	Phases  []string
}

var phaseNumRegex = regexp.MustCompile(`phase-(\d+)`)

// DetectPlan checks if a file is part of a plan directory.
func DetectPlan(filePath string) *PlanInfo {
	dir := filepath.Dir(filePath)
	planFile := filepath.Join(dir, "plan.md")

	if _, err := os.Stat(planFile); os.IsNotExist(err) {
		return &PlanInfo{}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return &PlanInfo{}
	}

	var phases []string
	for _, e := range entries {
		name := e.Name()
		if len(name) > 6 && name[:6] == "phase-" && filepath.Ext(name) == ".md" {
			phases = append(phases, filepath.Join(dir, name))
		}
	}

	sort.Slice(phases, func(i, j int) bool {
		numI := extractPhaseNum(filepath.Base(phases[i]))
		numJ := extractPhaseNum(filepath.Base(phases[j]))
		return numI < numJ
	})

	return &PlanInfo{
		IsPlan:   true,
		PlanDir:  dir,
		PlanFile: planFile,
		Phases:   phases,
	}
}

func extractPhaseNum(filename string) int {
	m := phaseNumRegex.FindStringSubmatch(filename)
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}
