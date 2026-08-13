package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"sort"
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
	port := flag.Int("port", 8080, "Port to serve on (1-65535)")
	dir := flag.String("dir", ".", "Directory to share")
	flag.Parse()

	if *port < 1 || *port > 65535 {
		log.Fatalf("%s✖ Invalid port:%s %d (must be between 1 and 65535)\n", colorRed, colorReset, *port)
	}

	absDir, err := filepath.Abs(expandHome(*dir))
	if err != nil {
		log.Fatalf("%s✖ Error resolving directory:%s %v\n", colorRed, colorReset, err)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		log.Fatalf("%s✖ Directory does not exist:%s %s\n", colorRed, colorReset, absDir)
	}
	if !info.IsDir() {
		log.Fatalf("%s✖ Not a directory:%s %s\n", colorRed, colorReset, absDir)
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
	fmt.Printf("  %s✔ Server is live! Open the link above in any browser on your network.%s\n\n", colorGreen, colorReset)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      newFileServer(absDir),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Printf("\n\n  %s🛑 Shutting down FsHost...%s\n", colorRed, colorReset)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("%s⚠ Graceful shutdown failed:%s %v\n", colorYellow, colorReset, err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("%s✖ Server error:%s %v\n", colorRed, colorReset, err)
	}

	fmt.Printf("  %s👋 FsHost stopped. Goodbye!%s\n\n", colorGreen, colorReset)
}

func expandHome(p string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	switch {
	case p == "~":
		return home
	case strings.HasPrefix(p, "~/"):
		return filepath.Join(home, strings.TrimPrefix(p, "~/"))
	default:
		return p
	}
}

func printBanner() {
	lines := strings.Split(banner, "\n")
	fmt.Println()
	for i, line := range lines {
		if i < len(lines)-1 {
			fmt.Printf("  %s%s%s%s\n", colorCyan, colorBold, line, colorReset)
		} else {
			fmt.Printf("  %s%s%s%s\n", colorDim, colorYellow, line, colorReset)
		}
	}
	fmt.Println()
	fmt.Printf("  %s%sFsHost v%s %s— Fast, beautiful local file sharing%s\n\n", colorBold, colorCyan, version, colorDim, colorReset)
}

func getLocalIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "localhost"
	}

	var candidates []net.IP
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
				candidates = append(candidates, ip4)
			}
		}
	}

	for _, ip := range candidates {
		if strings.HasPrefix(ip.String(), "192.") {
			return ip.String()
		}
	}
	if len(candidates) > 0 {
		return candidates[0].String()
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
	Href    string
}

type breadcrumbPart struct {
	Name string
	Path string
}

type pageData struct {
	Title     string
	Bread     []breadcrumbPart
	Files     []fileInfo
	IPAddress string
	Version   string
	NumFiles  int
}

func newFileServer(rootDir string) http.Handler {
	tmpl := template.Must(template.New("dir").Parse(theme.PageTemplate))
	fileServer := http.FileServer(http.Dir(rootDir))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fullPath, ok := safeJoin(rootDir, r.URL.Path)
		if !ok {
			http.NotFound(w, r)
			return
		}

		info, err := os.Stat(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, escapeURLPath(r.URL.Path)+"/", http.StatusMovedPermanently)
			return
		}

		entries, err := os.ReadDir(fullPath)
		if err != nil {
			http.Error(w, "Cannot read directory", http.StatusInternalServerError)
			return
		}

		base := escapeURLPath(strings.TrimSuffix(r.URL.Path, "/")) + "/"

		listing := make([]fileInfo, 0, len(entries)+1)
		if trimmed := strings.TrimSuffix(r.URL.Path, "/"); trimmed != "" {
			listing = append(listing, fileInfo{
				Name:    "..",
				Size:    "-",
				IsDir:   true,
				ModTime: "-",
				Icon:    "\u2B06",
				Ext:     "",
				Href:    escapeURLPath(path.Join(trimmed, "..")),
			})
		}

		var dirs, files []fileInfo
		count := 0
		for _, entry := range entries {
			if count >= 5000 {
				break
			}
			name := entry.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			entryInfo, err := entry.Info()
			if err != nil {
				continue
			}
			isDir := entry.IsDir()
			href := base + url.PathEscape(name)
			if isDir {
				href += "/"
			}
			f := fileInfo{
				Name:    name,
				Size:    formatSize(entryInfo.Size()),
				IsDir:   isDir,
				ModTime: entryInfo.ModTime().Format("2006-01-02 15:04"),
				Icon:    getFileIcon(name, isDir),
				Ext:     strings.ToLower(filepath.Ext(name)),
				Href:    href,
			}
			if isDir {
				dirs = append(dirs, f)
			} else {
				files = append(files, f)
			}
			count++
		}

		sort.Slice(dirs, func(i, j int) bool { return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name) })
		sort.Slice(files, func(i, j int) bool { return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name) })

		listing = append(listing, dirs...)
		listing = append(listing, files...)

		data := pageData{
			Title:     "FsHost \u2014 " + r.URL.Path,
			Bread:     buildBreadcrumbs(r.URL.Path),
			Files:     listing,
			IPAddress: getLocalIP(),
			Version:   version,
			NumFiles:  len(listing),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("%s✖ Template error:%s %v\n", colorRed, colorReset, err)
		}
	})

	return logMiddleware(mux)
}

func safeJoin(rootDir, urlPath string) (string, bool) {
	clean := path.Clean(urlPath)
	if clean == "." {
		clean = "/"
	}
	full := filepath.Join(rootDir, filepath.FromSlash(clean))
	rel, err := filepath.Rel(rootDir, full)
	if err != nil {
		return "", false
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return full, true
}

func escapeURLPath(p string) string {
	segments := strings.Split(p, "/")
	for i, s := range segments {
		segments[i] = url.PathEscape(s)
	}
	return strings.Join(segments, "/")
}

func buildBreadcrumbs(urlPath string) []breadcrumbPart {
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")
	out := make([]breadcrumbPart, 0, len(parts))
	accum := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		accum += "/" + url.PathEscape(part)
		out = append(out, breadcrumbPart{Name: part, Path: accum + "/"})
	}
	return out
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

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start)

		statusColor := colorGreen
		if wrapped.statusCode >= 400 {
			statusColor = colorYellow
		}
		if wrapped.statusCode >= 500 {
			statusColor = colorRed
		}

		fmt.Printf("  %s%s%d%s %s %s%s%s\n",
			colorDim, statusColor, wrapped.statusCode, colorReset,
			r.URL.Path,
			colorDim, duration.Round(time.Millisecond), colorReset)
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
