package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

type MainConfig struct {
	Features []string
}

const (
	FeatureBase      = "Base Tools"
	FeatureSoftware  = "Software Installer"
	FeatureDebloat   = "Debloat Gnome"
	FeatureGit       = "Git Setup"
	FeatureGitHub    = "GitHub Setup"
	FeatureAI        = "AI Tools"
	FeatureShell     = "Shell Configuration"
	FeatureHyprland  = "Hyprland Setup"
	FeatureI3        = "i3wm Setup"
	FeatureKeybinds  = "Keybindings"
	FeatureGnomePerf = "GNOME Optimization"
	FeatureFlatpak   = "Flatpak Setup"
)

func RunMainMenu(sysInfo system.Info) (MainConfig, error) {
	var cfg MainConfig

	// Customize keymap for MultiSelect to show "space" instead of "x"
	km := huh.NewDefaultKeyMap()
	km.MultiSelect.Toggle = key.NewBinding(key.WithKeys(" ", "x"), key.WithHelp("space", "toggle"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Linutils Rakesh").
				Description(fmt.Sprintf(
					"OS: %s %s\nDE: %s %s (%s)\n\nSelect the features you want to run.",
					sysInfo.OS, sysInfo.OSVersion, sysInfo.DE, sysInfo.DEVersion, sysInfo.SessionType,
				)),

			huh.NewMultiSelect[string]().
				Title("Main Menu").
				Options(
					huh.NewOption(FeatureBase, FeatureBase),
					huh.NewOption(FeatureSoftware, FeatureSoftware),
					huh.NewOption(FeatureDebloat, FeatureDebloat),
					huh.NewOption(FeatureGit, FeatureGit),
					huh.NewOption(FeatureGitHub, FeatureGitHub),
					huh.NewOption(FeatureAI, FeatureAI),
					huh.NewOption(FeatureShell, FeatureShell),
					huh.NewOption(FeatureHyprland, FeatureHyprland),
					huh.NewOption(FeatureI3, FeatureI3),
					huh.NewOption(FeatureKeybinds, FeatureKeybinds),
					huh.NewOption(FeatureGnomePerf, FeatureGnomePerf),
					huh.NewOption(FeatureFlatpak, FeatureFlatpak),
				).
				Value(&cfg.Features).
				WithKeyMap(km),
		),
	).WithKeyMap(km)

	err := form.Run()
	return cfg, err
}
