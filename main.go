package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/rakesh/linutils-rakesh/internal/modules"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

func main() {
	sysInfo := system.GetSystemInfo()

	// Check for standalone subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "theme":
			manager, _ := pkgmanager.GetManager(sysInfo.OS)
			if err := modules.RunStandaloneThemeSwitcher(manager, sysInfo); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}
	
	manager, err := pkgmanager.GetManager(sysInfo.OS)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Persistent state for selections
	mainConfig := tui.MainConfig{}
	var softwareItems []tui.ListItem

	for {
		cfg, err := tui.RunMainMenu(sysInfo, &mainConfig)
		if err != nil {
			log.Fatal(err)
		}

		if len(cfg.Features) == 0 {
			fmt.Println("No features selected. Use Space to select features.")
			fmt.Println("\nPress Enter to return to menu...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		// Check if "Exit" was chosen
		exitChosen := false
		for _, f := range cfg.Features {
			if f == tui.FeatureExit {
				exitChosen = true
				break
			}
		}
		if exitChosen {
			fmt.Println("Goodbye!")
			break
		}

		for _, feature := range cfg.Features {
			switch feature {
			case tui.FeatureQuickSetup:
				fmt.Println("\n>>> STARTING QUICK SETUP (Non-Theme) <<<")
				modules.RunInitialSetup(manager, sysInfo)
				installBaseTools(manager, sysInfo)
				if sysInfo.DE == "gnome" {
					modules.SetupGnomePerformance()
					modules.SetupGnomeKeybinds()
				}
				modules.SetupFlatpak(manager, sysInfo)
				modules.SetupShell(manager)
				modules.SetupFonts(manager)
				modules.InstallIconAssets(manager)
				modules.SetupEditors(manager)
				fmt.Println("\n>>> QUICK SETUP COMPLETE <<<")

			case tui.FeatureInitialSetup:
				modules.RunInitialSetup(manager, sysInfo)
			case tui.FeatureBase:
				installBaseTools(manager, sysInfo)
			case tui.FeatureSoftware:
				items, _ := modules.InstallSoftware(manager, sysInfo, softwareItems)
				softwareItems = items
			case tui.FeatureDebloat:
				modules.DebloatGnome(manager, sysInfo)
			case tui.FeatureGit:
				modules.SetupGit(manager)
			case tui.FeatureGitHub:
				modules.SetupGitHub(manager)
			case tui.FeatureShell:
				modules.SetupShell(manager)
			case tui.FeatureAlacritty:
				modules.SetupAlacritty(manager)
			case tui.FeatureHyprland:
				modules.SetupHyprland(manager, sysInfo)
			case tui.FeatureHyprlandExtra:
				modules.ConfigureHyprlandExtras(manager)
			case tui.FeatureI3:
				modules.SetupI3(manager, sysInfo)
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
			case tui.FeatureDotfiles:
				modules.SetupDotfiles(manager)
			case tui.FeatureFonts:
				modules.SetupFonts(manager)
			case tui.FeatureIcons:
				modules.InstallIconAssets(manager)
			case tui.FeatureRepos:
				modules.CloneRepos(manager)
			case tui.FeatureNvidia:
				modules.SetupNvidia(manager, sysInfo)
			case tui.FeatureBluetooth:
				modules.SetupBluetoothAndAudio(manager, sysInfo)
			case tui.FeatureSDDM:
				modules.SetupSDDM(manager, sysInfo)
			case tui.FeatureFileManagers:
				modules.SetupFileManagers(manager, sysInfo)
			case tui.FeatureEditors:
				modules.SetupEditors(manager)
			case tui.FeatureScripts:
				modules.InstallCustomScripts(manager)
			case tui.FeatureThemes:
				modules.ApplyThemes(manager)
			case tui.FeatureThemeSwitcher:
				modules.InstallThemeSwitcher(manager)
			case tui.FeatureThemeSetup:
				modules.IntegrateThemeSwitcher()
			case tui.FeatureThemeReset:
				modules.RestoreThemeDefaults(sysInfo)
			}
		}

		fmt.Println("\nSelected tasks complete! Press Enter to return to menu...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func installBaseTools(manager pkgmanager.PackageManager, sysInfo system.Info) {
	fmt.Println("\n--- Installing Base Tools ---")
	manager.Update()
	
	basePkgs := []string{
		"neovim", "grep", "ripgrep", "fzf", "zoxide", "curl", "wget", 
		"git", "vim", "micro", "btop", "htop", "nvtop", "fastfetch", "alacritty", "jq",
	}

	if sysInfo.OS == "debian" || sysInfo.OS == "ubuntu" {
		basePkgs = append(basePkgs, "batcat")
	} else {
		basePkgs = append(basePkgs, "bat")
	}
	
	if err := manager.Install(basePkgs...); err != nil {
		fmt.Printf("Error installing base packages: %v\n", err)
	}
}
