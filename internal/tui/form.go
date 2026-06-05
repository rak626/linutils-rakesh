package tui

import (
	"github.com/rakesh/linutils-rakesh/internal/system"
)

type MainConfig struct {
	Features []string
	Items    []ListItem
}

const (
	FeatureGnomeSetup    = "Full GNOME Desktop Setup"
	FeatureHyprlandSetup = "Full Hyprland Desktop Setup"
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
	FeatureGitCombined   = "Git & GitHub Setup"
	FeatureExit          = "Exit"
)

func RunMainMenu(sysInfo system.Info, state *MainConfig) (MainConfig, error) {
	if len(state.Items) == 0 {
		var items []ListItem

		items = append(items, ListItem{
			Key:         FeatureGnomeSetup,
			Name:        "Full GNOME Desktop Setup",
			Description: "Optimized GNOME with performance tweaks, debloating, custom keybinds, and essential dev tools.",
		})
		items = append(items, ListItem{
			Key:         FeatureHyprlandSetup,
			Name:        "Full Hyprland Desktop Setup",
			Description: "Complete Hyprland environment with Waybar, SDDM, custom configs, and essential dev tools.",
		})
		items = append(items, ListItem{
			Key:         FeatureGitCombined,
			Name:        "Git & GitHub Setup",
			Description: "Configure Git (user.name/email) and authenticate with GitHub CLI.",
		})
		items = append(items, ListItem{
			Key:         FeatureSoftware,
			Name:        "Software Installer",
			Description: "Automatically install all configured software, tools, and editors.",
		})
		items = append(items, ListItem{
			Key:         FeatureExit,
			Name:        "Exit",
			Description: "Close the application.",
		})

		state.Items = items
	}

	_, results, err := RunListUIWithDesc("System Presets", "Select a desktop environment to configure.", state.Items)
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
