<div align="center">

# FsHost

**A fast, zero-dependency file server for your local network.**

Share any folder over WiFi with a single command. Browse, navigate, and download files from any device — phone, tablet, or laptop — through a polished, dark-themed web interface.

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-333?style=for-the-badge)](https://github.com)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](https://opensource.org/licenses/MIT)

</div>

---

## Overview

FsHost is a single static binary that turns any folder into a website for your local network. Point it at a directory, and anyone connected to your WiFi can browse and download its contents from a browser — no accounts, no config files, no dependencies.

```
  ______   _    _           _
 |  ____| | |  | |         | |
 | |__ ___| |__| | ___  ___| |_
 |  __/ __|  __  |/ _ \/ __| __|
 | |  \__ \ |  | | (_) \__ \ |_
 |_|  |___/_|  |_|\___/|___/\__|
```

## Features

- **Single binary, zero dependencies** — statically compiled, nothing to install on the server.
- **Automatic IP detection** — picks your LAN address (`192.168.x.x`) and prints ready-to-open URLs.
- **Dark, responsive web UI** — animated, type-aware file icons, breadcrumbs, and a mobile-friendly layout.
- **Fast and concurrent** — pure Go `net/http`, handles many simultaneous connections and resumeable downloads (HTTP Range).
- **Contained by design** — serves only the directory you specify; path traversal is blocked at the filesystem level.
- **Hidden-file aware** — dotfiles are excluded from listings (but remain fetchable by direct URL if you know the name).

## Requirements

- Go **1.21+** (only needed to build from source)
- A network adapter with a routable IP (WiFi, Ethernet, VPN)

## Installation

### Build from source

```bash
git clone https://github.com/n11kol11c/FsHost.git
cd FsHost
go build -o FsHost .
```

### Install with `go install`

```bash
go install github.com/n11kol11c/FsHost@latest
```

### Embed the version at build time

The version string defaults to `1.0.0` but can be injected:

```bash
go build -ldflags "-X main.version=2.3.1" -o FsHost .
```

## Usage

```
FsHost [flags]

Flags:
  -dir string    Directory to share (default ".")
  -port int      Port to serve on (default 8080)
```

### Examples

```bash
# Share the current directory
FsHost

# Share a specific folder (supports ~)
FsHost -dir ~/Documents

# Serve on a custom port
FsHost -port 3000

# Share a build output for your team
FsHost -dir ./build -port 9090
```

### What you'll see

```
  📂 Serving:  /Users/you/Documents
  🌐 Network:  http://192.168.1.23:8080
  💻 Local:    http://localhost:8080
  ⏱  Started:  2026-08-13 21:30:00
  🛡  OS:       darwin

  ✔ Server is live! Open the link above in any browser on your network.
```

Open the **Network** URL from any device on the same network to browse and download files. Press `Ctrl+C` to stop — FsHost shuts down gracefully and drains in-flight requests.

## Web Interface

| Feature | Description |
|---|---|
| **Breadcrumbs** | Click any ancestor folder to jump back up the path |
| **Parent link** | A `..` entry at the top of each listing |
| **File icons** | Emoji matched to extension — images, video, audio, archives, code, docs, and more |
| **Sizes & dates** | Human-readable sizes (KB / MB / GB) and last-modified timestamps |
| **Direct download** | Click any file to download it (with resume support) |
| **Safe URLs** | Filenames with spaces, `#`, `?`, and other special characters are correctly escaped |
| **Responsive** | Optimized layout for phones, tablets, and desktops |

## Development

```bash
go build ./...     # compile
go vet ./...       # static checks
go build -o FsHost . && ./FsHost -dir ./testdata   # run locally
```

The web UI lives in [`theme/theme.go`](theme/theme.go) as a Go raw string; no build step or bundler is needed.

## Platform Notes

| OS | Notes |
|---|---|
| macOS | Works out of the box; may prompt for network permissions on first run |
| Linux | Works out of the box; use ports > 1024 to avoid needing root |
| Windows | Works out of the box; the firewall may ask to allow inbound access |

## Security Considerations

- FsHost has **no authentication** — anything on the network can read the shared folder. Use it only on networks you trust.
- It refuses to follow paths outside the shared root, even with crafted `..` or encoded traversal attempts.
- Symlinks inside the shared folder are followed like any file manager would.
- **Do not expose FsHost to the public internet.** It is designed for trusted local networks.
- Hidden files are hidden from listings but still reachable by direct URL; do not rely on this as access control.

## License

Released under the [MIT License](https://opensource.org/licenses/MIT).

---

<div align="center">

Built with Go. No frameworks. No bloat. Just works.

</div>
