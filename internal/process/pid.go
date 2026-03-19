package process

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	pidDir    = "/tmp"
	pidPrefix = "md-novel-viewer-"
)

// InstanceInfo holds metadata about a running mdview server.
type InstanceInfo struct {
	Port int    `json:"port"`
	Pid  int    `json:"pid"`
	Host string `json:"host,omitempty"`
	Path string `json:"path,omitempty"`
}

func getPidFilePath(port int) string {
	return filepath.Join(pidDir, fmt.Sprintf("%s%d.pid", pidPrefix, port))
}

// WritePidFile stores instance metadata as JSON.
func WritePidFile(info InstanceInfo) {
	data, _ := json.Marshal(info)
	os.WriteFile(getPidFilePath(info.Port), data, 0644)
}

// readInstanceInfo reads PID file and returns info.
// Handles both JSON (new) and plain-PID (old) formats.
func readInstanceInfo(port int) (InstanceInfo, bool) {
	data, err := os.ReadFile(getPidFilePath(port))
	if err != nil {
		return InstanceInfo{}, false
	}
	content := strings.TrimSpace(string(data))
	if content == "" {
		return InstanceInfo{}, false
	}

	// Try JSON first
	var info InstanceInfo
	if err := json.Unmarshal([]byte(content), &info); err == nil && info.Pid > 0 {
		info.Port = port
		return info, true
	}

	// Fallback: plain PID
	pid, err := strconv.Atoi(content)
	if err != nil || pid <= 0 {
		return InstanceInfo{}, false
	}
	return InstanceInfo{Port: port, Pid: pid}, true
}

func RemovePidFile(port int) {
	os.Remove(getPidFilePath(port))
}

// FindRunningInstances returns info about all live mdview servers.
func FindRunningInstances() []InstanceInfo {
	var instances []InstanceInfo

	entries, err := os.ReadDir(pidDir)
	if err != nil {
		return instances
	}

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, pidPrefix) || !strings.HasSuffix(name, ".pid") {
			continue
		}
		portStr := strings.TrimPrefix(strings.TrimSuffix(name, ".pid"), pidPrefix)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}
		info, ok := readInstanceInfo(port)
		if !ok {
			continue
		}
		// Check if process is running
		proc, err := os.FindProcess(info.Pid)
		if err != nil {
			RemovePidFile(port)
			continue
		}
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			RemovePidFile(port)
			continue
		}
		instances = append(instances, info)
	}
	return instances
}

func StopAllServers() int {
	instances := FindRunningInstances()
	stopped := 0
	for _, inst := range instances {
		proc, err := os.FindProcess(inst.Pid)
		if err != nil {
			RemovePidFile(inst.Port)
			continue
		}
		if err := proc.Signal(syscall.SIGTERM); err == nil {
			stopped++
		}
		RemovePidFile(inst.Port)
	}
	return stopped
}
