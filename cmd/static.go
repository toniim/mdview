package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bilabl/mdview/internal/process"
	"github.com/spf13/cobra"
)

var staticCmd = &cobra.Command{
	Use:   "static [path]",
	Short: "Serve a directory as a plain static file server",
	Long:  "Serve a directory using Go's built-in file server.\nNo markdown rendering — files are served as-is, with native directory listing.\nDefaults to the current directory if no path is given.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runStatic,
}

func init() {
	staticCmd.Flags().IntVarP(&port, "port", "p", 3456, "Server port (auto-increment if busy)")
	staticCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "Bind address")
	staticCmd.Flags().BoolVar(&public, "public", false, "Bind to 0.0.0.0 (shortcut for --host 0.0.0.0)")
	staticCmd.Flags().BoolVarP(&open, "open", "o", true, "Open browser")
	staticCmd.Flags().BoolVar(&noOpen, "no-open", false, "Don't open browser")
	staticCmd.Flags().BoolVar(&foreground, "foreground", false, "Foreground mode (JSON output for CC)")
}

func runStatic(cmd *cobra.Command, args []string) error {
	if noOpen {
		open = false
	}
	if public {
		host = "0.0.0.0"
	}

	inputPath := "."
	if len(args) == 1 {
		inputPath = args[0]
	}

	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("invalid path: %s", inputPath)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("path not found: %s", inputPath)
	}
	if !info.IsDir() {
		return fmt.Errorf("static mode requires a directory, got file: %s", inputPath)
	}

	actualPort, err := process.FindAvailablePort(port)
	if err != nil {
		return err
	}
	if actualPort != port {
		fmt.Fprintf(os.Stderr, "Port %d in use, using %d\n", port, actualPort)
	}

	displayHost := host
	if host == "0.0.0.0" {
		displayHost = "localhost"
	}
	url := fmt.Sprintf("http://%s:%d/", displayHost, actualPort)

	var networkURL string
	if ip := getLocalIP(); ip != "" {
		networkURL = fmt.Sprintf("http://%s:%d/", ip, actualPort)
	}

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, actualPort),
		Handler:      http.FileServer(http.Dir(absPath)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	process.WritePidFile(process.InstanceInfo{
		Port: actualPort,
		Pid:  os.Getpid(),
		Host: host,
		Path: absPath,
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		process.RemovePidFile(actualPort)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
		os.Exit(0)
	}()

	if foreground {
		result := serveResult{
			Success:    true,
			URL:        url,
			Path:       absPath,
			Port:       actualPort,
			Host:       host,
			Mode:       "static",
			NetworkURL: networkURL,
		}
		json.NewEncoder(os.Stdout).Encode(result)
	} else {
		fmt.Println("\nMarkdown Novel Viewer (Go) — static mode")
		fmt.Println("────────────────────────────────────────")
		fmt.Printf("URL:  %s\n", url)
		if networkURL != "" {
			fmt.Printf("Network: %s\n", networkURL)
		}
		fmt.Printf("Path: %s\n", absPath)
		fmt.Printf("Port: %d\n", actualPort)
		fmt.Printf("Host: %s\n", host)
		fmt.Println("Mode: static")
		fmt.Println("\nPress Ctrl+C to stop")
	}

	if open {
		openBrowser(url)
	}

	return httpServer.ListenAndServe()
}
