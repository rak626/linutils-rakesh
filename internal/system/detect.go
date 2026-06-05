package system

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

type Info struct {
	OS          string
	OSVersion   string
	DistroID    string // ID from /etc/os-release (e.g., fedora, ubuntu)
	DE          string // Pretty name for display
	DEID        string // Normalized ID for logic (e.g., gnome, plasma, hyprland)
	DEVersion   string
	SessionType string // wayland or x11
	CPU         string
	RAM         string
	Disk        string
	GPU         string
}

// GetSystemInfo detects the OS, Version, DE, Session Type, and Hardware.
func GetSystemInfo() Info {
	info := Info{
		OS:          "unknown",
		OSVersion:   "unknown",
		DistroID:    "unknown",
		DE:          "unknown",
		DEID:        "unknown",
		DEVersion:   "unknown",
		SessionType: "unknown",
		CPU:         "unknown",
		RAM:         "unknown",
		Disk:        "unknown",
		GPU:         "unknown",
	}

	// 1. Detect OS and Version from /etc/os-release
	if file, err := os.Open("/etc/os-release"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "ID=") {
				info.DistroID = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
				if info.OS == "unknown" {
					info.OS = info.DistroID
				}
			}
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				info.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			}
			if strings.HasPrefix(line, "VERSION_ID=") {
				info.OSVersion = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
			}
		}
	}

	// 2. Detect Session Type
	if val := os.Getenv("XDG_SESSION_TYPE"); val != "" {
		info.SessionType = strings.ToLower(val)
	}

	// 3. Detect DE and Version
	rawDE := os.Getenv("XDG_CURRENT_DESKTOP")
	if rawDE != "" {
		info.DE = rawDE
		lowerDE := strings.ToLower(rawDE)
		if strings.Contains(lowerDE, "gnome") {
			info.DE = "GNOME"
			info.DEID = "gnome"
			// Get Gnome Version
			out, err := exec.Command("gnome-shell", "--version").Output()
			if err == nil {
				// Output format: GNOME Shell 47.0
				parts := strings.Fields(string(out))
				if len(parts) >= 3 {
					info.DEVersion = parts[2]
				}
			}
		} else if strings.Contains(lowerDE, "hyprland") {
			info.DE = "Hyprland"
			info.DEID = "hyprland"
		} else if strings.Contains(lowerDE, "i3") {
			info.DE = "i3wm"
			info.DEID = "i3"
			out, err := exec.Command("i3", "--version").Output()
			if err == nil {
				// i3 version 4.22 (2023-01-02)
				parts := strings.Fields(string(out))
				if len(parts) >= 3 {
					info.DEVersion = parts[2]
				}
			}
		} else if strings.Contains(lowerDE, "sway") {
			info.DE = "Sway"
			info.DEID = "sway"
			out, err := exec.Command("sway", "--version").Output()
			if err == nil {
				// sway version 1.8.1
				parts := strings.Fields(string(out))
				if len(parts) >= 3 {
					info.DEVersion = parts[2]
				}
			}
		} else if strings.Contains(lowerDE, "kde") || strings.Contains(lowerDE, "plasma") {
			info.DE = "KDE Plasma"
			info.DEID = "plasma"
			out, err := exec.Command("plasmashell", "--version").Output()
			if err == nil {
				// plasmashell 6.1.4
				parts := strings.Fields(string(out))
				if len(parts) >= 2 {
					info.DEVersion = parts[1]
				}
			}
		} else {
			info.DEID = lowerDE
		}
	}
	
	// Fallback for Wayland specific (Hyprland might not set XDG_CURRENT_DESKTOP properly in all environments)
	if info.DEID == "unknown" || info.DEID == "hyprland" {
		if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" || info.DEID == "hyprland" {
			info.DE = "Hyprland"
			info.DEID = "hyprland"
			// Get Hyprland Version
			out, err := exec.Command("hyprctl", "version").Output()
			if err == nil {
				// Output contains "Tag: v0.41.2" or similar
				lines := strings.Split(string(out), "\n")
				for _, line := range lines {
					if strings.Contains(line, "Tag:") {
						parts := strings.Fields(line)
						if len(parts) >= 2 {
							info.DEVersion = parts[1]
							break
						}
					}
				}
			}
		}
	}

	// 4. Hardware Detection
	// CPU
	if out, err := exec.Command("bash", "-c", "grep -m 1 'model name' /proc/cpuinfo | cut -d: -f2 | sed 's/^[ \t]*//'").Output(); err == nil {
		info.CPU = strings.TrimSpace(string(out))
	}

	// RAM - More accurate total RAM
	if out, err := exec.Command("bash", "-c", "grep MemTotal /proc/meminfo | awk '{printf \"%.1fGi\", $2/1024/1024}'").Output(); err == nil {
		info.RAM = strings.TrimSpace(string(out))
	}

	// Disk
	if out, err := exec.Command("bash", "-c", "df -h / | tail -1 | awk '{print $2}'").Output(); err == nil {
		info.Disk = strings.TrimSpace(string(out))
	}

	// GPU
	if out, err := exec.Command("bash", "-c", "lspci | grep -i vga | cut -d: -f3 | sed 's/^[ \t]*//' | head -n 1").Output(); err == nil {
		info.GPU = strings.TrimSpace(string(out))
	}

	return info
}

