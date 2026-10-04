package tui

import (
	"os"

	"github.com/rak626/linutils-rakesh/internal/system"
)

type MainConfig struct {
	Features []string
	Items    []ListItem
}

const (
	FeatureGnomeSetup     = "Full GNOME Desktop Setup"
	FeatureHyprlandSetup  = "Full Hyprland Desktop Setup"
	FeatureQuickSetup     = "Full System Setup (Quick)"
	FeatureFreshInstall   = "Fresh Install"
	FeatureRestore        = "Restore My Settings"
	FeatureBackup         = "Backup Now to Git"
	FeatureInitialSetup   = "OS Initial Setup"
	FeatureBase           = "Base Tools"
	FeatureSoftware       = "Software Installer"
	FeatureSoftwarePicker = "Software Picker"
	FeatureDebloat        = "Debloat Gnome"
	FeatureGit            = "Git Setup"
	FeatureGitHub         = "GitHub Setup"
	FeatureShell          = "Shell Configuration"
	FeatureAlacritty      = "Alacritty Setup"
	FeatureHyprland       = "Hyprland Setup"
	FeatureHyprlandExtra  = "Hyprland Extra Config"
	FeatureI3             = "i3wm Setup"
	FeatureKeybinds       = "Keybindings"
	FeatureGnomePerf      = "GNOME Optimization"
	FeatureFlatpak        = "Flatpak Setup"
	FeatureDotfiles       = "Dotfiles Sync"
	FeatureFonts          = "Fonts Setup"
	FeatureIcons          = "Icons & Cursors"
	FeatureRepos          = "GitHub Repo Cloner"
	FeatureNvidia         = "NVIDIA Driver Setup"
	FeatureBluetooth      = "Bluetooth & Audio (Omarchy-style)"
	FeatureSDDM           = "SDDM Login Manager"
	FeatureFileManagers   = "File Managers (Thunar/Yazi)"
	FeatureEditors        = "Editor Config (NVim/Vim)"
	FeatureScripts        = "Custom Scripts"
	FeatureGitCombined    = "Git & GitHub Setup"
	FeatureExit           = "Exit"
)

func IsAdvancedEnabled() bool {
	return os.Getenv("LINUTILS_ADVANCED") == "1"
}

func RunMainMenu(sysInfo system.Info, state *MainConfig) (MainConfig, error) {
	if len(state.Items) == 0 {
		var items []ListItem

		items = append(items, ListItem{
			Key:         FeatureFreshInstall,
			Name:        "Fresh Install",
			Description: "New system: base + GNOME/Hyprland + perf + keybinds + apps. Asks each time.",
		})
		items = append(items, ListItem{
			Key:         FeatureRestore,
			Name:        "Restore My Settings",
			Description: "Fix messed-up system: re-apply dotfiles (git-only), keybinds, perf. Idempotent.",
		})
		items = append(items, ListItem{
			Key:         FeatureBackup,
			Name:        "Backup Now to Git",
			Description: "Dump dconf + package list and commit to dotfiles repo.",
		})
		items = append(items, ListItem{
			Key:         FeatureKeybinds,
			Name:        "Keybindings",
			Description: "Pick terminal/browser/editor, apply GNOME keybinds.",
		})
		items = append(items, ListItem{
			Key:         FeatureSoftwarePicker,
			Name:        "Software Picker",
			Description: "Choose what to install from catalog. Nothing auto-installed.",
		})

		if IsAdvancedEnabled() {
			// Hidden for later — nothing deleted, only shown with LINUTILS_ADVANCED=1.
			advanced := []ListItem{
				{Key: FeatureGnomeSetup, Name: "Full GNOME Desktop Setup", Description: "Legacy preset.", Hidden: false},
				{Key: FeatureHyprlandSetup, Name: "Full Hyprland Desktop Setup", Description: "Legacy preset.", Hidden: false},
				{Key: FeatureQuickSetup, Name: "Full System Setup (Quick)", Description: "Legacy preset.", Hidden: false},
				{Key: FeatureInitialSetup, Name: "OS Initial Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureBase, Name: "Base Tools", Description: "Advanced.", Hidden: false},
				{Key: FeatureSoftware, Name: "Software Installer", Description: "Advanced legacy auto-install.", Hidden: false},
				{Key: FeatureDebloat, Name: "Debloat Gnome", Description: "Advanced.", Hidden: false},
				{Key: FeatureGit, Name: "Git Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureGitHub, Name: "GitHub Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureGitCombined, Name: "Git & GitHub Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureShell, Name: "Shell Configuration", Description: "Advanced.", Hidden: false},
				{Key: FeatureAlacritty, Name: "Alacritty Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureHyprland, Name: "Hyprland Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureHyprlandExtra, Name: "Hyprland Extra Config", Description: "Advanced.", Hidden: false},
				{Key: FeatureI3, Name: "i3wm Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureGnomePerf, Name: "GNOME Optimization", Description: "Advanced.", Hidden: false},
				{Key: FeatureFlatpak, Name: "Flatpak Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureDotfiles, Name: "Dotfiles Sync", Description: "Advanced.", Hidden: false},
				{Key: FeatureFonts, Name: "Fonts Setup", Description: "Advanced.", Hidden: false},
				{Key: FeatureIcons, Name: "Icons & Cursors", Description: "Hidden for later.", Hidden: false},
				{Key: FeatureRepos, Name: "GitHub Repo Cloner", Description: "Hidden for later.", Hidden: false},
				{Key: FeatureNvidia, Name: "NVIDIA Driver Setup", Description: "Hidden for later.", Hidden: true},
				{Key: FeatureBluetooth, Name: "Bluetooth & Audio (Omarchy-style)", Description: "Hidden for later.", Hidden: false},
				{Key: FeatureSDDM, Name: "SDDM Login Manager", Description: "Hidden for later.", Hidden: false},
				{Key: FeatureFileManagers, Name: "File Managers (Thunar/Yazi)", Description: "Hidden for later.", Hidden: false},
				{Key: FeatureEditors, Name: "Editor Config (NVim/Vim)", Description: "Advanced.", Hidden: false},
				{Key: FeatureScripts, Name: "Custom Scripts", Description: "Hidden for later.", Hidden: false},
			}
			items = append(items, advanced...)
		}

		items = append(items, ListItem{
			Key:         FeatureExit,
			Name:        "Exit",
			Description: "Close the application.",
		})

		state.Items = items
	}

	action, results, err := RunListUIWithDesc("System Presets", "Select a desktop environment to configure.", state.Items)
	if err != nil {
		return *state, err
	}

	// q / ctrl+c or Esc on the top-level menu closes the application
	// instead of looping back with "No features selected".
	if action == "quit" || action == "back" {
		state.Items = results
		state.Features = []string{FeatureExit}
		return *state, nil
	}

	state.Items = results
	state.Features = []string{}
	for i, item := range results {
		if item.Selected {
			state.Features = append(state.Features, item.Key)
			// If it was a single execution via Enter (no previous Space selection),
			// don't persist the selected state.
			if action == "i_single" {
				state.Items[i].Selected = false
			}
		}
	}

	return *state, nil
}
