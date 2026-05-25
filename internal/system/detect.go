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
	DE          string
	DEVersion   string
	SessionType string // wayland or x11
}

// GetSystemInfo detects the OS, Version, DE, and Session Type.
func GetSystemInfo() Info {
	info := Info{
		OS:          "unknown",
		OSVersion:   "unknown",
		DE:          "unknown",
		DEVersion:   "unknown",
		SessionType: "unknown",
	}

	// 1. Detect OS and Version from /etc/os-release
	if file, err := os.Open("/etc/os-release"); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "ID=") {
				info.OS = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
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
	if val := os.Getenv("XDG_CURRENT_DESKTOP"); val != "" {
		info.DE = strings.ToLower(val)
		if strings.Contains(info.DE, "gnome") {
			info.DE = "gnome"
			// Get Gnome Version
			out, err := exec.Command("gnome-shell", "--version").Output()
			if err == nil {
				// Output format: GNOME Shell 47.0
				parts := strings.Fields(string(out))
				if len(parts) >= 3 {
					info.DEVersion = parts[2]
				}
			}
		} else if strings.Contains(info.DE, "i3") {
			info.DE = "i3"
		} else if strings.Contains(info.DE, "sway") {
			info.DE = "sway"
		}
	}
	
	// Fallback for Wayland specific (Hyprland might not set XDG_CURRENT_DESKTOP properly in all environments)
	if info.DE == "unknown" {
		if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
			info.DE = "hyprland"
		}
	}

	return info
}
