package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"FsHost/theme"
)

var (
	version = "1.0.0"
	banner  = `
  ______   _    _           _   
 |  ____| | |  | |         | |  
 | |__ ___| |__| | ___  ___| |_ 
 |  __/ __|  __  |/ _ \/ __| __|
 | |  \__ \ |  | | (_) \__ \ |_ 
 |_|  |___/_|  |_|\___/|___/\__|
                                
    `

	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorPurple = "\033[35m"
)

func main() {
	port := flag.Int("port", 8080, "Port to serve on")
	dir := flag.String("dir", ".", "Directory to share")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("%s✖ Error resolving directory:%s %v\n", colorRed, colorReset, err)
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		log.Fatalf("%s✖ Directory does not exist:%s %s\n", colorRed, colorReset, absDir)
	}

	ip := getLocalIP()

	printBanner()

	fmt.Printf("  %s📂 Serving:%s  %s%s%s\n", colorCyan, colorReset, colorBold, absDir, colorReset)
	fmt.Printf("  %s🌐 Network:%s  %shttp://%s:%d%s\n", colorCyan, colorReset, colorGreen, ip, *port, colorReset)
	fmt.Printf("  %s💻 Local:%s    %shttp://localhost:%d%s\n", colorCyan, colorReset, colorGreen, *port, colorReset)
	fmt.Printf("  %s⏱  Started:%s  %s%s%s\n", colorCyan, colorReset, colorPurple, time.Now().Format("2006-01-02 15:04:05"), colorReset)
	fmt.Printf("  %s🛡  OS:%s       %s%s%s\n", colorCyan, colorReset, colorDim, runtime.GOOS, colorReset)
	fmt.Println()
	fmt.Printf("  %s%sPress Ctrl+C to stop serving%s\n", colorDim, colorYellow, colorReset)
	fmt.Println()

	handler := newFileServer(absDir)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Printf("\n\n  %s🛑 Shutting down GoDir...%s\n", colorRed, colorReset)
		server.Close()
	}()

	fmt.Printf("  %s✔ Server is live! Open the link above in any browser on your network.%s\n\n", colorGreen, colorReset)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("%s✖ Server error:%s %v\n", colorRed, colorReset, err)
	}

	fmt.Printf("  %s👋 GoDir stopped. Goodbye!%s\n\n", colorGreen, colorReset)
}

func printBanner() {
	lines := strings.Split(banner, "\n")
	cyan := "\033[36m"
	bold := "\033[1m"
	reset := "\033[0m"
	dim := "\033[2m"
	yellow := "\033[33m"

	fmt.Println()
	for i, line := range lines {
		if i < len(lines)-1 {
			fmt.Printf("  %s%s%s%s\n", cyan, bold, line, reset)
		} else {
			fmt.Printf("  %s%s%s%s\n", dim, yellow, line, reset)
		}
	}
	fmt.Println()
	fmt.Printf("  %s%s v%s %s— Fast, beautiful local file sharing%s\n\n", bold, cyan, version, dim, reset)
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "localhost"
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if strings.Contains(iface.Name, "lo") || strings.Contains(iface.Name, "utun") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				octets := strings.Split(ip4.String(), ".")
				if len(octets) == 4 && octets[0] == "192" {
					return ip4.String()
				}
			}
		}
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				return ip4.String()
			}
		}
	}
	return "localhost"
}

type fileInfo struct {
	Name    string
	Size    string
	IsDir   bool
	ModTime string
	Icon    string
	Ext     string
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func getFileIcon(name string, isDir bool) string {
	if isDir {
		return "\U0001F4C1"
	}
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", ".ico":
		return "\U0001F5BC"
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm":
		return "\U0001F3AC"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".wma", ".m4a":
		return "\U0001F3B5"
	case ".pdf":
		return "\U0001F4C4"
	case ".doc", ".docx":
		return "\U0001F4DD"
	case ".xls", ".xlsx", ".csv":
		return "\U0001F4CA"
	case ".ppt", ".pptx":
		return "\U0001F4D1"
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2":
		return "\U0001F4E6"
	case ".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".rs", ".rb", ".php", ".sh":
		return "\U0001F4BB"
	case ".txt", ".md", ".log", ".ini", ".cfg", ".conf":
		return "\U0001F4DC"
	case ".html", ".css", ".xml", ".json", ".yaml", ".yml":
		return "\U0001F310"
	case ".exe", ".msi", ".dmg", ".app":
		return "\u2699"
	case ".ttf", ".otf", ".woff", ".woff2":
		return "\U0001F524"
	default:
		return "\U0001F4C4"
	}
}

type breadcrumbPart struct {
	Name string
	Path string
}

type pageData struct {
	Title     string
	Path      string
	Bread     []breadcrumbPart
	Files     []fileInfo
	ParentDir string
	IPAddress string
	Version   string
	NumFiles  int
}

func newFileServer(rootDir string) http.Handler {
	mux := http.NewServeMux()

	funcMap := template.FuncMap{
		"fileCount": func(f []fileInfo) int { return len(f) },
	}

	tmpl := template.Must(template.New("dir").Funcs(funcMap).Parse(theme.PageTemplate))

	fileServer := http.FileServer(http.Dir(rootDir))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cleanPath := filepath.Clean(r.URL.Path)
		if cleanPath == "/" || cleanPath == "." {
			cleanPath = "/"
		}

		fullPath := filepath.Join(rootDir, cleanPath)
		info, err := os.Stat(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		entries, err := os.ReadDir(fullPath)
		if err != nil {
			http.Error(w, "Cannot read directory", http.StatusInternalServerError)
			return
		}

		files := make([]fileInfo, 0, len(entries))

		if cleanPath != "/" {
			files = append(files, fileInfo{
				Name:    "..",
				Size:    "-",
				IsDir:   true,
				ModTime: "-",
				Icon:    "\u2B06",
				Ext:     "",
			})
		}

		dirs := make([]fileInfo, 0)
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			entryInfo, err := entry.Info()
			if err != nil {
				continue
			}
			f := fileInfo{
				Name:    name,
				Size:    formatSize(entryInfo.Size()),
				IsDir:   entry.IsDir(),
				ModTime: entryInfo.ModTime().Format("2006-01-02 15:04"),
				Icon:    getFileIcon(name, entry.IsDir()),
				Ext:     strings.ToLower(filepath.Ext(name)),
			}
			if entry.IsDir() {
				dirs = append(dirs, f)
			} else {
				files = append(files, f)
			}
			if len(dirs)+len(files)-1 > 5000 {
				break
			}
		}
		files = append(dirs, files...)

		trimmedPath := strings.Trim(cleanPath, "/")
		breadParts := make([]breadcrumbPart, 0)
		if trimmedPath != "" {
			parts := strings.Split(trimmedPath, "/")
			accumPath := ""
			for _, part := range parts {
				accumPath += "/" + part + "/"
				breadParts = append(breadParts, breadcrumbPart{
					Name: part,
					Path: accumPath,
				})
			}
		}

		ip := getLocalIP()

		data := pageData{
			Title:     "GoDir \u2014 " + cleanPath,
			Path:      cleanPath,
			Bread:     breadParts,
			Files:     files,
			ParentDir: filepath.Dir(strings.TrimSuffix(cleanPath, "/")),
			IPAddress: ip,
			Version:   version,
			NumFiles:  len(files),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, data)
	})

	return logMiddleware(mux)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		status := "\033[32m"
		if wrapped.statusCode >= 400 {
			status = "\033[33m"
		}
		if wrapped.statusCode >= 500 {
			status = "\033[31m"
		}

		fmt.Printf("  %s%s%d%s %s %s%s%s\n",
			"\033[2m", status, wrapped.statusCode, "\033[0m",
			r.URL.Path,
			"\033[2m", duration.Round(time.Millisecond), "\033[0m")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
