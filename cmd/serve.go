package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/bilabl/mdview/internal/process"
	"github.com/bilabl/mdview/internal/server"
	"github.com/spf13/cobra"
)

var (
	port       int
	host       string
	open       bool
	noOpen     bool
	foreground bool
)

var serveCmd = &cobra.Command{
	Use:   "serve [path]",
	Short: "Serve a file or directory",
	Long:  "Serve a markdown file, code file, or directory for viewing in the browser.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 3456, "Server port (auto-increment if busy)")
	serveCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Bind address (0.0.0.0 for remote)")
	serveCmd.Flags().BoolVarP(&open, "open", "o", true, "Open browser")
	serveCmd.Flags().BoolVar(&noOpen, "no-open", false, "Don't open browser")
	serveCmd.Flags().BoolVar(&foreground, "foreground", false, "Foreground mode (JSON output for CC)")
}

type serveResult struct {
	Success    bool   `json:"success"`
	URL        string `json:"url"`
	Path       string `json:"path"`
	Port       int    `json:"port"`
	Host       string `json:"host"`
	Mode       string `json:"mode"`
	NetworkURL string `json:"networkUrl,omitempty"`
}

func runServe(cmd *cobra.Command, args []string) error {
	if noOpen {
		open = false
	}

	if len(args) == 0 {
		return fmt.Errorf("path argument required\nUsage: mdview serve <path> [--port 3456] [--no-open]")
	}

	inputPath := args[0]

	// Resolve path
	resolved, err := resolvePath(inputPath)
	if err != nil {
		return err
	}

	// Find available port
	actualPort, err := process.FindAvailablePort(port)
	if err != nil {
		return err
	}
	if actualPort != port {
		fmt.Fprintf(os.Stderr, "Port %d in use, using %d\n", port, actualPort)
	}

	// Build URL
	displayHost := host
	if host == "0.0.0.0" {
		displayHost = "localhost"
	}

	var urlPath string
	switch resolved.Type {
	case "file":
		urlPath = "/view?file=" + netURLEncode(resolved.Path)
	case "directory":
		urlPath = "/browse?dir=" + netURLEncode(resolved.Path)
	}

	url := fmt.Sprintf("http://%s:%d%s", displayHost, actualPort, urlPath)

	var networkURL string
	if ip := getLocalIP(); ip != "" {
		networkURL = fmt.Sprintf("http://%s:%d%s", ip, actualPort, urlPath)
	}

	// Build allowed directories
	allowedDirs := []string{filepath.Dir(resolved.Path)}
	if resolved.Type == "directory" {
		allowedDirs = []string{resolved.Path}
	}

	// Create and start server
	srv := server.New(server.Config{
		Host:        host,
		Port:        actualPort,
		AllowedDirs: allowedDirs,
	})

	// Write PID file
	process.WritePidFile(actualPort, os.Getpid())

	// Setup shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		process.RemovePidFile(actualPort)
		srv.Shutdown()
		os.Exit(0)
	}()

	// Output
	if foreground {
		result := serveResult{
			Success:    true,
			URL:        url,
			Path:       resolved.Path,
			Port:       actualPort,
			Host:       host,
			Mode:       resolved.Type,
			NetworkURL: networkURL,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(result)
	} else {
		fmt.Println("\nMarkdown Novel Viewer (Go)")
		fmt.Println("────────────────────────────────────────")
		fmt.Printf("URL:  %s\n", url)
		if networkURL != "" {
			fmt.Printf("Network: %s\n", networkURL)
		}
		fmt.Printf("Path: %s\n", resolved.Path)
		fmt.Printf("Port: %d\n", actualPort)
		fmt.Printf("Host: %s\n", host)
		fmt.Printf("Mode: %s\n", resolved.Type)
		fmt.Println("\nPress Ctrl+C to stop")
	}

	// Open browser
	if open {
		openBrowser(url)
	}

	// Start server (blocks)
	return srv.ListenAndServe()
}

type resolvedPath struct {
	Type string // "file" or "directory"
	Path string
}

func resolvePath(input string) (*resolvedPath, error) {
	abs, err := filepath.Abs(input)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", input)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("path not found: %s", input)
	}

	if info.IsDir() {
		return &resolvedPath{Type: "directory", Path: abs}, nil
	}
	return &resolvedPath{Type: "file", Path: abs}, nil
}

func netURLEncode(s string) string {
	// Simple URL encoding for path parameters
	var result []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEncode(c) {
			result = append(result, '%')
			result = append(result, "0123456789ABCDEF"[c>>4])
			result = append(result, "0123456789ABCDEF"[c&0x0f])
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

func shouldEncode(c byte) bool {
	// RFC 3986 unreserved characters
	if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '~', '/':
		return false
	}
	return true
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}
