# Project Overview: Linutils Rakesh

A TUI-based Linux tool installer and system configurator written in Go. It provides a user-friendly interface to set up a fresh Linux installation with essential tools, developer software, and customized desktop environment configurations.

## Core Technologies
- **Language:** Go (1.26+)
- **TUI Framework:** [huh](https://github.com/charmbracelet/huh) (Charmbracelet)
- **Target OS:** Linux (Debian/Ubuntu, Arch, Fedora)

## Architecture & Community Domains
The project is organized into several functional domains (communities) that handle specific aspects of the system setup:

- **Package Management Core:** Handles abstraction for `apt`, `pacman`, and `dnf`. Core functions like `runCommand()` serve as bridges across these implementations.
- **System Detection:** Detects OS, version, Desktop Environment (Gnome, Hyprland, i3), and session type (Wayland/X11).
- **Orchestration:** The `main()` function acts as the central bridge, connecting system detection with various installation modules.
- **Software Modules:** Specialized installers for general software, manual `curl`-based installs, and WebApp creation.
- **Environment Configuration:** Modules dedicated to Git, GitHub, and Window Manager (Hyprland/i3) setup.

### Core Abstractions (God Nodes)
- `runCommand()`: The primary execution engine for shell-based tasks.
- `main()`: Orchestrates the entire flow from detection to TUI interaction and execution.
- `PackageManager` (Interface): Abstraction for distro-specific installers (`AptManager`, `PacmanManager`, `DnfManager`).

---

## Building and Running

### Prerequisites
- Go installed (version 1.26 or later recommended)
- Standard Linux build tools

### Build
To build the project into a single binary:
```bash
go build -o linutils-rakesh main.go
```

### Run
To run without building:
```bash
go run main.go
```
Or run the built binary:
```bash
./linutils-rakesh
```

---

## Development Conventions

### Knowledge Graph Maintenance
Always update the knowledge graph using `/graphify --update` after creating or significantly updating code nodes (functions, types, modules) to ensure the project's architectural map remains current.

### Package Management
All package management operations should go through the `PackageManager` interface found in `internal/pkgmanager/manager.go`. This ensures cross-distro compatibility.

### System Detection
Use `system.GetSystemInfo()` from `internal/system/detect.go` to retrieve system-specific metadata before making configuration decisions.

### Adding New Modules
1. Define the module logic in `internal/modules/` or `internal/config/`.
2. Add the module as a constant in `internal/tui/form.go`.
3. Update the `RunMainMenu` form in `internal/tui/form.go` to include the new option.
4. Update the `switch` statement in `main.go` to handle the new feature.

### UI Style
Maintain a consistent TUI style using the `huh` library. Group related options and provide clear descriptions for each task.
