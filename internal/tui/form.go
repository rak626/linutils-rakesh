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
	FeatureSDKMan    = "SDKMan Setup"
	FeatureDotfiles  = "Dotfiles Sync"
	FeatureRepos     = "GitHub Repo Cloner"
	FeatureExit      = "Exit"
)

func RunMainMenu(sysInfo system.Info, state *MainConfig) (MainConfig, error) {
	if len(state.Items) == 0 {
		state.Items = []ListItem{
			{Key: FeatureBase, Name: FeatureBase},
			{Key: FeatureSoftware, Name: FeatureSoftware},
			{Key: FeatureDebloat, Name: FeatureDebloat},
			{Key: FeatureGit, Name: FeatureGit},
			{Key: FeatureGitHub, Name: FeatureGitHub},
			{Key: FeatureAI, Name: FeatureAI},
			{Key: FeatureShell, Name: FeatureShell},
			{Key: FeatureHyprland, Name: FeatureHyprland},
			{Key: FeatureI3, Name: FeatureI3},
			{Key: FeatureKeybinds, Name: FeatureKeybinds},
			{Key: FeatureGnomePerf, Name: FeatureGnomePerf},
			{Key: FeatureFlatpak, Name: FeatureFlatpak},
			{Key: FeatureSDKMan, Name: FeatureSDKMan},
			{Key: FeatureDotfiles, Name: FeatureDotfiles},
			{Key: FeatureRepos, Name: FeatureRepos},
			{Key: FeatureExit, Name: FeatureExit},
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
