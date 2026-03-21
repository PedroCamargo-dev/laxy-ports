# laxy-ports

A lightweight, keyboard-driven TUI for monitoring open ports and killing processes on Linux — reads directly from `/proc`, no external tools required.

```
╔══════════════════════════════════════════════════════════════════════════════════╗
║                             LAXY-PORTS  (TUI)                                  ║
╚══════════════════════════════════════════════════════════════════════════════════╝
┌─ OPEN PORTS (12)──────────────────────────────────────────────────────────────┐
│  PORT   PROTO   PID      PROCESS                                               │
│  ────────────────────────────────────────────────────────────────────────────  │
│✓ 22     TCP     1234     sshd                                                  │
│✓ 80     TCP     5678     nginx                                                 │
│◆ 53     UDP     910      systemd-resolved                                      │
│▶ 443    TCP     5678     nginx                                                 │
│✓ 3000   TCP     2345     [docker] my-app                                       │
│○ 8080   TCP     -        kernel                                                │
└───────────────────────────────────────────────────────────────────────────────┘
 [↑↓/jk] navigate  [K] kill  [R] refresh  [Q] quit
```

## Features

- **Zero external dependencies** — reads `/proc/net/tcp`, `tcp6`, `udp`, `udp6` directly; never shells out to `ss`, `netstat`, or `lsof`
- **Process resolution** — maps socket inodes to PIDs via `/proc/<pid>/fd`, then reads `/proc/<pid>/comm` for the process name
- **Docker enrichment** — queries the Docker socket (`/var/run/docker.sock`) to label ports exposed by containers
- **Kill on demand** — `K` sends `SIGTERM` to the owning process, escalates to `SIGKILL` after 500 ms if still alive
- **Auto-refresh** — list refreshes every 2 seconds; cursor follows the selected entry across refreshes
- **Mouse support** — click to select a row; scroll wheel navigates the list

## Requirements

- Linux (reads `/proc` — macOS and BSD are not supported)
- Go 1.21+ (to build from source)
- Optional: Docker socket at `/var/run/docker.sock` for container labels
- Killing processes requires root privileges or appropriate capabilities

## Installation

### One-liner (Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/PedroCamargo-dev/laxy-ports/main/install.sh | sh
```

Downloads and installs the latest pre-built binary for your architecture (`amd64` or `arm64`) to `/usr/local/bin/laxy-ports`. Override the install directory with `INSTALL_DIR=/your/path`.

### Download manually

Grab the archive for your platform from the [Releases](https://github.com/PedroCamargo-dev/laxy-ports/releases/latest) page:

| Platform | File |
|----------|------|
| Linux x86-64 | `laxy-ports_<version>_linux_amd64.tar.gz` |
| Linux ARM64  | `laxy-ports_<version>_linux_arm64.tar.gz` |

```bash
tar -xzf laxy-ports_*_linux_amd64.tar.gz
sudo mv laxy-ports /usr/local/bin/
```

### From source

```bash
git clone https://github.com/PedroCamargo-dev/laxy-ports.git
cd laxy-ports
go build -o laxy-ports ./cmd/laxy-ports/
sudo mv laxy-ports /usr/local/bin/
```

### Run directly

```bash
go run ./cmd/laxy-ports/
```

> **Note:** killing processes requires root privileges or `CAP_KILL`. Run with `sudo` if needed.

## Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `K` | Kill the process owning the selected port (SIGTERM → SIGKILL) |
| `R` / `r` | Refresh port list immediately |
| `Q` / `ctrl+c` | Quit |

### Mouse

| Action | Effect |
|--------|--------|
| Left click on row | Select that port entry |
| Scroll wheel | Move cursor up / down |

## Architecture

```
laxy-ports/
├── cmd/laxy-ports/     # Entry point — initialises the TUI
└── internal/
    ├── models/         # Shared domain type: PortEntry
    ├── network/        # /proc reader, inode→PID resolution, Docker enrichment
    └── tui/
        ├── model.go    # Bubble Tea Model, Update loop, key & mouse handlers
        ├── view.go     # Port list rendering, status bar
        └── styles.go   # Lipgloss colour palette and style definitions
```

**Design principles:**
- No `os.Exec` — all data comes from `/proc` or the Docker Unix socket
- No polling goroutines — Bubble Tea tick drives periodic refresh; CPU idles between ticks
- No comments — naming is the documentation

## Stack

| | |
|-|-|
| Language | Go 1.24 |
| TUI framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) v0.27 |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) v1.0 |

## Contributing

Contributions are welcome. Please open an issue before submitting a pull request for significant changes.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Open a pull request

## License

[MIT](LICENSE) — © 2026 Pedro Camargo
