# 🚀 Linutils Rakesh

A powerful, TUI-based Linux system orchestrator and global theme switcher inspired by the **Omarchy** aesthetic. Built with Go and the Charmbracelet `huh` library, this tool transforms a fresh Linux installation into a fully-configured, high-performance workstation with a single command.

---

## ✨ Features

### 🎨 The Ultimate Theme Orchestrator
Sync your entire aesthetic across **15+ applications** instantly. Selecting a theme (like *Rose Pine*, *Everforest*, or *Catppuccin*) cascades your colors to:
- **Terminals**: Alacritty, Kitty, Ghostty
- **Editors**: Neovim, Vim, Zed, VSCodium
- **System UI**: Hyprland (Borders), i3 (Window Decos), GTK (Apps), GNOME Shell, SDDM
- **Utilities**: Waybar, Mako, btop, Starship, Ulauncher, SwayOSD, Hyprlock
- **Icons & Cursors**: Automated `gsettings` sync for matching icon packs and mouse pointers.

### 🛠️ System Core & Hardware
- **Arch + Hyprland Focus**: Automated setup for Hyprland with specific **NVIDIA/DKMS** performance optimizations.
- **TUI Hardware Managers**: Keyboard-driven Bluetooth and Audio selection (Omarchy-style).
- **Package Management**: Cross-distro abstraction (Apt, DNF, Pacman) with seamless **AUR integration** (yay/paru).
- **Initial Setup**: One-click optimization (DNF speedup, Reflector, DNS config, Debloating).

### 📦 Software & Workflow
- **Categorized Installer**: Browse and install developer tools, AI agents, and Flatpaks.
- **Zero-Config Integration**: One-click tool to "plumb" the theme switcher into your existing `.lua` and `.conf` files.
- **Custom Scripts**: Deploy pre-configured utility scripts for screenshots, power menus, and more.

---

## 🚀 Quick Start

Run the installer directly from the web:

```bash
curl -fsSL https://raw.githubusercontent.com/rak626/linutils-rakesh/main/install.sh | bash
```

Alternatively, build from source:

```bash
git clone https://github.com/rak626/linutils-rakesh.git
cd linutils-rakesh
go build -o linutils-rakesh main.go
./linutils-rakesh
```

---

## ⌨️ Standalone Theme Switcher

Once installed, you can launch the theme switcher instantly via a keybind (default: `Super + Alt + T`) or via CLI:

```bash
linutils-rakesh theme
```

---

## 🏗️ Architecture

- **Go**: Core logic and orchestration.
- **Huh (Charmbracelet)**: Beautiful, accessible TUI components.
- **Environment Aware**: Detects your DE (GNOME, Hyprland, i3) and applies environment-specific tweaks.
- **Stateless/Dynamic**: Generates `active_theme` files that your configs source, ensuring live-reloading without restarts.

---

## 🤝 Contributing

Community themes are welcome! You can import them directly from the TUI by pasting a GitHub URL.

---

> *"Where system configuration meets aesthetic perfection—effortless orchestration for the modern Linux pioneer."*
