package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

// ConfigureHyprlandExtras handles the installation of additional Hyprland tools
// and manages the symlinking of configurations from dotfiles.
func ConfigureHyprlandExtras(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Granular Hyprland Setup ---")

	// 1. Install additional Hyprland tools
	pkgs := []string{"hypridle", "hyprlock", "hyprpaper", "hyprsunset"}
	fmt.Printf("Installing Hyprland tools: %v\n", pkgs)

	// We try to install all, but we should be aware that some might not be available
	// on all distributions yet (e.g. hyprsunset).
	if err := manager.Install(pkgs...); err != nil {
		fmt.Printf("Note: Some packages might not be available in your distribution's repositories: %v\n", err)
	}

	// 2. Symlink configuration
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	dotfilesHyprPath := filepath.Join(home, ".dotfiles", "hyprland")
	configHyprPath := filepath.Join(home, ".config", "hypr")

	// Verify if the source directory exists in ~/.dotfiles
	if _, err := os.Stat(dotfilesHyprPath); os.IsNotExist(err) {
		fmt.Printf("Configuration source not found at %s. Skipping symlink step.\n", dotfilesHyprPath)
		return nil
	}

	var confirmSymlink bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Symlink Hyprland Configuration?").
				Description(fmt.Sprintf("This will link %s to %s", dotfilesHyprPath, configHyprPath)).
				Value(&confirmSymlink),
		),
	)

	if err := confirmForm.Run(); err != nil {
		return err
	}

	if confirmSymlink {
		// Check if destination already exists
		if _, err := os.Lstat(configHyprPath); err == nil {
			var confirmOverwrite bool
			overwriteForm := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title("Target Already Exists").
						Description(fmt.Sprintf("%s already exists. Overwrite with symlink?", configHyprPath)).
						Value(&confirmOverwrite),
				),
			)
			if err := overwriteForm.Run(); err != nil {
				return err
			}

			if !confirmOverwrite {
				fmt.Println("Symlink operation cancelled by user.")
				return nil
			}

			// Backup existing configuration
			backupPath := configHyprPath + ".bak"
			fmt.Printf("Backing up existing configuration to %s\n", backupPath)
			// Remove existing backup if it exists to avoid rename failure
			os.RemoveAll(backupPath)
			if err := os.Rename(configHyprPath, backupPath); err != nil {
				// If rename fails, it might be because it's a symlink already or cross-device
				// Try to remove it if it's a symlink or if user confirmed overwrite
				if err := os.RemoveAll(configHyprPath); err != nil {
					return fmt.Errorf("failed to remove existing configuration: %v", err)
				}
			}
		}

		// Ensure the parent directory (~/.config) exists
		if err := os.MkdirAll(filepath.Dir(configHyprPath), 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}

		// Create the symlink
		fmt.Printf("Creating symlink: %s -> %s\n", dotfilesHyprPath, configHyprPath)
		if err := os.Symlink(dotfilesHyprPath, configHyprPath); err != nil {
			return fmt.Errorf("failed to create symlink: %v", err)
		}
		fmt.Println("Hyprland configuration symlinked successfully.")
	}

	return nil
}
