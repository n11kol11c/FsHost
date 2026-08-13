package theme

const PageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&family=Inter:wght@400;500;600;700&display=swap');

        :root {
            --bg-primary: #0a0a0f;
            --bg-secondary: #12121a;
            --bg-tertiary: #1a1a2e;
            --bg-hover: #1e1e36;
            --border: #2a2a40;
            --accent: #6c5ce7;
            --accent-light: #a29bfe;
            --accent-glow: rgba(108, 92, 231, 0.3);
            --text-primary: #e2e8f0;
            --text-secondary: #94a3b8;
            --text-muted: #64748b;
            --green: #00b894;
            --cyan: #00cec9;
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg-primary);
            color: var(--text-primary);
            min-height: 100vh;
            line-height: 1.6;
        }

        .bg-grid {
            position: fixed; top: 0; left: 0; width: 100%; height: 100%;
            background-image:
                linear-gradient(rgba(108, 92, 231, 0.03) 1px, transparent 1px),
                linear-gradient(90deg, rgba(108, 92, 231, 0.03) 1px, transparent 1px);
            background-size: 50px 50px;
            pointer-events: none; z-index: 0;
        }

        .bg-glow {
            position: fixed; top: -50%; left: -50%; width: 200%; height: 200%;
            background: radial-gradient(circle at 50% 50%, rgba(108, 92, 231, 0.06) 0%, transparent 50%);
            pointer-events: none; z-index: 0;
            animation: glowPulse 8s ease-in-out infinite alternate;
        }

        @keyframes glowPulse { 0% { opacity: 0.5; } 100% { opacity: 1; } }

        .container {
            position: relative; z-index: 1;
            max-width: 1100px; margin: 0 auto; padding: 24px 20px;
        }

        header {
            display: flex; align-items: center; justify-content: space-between;
            margin-bottom: 32px; padding: 20px 28px;
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 16px; backdrop-filter: blur(20px);
        }

        .logo { display: flex; align-items: center; gap: 14px; }

        .logo-icon {
            width: 44px; height: 44px;
            background: linear-gradient(135deg, var(--accent), var(--accent-light));
            border-radius: 12px; display: flex; align-items: center; justify-content: center;
            font-size: 22px; box-shadow: 0 0 20px var(--accent-glow);
        }

        .logo h1 {
            font-family: 'JetBrains Mono', monospace;
            font-size: 1.5rem; font-weight: 700;
            background: linear-gradient(135deg, var(--accent-light), var(--cyan));
            -webkit-background-clip: text; -webkit-text-fill-color: transparent;
            background-clip: text; letter-spacing: -0.5px;
        }

        .logo .version {
            font-size: 0.7rem; color: var(--text-muted);
            font-family: 'JetBrains Mono', monospace;
        }

        .header-info { display: flex; align-items: center; gap: 16px; }

        .status-badge {
            display: flex; align-items: center; gap: 8px;
            padding: 6px 14px; background: rgba(0, 184, 148, 0.1);
            border: 1px solid rgba(0, 184, 148, 0.2); border-radius: 100px;
            font-size: 0.75rem; font-weight: 500; color: var(--green);
        }

        .status-dot {
            width: 7px; height: 7px; background: var(--green);
            border-radius: 50%; animation: pulse 2s ease-in-out infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; box-shadow: 0 0 0 0 rgba(0, 184, 148, 0.4); }
            50% { opacity: 0.8; box-shadow: 0 0 0 6px rgba(0, 184, 148, 0); }
        }

        .breadcrumb {
            display: flex; align-items: center; gap: 4px; flex-wrap: wrap;
            padding: 16px 28px; margin-bottom: 24px;
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 12px; font-family: 'JetBrains Mono', monospace; font-size: 0.82rem;
        }

        .breadcrumb a {
            color: var(--accent-light); text-decoration: none;
            padding: 3px 8px; border-radius: 6px; transition: all 0.2s;
        }
        .breadcrumb a:hover { background: var(--bg-tertiary); color: var(--cyan); }
        .breadcrumb .sep { color: var(--text-muted); font-size: 0.7rem; }
        .breadcrumb .current { color: var(--text-primary); font-weight: 600; }

        .file-table {
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 16px; overflow: hidden; backdrop-filter: blur(20px);
        }

        .file-table table { width: 100%; border-collapse: collapse; }

        .file-table thead th {
            padding: 14px 24px; text-align: left;
            font-size: 0.72rem; font-weight: 600; text-transform: uppercase;
            letter-spacing: 1px; color: var(--text-muted);
            background: var(--bg-tertiary); border-bottom: 1px solid var(--border);
        }

        .file-table tbody tr {
            transition: all 0.15s ease; border-bottom: 1px solid rgba(42, 42, 64, 0.4);
        }
        .file-table tbody tr:last-child { border-bottom: none; }
        .file-table tbody tr:hover { background: var(--bg-hover); }
        .file-table td { padding: 12px 24px; font-size: 0.88rem; }

        .file-name { display: flex; align-items: center; gap: 12px; }
        .file-name a {
            color: var(--text-primary); text-decoration: none;
            font-weight: 500; transition: color 0.2s;
        }
        .file-name a:hover { color: var(--accent-light); }
        .file-name .icon { font-size: 1.2rem; width: 28px; text-align: center; flex-shrink: 0; }
        .file-size { font-family: 'JetBrains Mono', monospace; font-size: 0.8rem; color: var(--text-secondary); }
        .file-date { color: var(--text-muted); font-size: 0.82rem; }
        .file-ext {
            display: inline-block; padding: 2px 8px; background: var(--bg-tertiary);
            border-radius: 6px; font-family: 'JetBrains Mono', monospace;
            font-size: 0.7rem; color: var(--text-muted); border: 1px solid var(--border);
        }

        .empty-state { text-align: center; padding: 60px 20px; color: var(--text-muted); }
        .empty-state .empty-icon { font-size: 3rem; margin-bottom: 16px; opacity: 0.5; }
        .empty-state h3 { font-size: 1.1rem; color: var(--text-secondary); margin-bottom: 8px; }

        footer {
            margin-top: 32px; padding: 16px 28px;
            background: var(--bg-secondary); border: 1px solid var(--border);
            border-radius: 12px; display: flex; align-items: center;
            justify-content: space-between; font-size: 0.75rem; color: var(--text-muted);
        }
        footer a { color: var(--accent-light); text-decoration: none; }
        footer a:hover { color: var(--cyan); }

        .stats { display: flex; gap: 20px; align-items: center; }
        .stat { display: flex; align-items: center; gap: 6px; }

        @media (max-width: 768px) {
            header { flex-direction: column; gap: 16px; text-align: center; }
            .header-info { justify-content: center; }
            .file-table thead { display: none; }
            .file-table tbody tr { display: flex; flex-direction: column; padding: 16px 20px; gap: 6px; }
            .file-table td { padding: 0; }
            .file-table td:nth-child(3), .file-table td:nth-child(4) { display: none; }
            footer { flex-direction: column; gap: 12px; text-align: center; }
        }

        tr { animation: fadeIn 0.2s ease backwards; }
        tr:nth-child(1) { animation-delay: 0.02s; }
        tr:nth-child(2) { animation-delay: 0.04s; }
        tr:nth-child(3) { animation-delay: 0.06s; }
        tr:nth-child(4) { animation-delay: 0.08s; }
        tr:nth-child(5) { animation-delay: 0.10s; }
        tr:nth-child(6) { animation-delay: 0.12s; }
        tr:nth-child(7) { animation-delay: 0.14s; }
        tr:nth-child(8) { animation-delay: 0.16s; }
        tr:nth-child(9) { animation-delay: 0.18s; }
        tr:nth-child(10) { animation-delay: 0.20s; }

        @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

        .back-row td { border-bottom: 1px solid var(--border) !important; }
    </style>
</head>
<body>
    <div class="bg-grid"></div>
    <div class="bg-glow"></div>

    <div class="container">
        <header>
            <div class="logo">
                <div class="logo-icon">&#128193;</div>
                <div>
                    <h1>FsHost</h1>
                    <span class="version">v{{.Version}}</span>
                </div>
            </div>
            <div class="header-info">
                <div class="status-badge">
                    <span class="status-dot"></span>
                    Serving on {{.IPAddress}}
                </div>
            </div>
        </header>

        <nav class="breadcrumb">
            <a href="/">&#127968; root</a>
            {{range .Bread}}
                <span class="sep">&#8250;</span>
                <a href="{{.Path}}">{{.Name}}</a>
            {{end}}
        </nav>

        <div class="file-table">
            {{if .Files}}
            <table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Size</th>
                        <th>Modified</th>
                        <th>Type</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Files}}
                    <tr{{if eq .Name ".."}} class="back-row"{{end}}>
                        <td>
                            <div class="file-name">
                                <span class="icon">{{.Icon}}</span>
                                {{if eq .Name ".."}}
                                    <a href="{{.Href}}">..</a>
                                {{else if .IsDir}}
                                    <a href="{{.Href}}">{{.Name}}/</a>
                                {{else}}
                                    <a href="{{.Href}}" download>{{.Name}}</a>
                                {{end}}
                            </div>
                        </td>
                        <td class="file-size">{{.Size}}</td>
                        <td class="file-date">{{.ModTime}}</td>
                        <td>
                            {{if .Ext}}
                                <span class="file-ext">{{.Ext}}</span>
                            {{end}}
                        </td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
            {{else}}
            <div class="empty-state">
                <div class="empty-icon">&#128381;</div>
                <h3>This folder is empty</h3>
                <p>No files or folders to display</p>
            </div>
            {{end}}
        </div>

        <footer>
            <div class="stats">
                <div class="stat">
                    <span>&#128193;</span>
                    {{.NumFiles}} items
                </div>
                <div class="stat">
                    <span>&#127760;</span>
                    {{.IPAddress}}
                </div>
            </div>
            <div>
                Powered by <a href="#">FsHost</a> &mdash; Fast local file sharing
            </div>
        </footer>
    </div>
</body>
</html>`
