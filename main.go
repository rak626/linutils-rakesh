package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rakesh/linutils-rakesh/internal/config"
	"github.com/rakesh/linutils-rakesh/internal/modules"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

func main() {
	sysInfo := system.GetSystemInfo()
	
	cfg, err := tui.RunMainMenu(sysInfo)
	if err != nil {
		log.Fatal(err)
	}

	manager, err := pkgmanager.GetManager(sysInfo.OS)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, feature := range cfg.Features {
		switch feature {
		case tui.FeatureBase:
			installBaseTools(manager, sysInfo)
		case tui.FeatureSoftware:
			modules.InstallSoftware(manager, sysInfo)
		case tui.FeatureDebloat:
			modules.DebloatGnome(manager, sysInfo)
		case tui.FeatureGit:
			modules.SetupGit(manager)
		case tui.FeatureGitHub:
			modules.SetupGitHub(manager)
		case tui.FeatureAI:
			// For now integrated in Software, but can be split
			fmt.Println("AI Tools selection integrated in Software Installer for now.")
		case tui.FeatureShell:
			modules.SetupShell(manager)
		case tui.FeatureHyprland:
			configurator := &config.HyprlandConfigurator{SysInfo: sysInfo}
			configurator.Setup(manager)
		case tui.FeatureI3:
			configurator := &config.I3Configurator{SysInfo: sysInfo}
			configurator.Setup(manager)
		case tui.FeatureKeybinds:
			if err := modules.SetupGnomeKeybinds(); err != nil {
				fmt.Printf("Error setting up keybindings: %v\n", err)
			}
		case tui.FeatureGnomePerf:
			if err := modules.SetupGnomePerformance(); err != nil {
				fmt.Printf("Error setting up GNOME performance: %v\n", err)
			}
		case tui.FeatureFlatpak:
			if err := modules.SetupFlatpak(manager, sysInfo); err != nil {
				fmt.Printf("Error configuring Flatpak: %v\n", err)
			}
		case tui.FeatureSDKMan:
			modules.SetupSDKMan()
		}
	}

	fmt.Println("\nAll selected tasks complete!")
}

func installBaseTools(manager pkgmanager.PackageManager, sysInfo system.Info) {
	fmt.Println("\n--- Installing Base Tools ---")
	manager.Update()
	
	basePkgs := []string{
		"neovim", "grep", "ripgrep", "fzf", "zoxide", "curl", "wget", 
		"git", "vim", "micro", "btop", "htop", "nvtop", "fastfetch", "alacritty", "jq",
	}

	// 'bat' is called 'batcat' on Debian/Ubuntu but 'bat' on others
	if sysInfo.OS == "debian" || sysInfo.OS == "ubuntu" {
		basePkgs = append(basePkgs, "batcat")
	} else {
		basePkgs = append(basePkgs, "bat")
	}
	
	if err := manager.Install(basePkgs...); err != nil {
		fmt.Printf("Error installing base packages: %v\n", err)
	}
}
