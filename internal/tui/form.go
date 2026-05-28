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
	FeatureExit          = "Exit"
)

func RunMainMenu(sysInfo system.Info, state *MainConfig) (MainConfig, error) {
	if len(state.Items) == 0 {
		state.Items = []ListItem{
			// --- System Core ---
			{Key: FeatureInitialSetup, Name: FeatureInitialSetup, Category: "System Core"},
			{Key: FeatureNvidia, Name: FeatureNvidia, Category: "System Core"},
			{Key: FeatureBluetooth, Name: FeatureBluetooth, Category: "System Core"},
			{Key: FeatureFonts, Name: FeatureFonts, Category: "System Core"},
			{Key: FeatureIcons, Name: FeatureIcons, Category: "System Core"},
			{Key: FeatureDebloat, Name: FeatureDebloat, Category: "System Core"},

			// --- Desktop Environment ---
			{Key: FeatureHyprland, Name: FeatureHyprland, Category: "Desktop Environment"},
			{Key: FeatureHyprlandExtra, Name: FeatureHyprlandExtra, Category: "Desktop Environment"},
			{Key: FeatureI3, Name: FeatureI3, Category: "Desktop Environment"},
			{Key: FeatureSDDM, Name: FeatureSDDM, Category: "Desktop Environment"},
			{Key: FeatureFileManagers, Name: FeatureFileManagers, Category: "Desktop Environment"},
			{Key: FeatureKeybinds, Name: FeatureKeybinds, Category: "Desktop Environment"},
			{Key: FeatureGnomePerf, Name: FeatureGnomePerf, Category: "Desktop Environment"},

			// --- Software & Setup ---
			{Key: FeatureBase, Name: FeatureBase, Category: "Software & Setup"},
			{Key: FeatureSoftware, Name: FeatureSoftware, Category: "Software & Setup"},
			{Key: FeatureEditors, Name: FeatureEditors, Category: "Software & Setup"},
			{Key: FeatureGit, Name: FeatureGit, Category: "Software & Setup"},
			{Key: FeatureGitHub, Name: FeatureGitHub, Category: "Software & Setup"},
			{Key: FeatureRepos, Name: FeatureRepos, Category: "Software & Setup"},
			{Key: FeatureFlatpak, Name: FeatureFlatpak, Category: "Software & Setup"},
			{Key: FeatureScripts, Name: FeatureScripts, Category: "Software & Setup"},
			{Key: FeatureShell, Name: FeatureShell, Category: "Software & Setup"},
			{Key: FeatureAlacritty, Name: FeatureAlacritty, Category: "Software & Setup"},

			// --- Theming (Omarchy-style) ---
			{Key: FeatureThemeSwitcher, Name: FeatureThemeSwitcher, Category: "Theming (Omarchy-style)"},
			{Key: FeatureThemeSetup, Name: FeatureThemeSetup, Category: "Theming (Omarchy-style)"},
			{Key: FeatureThemes, Name: FeatureThemes, Category: "Theming (Omarchy-style)"},
			{Key: FeatureDotfiles, Name: FeatureDotfiles, Category: "Theming (Omarchy-style)"},

			// --- Exit ---
			{Key: FeatureExit, Name: FeatureExit, Category: "Other"},
		}
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
