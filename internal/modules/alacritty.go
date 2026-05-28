package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

// SetupAlacritty installs Alacritty and handles its configuration.
// It offers to symlink ~/.dotfiles/alacritty to ~/.config/alacritty if it exists.
func SetupAlacritty(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Terminal Enhancements: Alacritty ---")

	// 1. Install Alacritty
	if !manager.IsInstalled("alacritty") {
		fmt.Println("Installing Alacritty...")
		if err := manager.Install("alacritty"); err != nil {
			return fmt.Errorf("failed to install alacritty: %v", err)
		}
	} else {
		fmt.Println("Alacritty is already installed.")
	}

	// 2. Handle Configuration
	home, _ := os.UserHomeDir()
	dotfilesAlacritty := filepath.Join(home, ".dotfiles", "alacritty")

	if _, err := os.Stat(dotfilesAlacritty); err == nil {
		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Alacritty configuration found in ~/.dotfiles.").
					Description("Would you like to symlink it using GNU Stow?").
					Value(&confirm),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		if confirm {
			// Ensure stow is installed
			if !pkgmanager.IsCommandAvailable("stow") {
				fmt.Println("Installing GNU Stow for symlinking...")
				if err := manager.Install("stow"); err != nil {
					return fmt.Errorf("failed to install stow: %v", err)
				}
			}

			fmt.Println("Symlinking Alacritty configuration...")
			dotfilesDir := filepath.Join(home, ".dotfiles")
			// stow -v -R -t ~ alacritty
			cmd := exec.Command("stow", "-v", "-R", "-t", home, "alacritty")
			cmd.Dir = dotfilesDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to stow alacritty: %v", err)
			}
			fmt.Println("Alacritty configuration symlinked successfully!")
		}
	} else {
		fmt.Printf("No Alacritty configuration found at %s. Skipping symlink.\n", dotfilesAlacritty)
	}

	return nil
}
