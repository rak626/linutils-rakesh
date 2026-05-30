package tui

import (
	"fmt"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

type MainConfig struct {
	Features []string
	Items    []ListItem
}

const (
	FeatureQuickSetup    = "Full System Setup (Quick)"
	FeatureInitialSetup  = "OS Initial Setup"
	FeatureBase          = "Base Tools"
	FeatureSoftware      = "Software Installer"
	FeatureDebloat       = "Debloat Gnome"
	FeatureGit           = "Git Setup"
	FeatureGitHub        = "GitHub Setup"
	FeatureShell         = "Shell Configuration"
	FeatureAlacritty     = "Alacritty Setup"
	FeatureHyprland      = "Hyprland Setup"
	FeatureHyprlandExtra = "Hyprland Extra Config"
	FeatureI3            = "i3wm Setup"
	FeatureKeybinds      = "Keybindings"
	FeatureGnomePerf     = "GNOME Optimization"
	FeatureFlatpak       = "Flatpak Setup"
	FeatureDotfiles      = "Dotfiles Sync"
	FeatureFonts         = "Fonts Setup"
	FeatureIcons         = "Icons & Cursors"
	FeatureRepos         = "GitHub Repo Cloner"
	FeatureNvidia        = "NVIDIA Driver Setup"
	FeatureBluetooth     = "Bluetooth & Audio (Omarchy-style)"
	FeatureSDDM          = "SDDM Login Manager"
	FeatureFileManagers  = "File Managers (Thunar/Yazi)"
	FeatureEditors       = "Editor Config (NVim/Vim)"
	FeatureScripts       = "Custom Scripts"
	FeatureThemes        = "Application Themes"
	FeatureThemeSwitcher = "Install Global Theme Switcher"
	FeatureThemeSetup    = "Integrate Theme Switcher with Configs"
	FeatureThemeReset    = "Restore Original Configs (Reset Themes)"
	FeatureExit          = "Exit"
)

func RunMainMenu(sysInfo system.Info, state *MainConfig) (MainConfig, error) {
	if len(state.Items) == 0 {
		var items []ListItem

		// --- Quick Actions ---
		items = append(items, ListItem{
			Key:         FeatureQuickSetup,
			Name:        FeatureQuickSetup,
			Category:    "Quick Actions",
			Description: "One-click install: OS Setup, Base Tools, Flatpak, Shell, Fonts, Icons, and Editors.",
		})

		// --- System Core ---
		items = append(items, ListItem{Key: FeatureInitialSetup, Name: FeatureInitialSetup, Category: "System Core", Description: "DNF/Pacman/Apt optimization, RPM Fusion, and system updates."})
		items = append(items, ListItem{Key: FeatureNvidia, Name: FeatureNvidia, Category: "System Core", Description: "Install NVIDIA drivers and configuration for your GPU."})
		items = append(items, ListItem{Key: FeatureBluetooth, Name: FeatureBluetooth, Category: "System Core", Description: "Enable Bluetooth services and audio improvements."})
		items = append(items, ListItem{Key: FeatureFonts, Name: FeatureFonts, Category: "System Core", Description: "Install JetBrains Mono, Nerd Fonts, and emoji sets."})
		items = append(items, ListItem{Key: FeatureIcons, Name: FeatureIcons, Category: "System Core", Description: "Install Papirus icons and custom cursor themes."})
		
		if sysInfo.DE == "gnome" {
			items = append(items, ListItem{Key: FeatureDebloat, Name: FeatureDebloat, Category: "System Core", Description: "Remove pre-installed GNOME apps and services."})
		}

		// --- Desktop Environment ---
		items = append(items, ListItem{Key: FeatureHyprland, Name: FeatureHyprland, Category: "Desktop Environment"})
		items = append(items, ListItem{Key: FeatureHyprlandExtra, Name: FeatureHyprlandExtra, Category: "Desktop Environment"})
		items = append(items, ListItem{Key: FeatureI3, Name: FeatureI3, Category: "Desktop Environment"})
		items = append(items, ListItem{Key: FeatureSDDM, Name: FeatureSDDM, Category: "Desktop Environment"})
		items = append(items, ListItem{Key: FeatureFileManagers, Name: FeatureFileManagers, Category: "Desktop Environment"})
		
		if sysInfo.DE == "gnome" {
			items = append(items, ListItem{Key: FeatureKeybinds, Name: FeatureKeybinds, Category: "Desktop Environment"})
			items = append(items, ListItem{Key: FeatureGnomePerf, Name: FeatureGnomePerf, Category: "Desktop Environment"})
		}

		// --- Software & Setup ---
		items = append(items, ListItem{Key: FeatureBase, Name: FeatureBase, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureSoftware, Name: FeatureSoftware, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureEditors, Name: FeatureEditors, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureGit, Name: FeatureGit, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureGitHub, Name: FeatureGitHub, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureRepos, Name: FeatureRepos, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureFlatpak, Name: FeatureFlatpak, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureScripts, Name: FeatureScripts, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureShell, Name: FeatureShell, Category: "Software & Setup"})
		items = append(items, ListItem{Key: FeatureAlacritty, Name: FeatureAlacritty, Category: "Software & Setup"})

		// --- Theming (Omarchy-style) ---
		items = append(items, ListItem{Key: FeatureThemeSwitcher, Name: FeatureThemeSwitcher, Category: "Theming (Omarchy-style)"})
		items = append(items, ListItem{Key: FeatureThemeSetup, Name: FeatureThemeSetup, Category: "Theming (Omarchy-style)"})
		items = append(items, ListItem{Key: FeatureThemeReset, Name: FeatureThemeReset, Category: "Theming (Omarchy-style)", Description: "Restore original application configurations and GNOME theme defaults."})
		items = append(items, ListItem{Key: FeatureThemes, Name: FeatureThemes, Category: "Theming (Omarchy-style)"})
		items = append(items, ListItem{Key: FeatureDotfiles, Name: FeatureDotfiles, Category: "Theming (Omarchy-style)"})

		// --- Exit ---
		items = append(items, ListItem{Key: FeatureExit, Name: FeatureExit, Category: "Other"})

		state.Items = items
	}

	desc := fmt.Sprintf(
		"OS: %s %s\nDE: %s %s (%s)\n\nSelect the features you want to run.",
		sysInfo.OS, sysInfo.OSVersion, sysInfo.DE, sysInfo.DEVersion, sysInfo.SessionType,
	)

	_, results, err := RunListUIWithDesc("Linutils Rakesh", desc, state.Items)
	if err != nil {
		return *state, err
	}

	state.Items = results
	state.Features = []string{}
	for _, item := range results {
		if item.Selected {
			state.Features = append(state.Features, item.Key)
		}
	}

	return *state, nil
}
