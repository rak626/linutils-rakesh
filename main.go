package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/rak626/linutils-rakesh/internal/config"
	"github.com/rak626/linutils-rakesh/internal/modules"
	"github.com/rak626/linutils-rakesh/internal/pkgmanager"
	"github.com/rak626/linutils-rakesh/internal/system"
	"github.com/rak626/linutils-rakesh/internal/tui"
)

func main() {
	config.LoadVariables()
	sysInfo := system.GetSystemInfo()

	// Check for standalone subcommands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		// Add future subcommands here
		}
	}

	manager, err := pkgmanager.GetManager(sysInfo.DistroID)
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

		// Ensure sudo privileges once before starting automated operations
		if err := pkgmanager.ValidateSudo(); err != nil {
			fmt.Printf("Error: sudo validation failed: %v\n", err)
			fmt.Println("\nPress Enter to return to menu...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		for _, feature := range cfg.Features {
			switch feature {
			case tui.FeatureFreshInstall:
				runFreshInstall(manager, sysInfo)
			case tui.FeatureRestore:
				runRestore(manager, sysInfo)
			case tui.FeatureBackup:
				if err := modules.BackupToGit(manager, sysInfo); err != nil {
					fmt.Printf("Backup failed: %v\n", err)
				}
			case tui.FeatureSoftwarePicker:
				if err := modules.InstallSoftwarePicker(manager, sysInfo); err != nil {
					fmt.Printf("Software picker failed: %v\n", err)
				}
			case tui.FeatureGnomeSetup:
				fmt.Println("\n>>> STARTING FULL GNOME SETUP <<<")
				modules.RunInitialSetup(manager, sysInfo)
				installBaseTools(manager, sysInfo)
				modules.DebloatGnome(manager, sysInfo)
				modules.SetupDotfiles(manager)
				modules.SetupGnomePerformance()
				modules.SetupGnomeKeybinds()
				// NVIDIA hidden for later — kept in code, not run by default.
				modules.SetupShell(manager)
				modules.SetupFonts(manager)
				modules.SetupEditors(manager)
				modules.InstallSoftware(manager, sysInfo, nil)
				fmt.Println("\n>>> FULL GNOME SETUP COMPLETE <<<")

			case tui.FeatureHyprlandSetup:
				if !hyprlandAllowed(sysInfo) {
					fmt.Println("Hyprland is only supported on Arch/Fedora. Blocked on Ubuntu/Debian.")
					continue
				}
				fmt.Println("\n>>> STARTING FULL HYPRLAND SETUP <<<")
				modules.RunInitialSetup(manager, sysInfo)
				installBaseTools(manager, sysInfo)
				modules.SetupHyprland(manager, sysInfo)
				modules.SetupDotfiles(manager)
				modules.ConfigureHyprlandExtras(manager)
				modules.SetupSDDM(manager, sysInfo)
				modules.SetupShell(manager)
				modules.SetupFonts(manager)
				modules.SetupEditors(manager)
				modules.InstallSoftware(manager, sysInfo, nil)
				fmt.Println("\n>>> FULL HYPRLAND SETUP COMPLETE <<<")

			case tui.FeatureQuickSetup:
				fmt.Println("\n>>> STARTING QUICK SETUP <<<")
				modules.RunInitialSetup(manager, sysInfo)
				installBaseTools(manager, sysInfo)
				if sysInfo.DEID == "gnome" {
					modules.SetupGnomePerformance()
					modules.SetupGnomeKeybinds()
				}
				modules.SetupFlatpak(manager, sysInfo)
				modules.SetupShell(manager)
				modules.SetupFonts(manager)
				modules.SetupEditors(manager)
				modules.InstallSoftware(manager, sysInfo, nil)
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
			case tui.FeatureGitCombined:
				modules.SetupGit(manager)
				modules.SetupGitHub(manager)
			case tui.FeatureShell:
				modules.SetupShell(manager)
			case tui.FeatureAlacritty:
				modules.SetupAlacritty(manager)
			case tui.FeatureHyprland:
				if !hyprlandAllowed(sysInfo) {
					fmt.Println("Hyprland is only supported on Arch/Fedora. Blocked on Ubuntu/Debian.")
					continue
				}
				modules.SetupHyprland(manager, sysInfo)
			case tui.FeatureHyprlandExtra:
				modules.ConfigureHyprlandExtras(manager)
			case tui.FeatureI3:
				modules.SetupI3(manager, sysInfo)
			case tui.FeatureKeybinds:
				if err := modules.RunInteractiveGnomeKeybinds(); err != nil {
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
				// Hidden for later — kept in code, not run by default.
				fmt.Println("NVIDIA setup is hidden for later (LINUTILS_ADVANCED=1). Skipping.")
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
			}
		}

		fmt.Println("\n[OK] Selected tasks complete. Press Enter to return to menu...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func hyprlandAllowed(sysInfo system.Info) bool {
	id := sysInfo.DistroID
	return id == "arch" || id == "manjaro" || id == "endeavouros" || id == "fedora" || id == "nobara"
}

func askDesktopEnv(sysInfo system.Info) string {
	items := []tui.ListItem{
		{Key: "gnome", Name: "GNOME (primary)", Description: "Perf tweaks + debloat apps + keybinds."},
		{Key: "i3", Name: "i3 (sometimes)", Description: "i3wm packages + dotfiles stow."},
		{Key: "hyprland", Name: "Hyprland (Arch/Fedora only)", Description: "Blocked on Ubuntu/Debian."},
	}
	action, results, err := tui.RunListUIWithDesc("Choose Desktop", "Fresh install target. Current: "+sysInfo.DE, items)
	if err != nil || action == "" || action == "back" || action == "quit" {
		return "gnome"
	}
	for _, it := range results {
		if it.Selected {
			return it.Key
		}
	}
	return "gnome"
}

func runFreshInstall(manager pkgmanager.PackageManager, sysInfo system.Info) {
	de := askDesktopEnv(sysInfo)
	if de == "hyprland" && !hyprlandAllowed(sysInfo) {
		fmt.Println("Hyprland is only supported on Arch/Fedora. Blocked on Ubuntu/Debian.")
		return
	}
	fmt.Printf("\n>>> FRESH INSTALL (%s on %s) <<<\n", de, sysInfo.DistroID)
	steps := []string{"OS setup", "Base tools", "Desktop", "Dotfiles", "Perf+Keybinds", "Shell/Fonts/Editors", "Git", "Software"}
	total := len(steps)
	done := 0

	modules.RunInitialSetup(manager, sysInfo)
	done++
	fmt.Printf("%s step %d/%d done: OS setup\n", tui.MarkOK, done, total)

	installBaseTools(manager, sysInfo)
	done++
	fmt.Printf("%s step %d/%d done: Base tools\n", tui.MarkOK, done, total)

	switch de {
	case "gnome":
		modules.DebloatGnome(manager, sysInfo)
		modules.SetupDotfiles(manager)
		modules.SetupGnomePerformance()
		if err := modules.RunInteractiveGnomeKeybinds(); err != nil {
			fmt.Printf("Keybinds: %v\n", err)
		}
	case "i3":
		modules.SetupI3(manager, sysInfo)
		modules.SetupDotfiles(manager)
	case "hyprland":
		modules.SetupHyprland(manager, sysInfo)
		modules.SetupDotfiles(manager)
		modules.ConfigureHyprlandExtras(manager)
	}
	done++
	fmt.Printf("%s step %d/%d done: Desktop\n", tui.MarkOK, done, total)
	done++
	fmt.Printf("%s step %d/%d done: Dotfiles\n", tui.MarkOK, done, total)
	done++
	fmt.Printf("%s step %d/%d done: Perf+Keybinds\n", tui.MarkOK, done, total)

	modules.SetupFlatpak(manager, sysInfo)
	modules.SetupShell(manager)
	modules.SetupFonts(manager)
	modules.SetupEditors(manager)
	done++
	fmt.Printf("%s step %d/%d done: Shell/Fonts/Editors\n", tui.MarkOK, done, total)

	modules.SetupGit(manager)
	// GitHub CLI auth only on fresh install, skipped if already authed.
	if modules.IsGitHubAuthenticated() {
		fmt.Println("[OK] GitHub already authenticated — skipping gh login.")
	} else {
		modules.SetupGitHub(manager)
	}
	done++
	fmt.Printf("%s step %d/%d done: Git\n", tui.MarkOK, done, total)

	if err := modules.InstallSoftwarePicker(manager, sysInfo); err != nil {
		fmt.Printf("Software picker: %v\n", err)
	}
	done++
	fmt.Printf("%s step %d/%d done: Software\n", tui.MarkOK, done, total)
	fmt.Printf("\n%s FRESH INSTALL COMPLETE %s\n", tui.MarkOK, tui.ProgressBar(total, total, 20))
}

func runRestore(manager pkgmanager.PackageManager, sysInfo system.Info) {
	fmt.Println("\n>>> RESTORE MY SETTINGS (idempotent, git-only) <<<")
	// No OS initial setup, no base reinstall unless missing, no gh unless missing.
	modules.SetupDotfiles(manager)
	if sysInfo.DEID == "gnome" {
		modules.SetupGnomePerformance()
		modules.SetupGnomeKeybinds()
	}
	modules.SetupShell(manager)
	modules.SetupEditors(manager)
	if !modules.IsGitHubAuthenticated() {
		fmt.Println("GitHub not authenticated — skipping (only required on Fresh Install). Run Fresh or gh auth login manually if needed.")
	}
	fmt.Printf("\n%s RESTORE COMPLETE — re-run anytime it breaks.\n", tui.MarkOK)
}

func installBaseTools(manager pkgmanager.PackageManager, sysInfo system.Info) {
	fmt.Println("\n--- Installing Base Tools ---")
	manager.Update()

	basePkgs := []string{
		"neovim", "grep", "ripgrep", "fzf", "zoxide", "curl", "wget",
		"git", "vim", "micro", "btop", "htop", "nvtop", "fastfetch", "alacritty", "jq", "wofi",
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
