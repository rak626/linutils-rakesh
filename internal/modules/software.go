package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rakesh/linutils-rakesh/internal/config"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

func InstallSoftware(manager pkgmanager.PackageManager, sysInfo system.Info, items []tui.ListItem) ([]tui.ListItem, error) {
	fmt.Println("\n--- Starting Automatic Software Installation ---")

	// 1. General Software
	fmt.Println("\n[General Software]")
	fmt.Println("Installing Chromium Browser...")
	manager.Install("chromium")

	// 2. Manual Installs
	fmt.Println("\n[Manual Installs]")
	for _, inst := range config.ManualInstalls {
		installFromConfig(manager, sysInfo, inst)
	}

	// 3. AI Tools
	fmt.Println("\n[AI Tools]")
	for _, inst := range config.AIInstalls {
		installFromConfig(manager, sysInfo, inst)
	}

	// 4. Helper Tools
	fmt.Println("\n[Helper Tools]")
	for _, inst := range config.HelperInstalls {
		installFromConfig(manager, sysInfo, inst)
	}

	// 5. Flatpak Installs
	fmt.Println("\n[Flatpak Installs]")
	if !isFlatpakReady() {
		fmt.Println("Setting up Flatpak first...")
		if err := SetupFlatpak(manager, sysInfo); err != nil {
			fmt.Printf("Error setting up Flatpak: %v\n", err)
		}
	}
	if isFlatpakReady() {
		for _, inst := range config.FlatpakInstalls {
			installFromConfig(manager, sysInfo, inst)
		}
	}

	fmt.Println("\n--- Automatic Software Installation Complete ---")
	return items, nil
}

func isFlatpakReady() bool {
	_, err := exec.LookPath("flatpak")
	if err != nil {
		return false
	}
	// Also check if flathub is added
	cmd := exec.Command("flatpak", "remotes")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "flathub")
}

func osGroup(osName string) string {
	switch osName {
	case "debian", "ubuntu", "pop", "linuxmint":
		return "apt"
	case "arch", "manjaro":
		return "arch"
	case "fedora":
		return "fedora"
	}
	return "default"
}

func installFromConfig(manager pkgmanager.PackageManager, sysInfo system.Info, inst config.InstallConfig) {
	if inst.Check != "" {
		if strings.HasPrefix(inst.Check, "~/") {
			home := os.Getenv("HOME")
			path := filepath.Join(home, inst.Check[2:])
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("✓ %s is already installed.\n", inst.Name)
				return
			}
		} else {
			if _, err := exec.LookPath(inst.Check); err == nil {
				fmt.Printf("✓ %s is already installed.\n", inst.Name)
				return
			}
		}
	}

	// Install dependencies if any
	if len(inst.Deps) > 0 {
		fmt.Printf("Checking dependencies for %s: %v\n", inst.Name, inst.Deps)
		var missingDeps []string
		for _, dep := range inst.Deps {
			if !manager.IsInstalled(dep) {
				missingDeps = append(missingDeps, dep)
			}
		}

		if len(missingDeps) > 0 {
			fmt.Printf("Installing missing dependencies: %v\n", missingDeps)
			if err := manager.Install(missingDeps...); err != nil {
				fmt.Printf("Warning: Failed to install dependencies: %v\n", err)
			}
		}
	}

	// Pick commands: CommandByOS takes priority, fall back to Command
	cmds := inst.Command
	group := osGroup(sysInfo.OS)
	if osCmds, ok := inst.CommandByOS[group]; ok {
		cmds = osCmds
	} else if osCmds, ok := inst.CommandByOS[sysInfo.OS]; ok {
		cmds = osCmds
	}

	fmt.Printf("Installing %s...\n", inst.Name)
	runCommands(cmds)
}

func runCommands(commands []string) {
	for _, cmdStr := range commands {
		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Step failed: %v\n", err)
			return
		}
	}
}
