package process

import (
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

func getPidFilePath(port int) string {
	return filepath.Join(pidDir, fmt.Sprintf("%s%d.pid", pidPrefix, port))
}

func WritePidFile(port, pid int) {
	os.WriteFile(getPidFilePath(port), []byte(strconv.Itoa(pid)), 0644)
}

func ReadPidFile(port int) int {
	data, err := os.ReadFile(getPidFilePath(port))
	if err != nil {
		return 0
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return pid
}

func RemovePidFile(port int) {
	os.Remove(getPidFilePath(port))
}

func FindRunningInstances() []struct{ Port, Pid int } {
	var instances []struct{ Port, Pid int }

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
		pid := ReadPidFile(port)
		if pid == 0 {
			continue
		}
		// Check if process is running
		proc, err := os.FindProcess(pid)
		if err != nil {
			RemovePidFile(port)
			continue
		}
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			RemovePidFile(port)
			continue
		}
		instances = append(instances, struct{ Port, Pid int }{port, pid})
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
