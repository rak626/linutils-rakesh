package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

// ApplyThemes handles the application of Rose Pine and Everforest themes
// by symlinking theme files from the dotfiles directory to their respective
// application configuration paths.
func ApplyThemes(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Application Themes ---")

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	// Identify dotfiles location
	dotfilesBase := filepath.Join(home, ".dotfiles")
	if _, err := os.Stat(dotfilesBase); os.IsNotExist(err) {
		return fmt.Errorf("dotfiles directory not found in ~/.dotfiles")
	}

	// Let the user select which themes to focus on
	var selectedThemes []string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Themes to Apply").
				Description("Focusing on Rose Pine and Everforest themes found in ~/.dotfiles.").
				Options(
					huh.NewOption("Rose Pine", "rose-pine"),
					huh.NewOption("Everforest", "everforest"),
				).
				Value(&selectedThemes),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if len(selectedThemes) == 0 {
		fmt.Println("No themes selected.")
		return nil
	}

	// Apply folder-level symlinks as per instructions
	fmt.Println("\nConfiguring theme directories...")

	// 1. Rofi Themes
	// If ~/.dotfiles/rofi/themes exists, symlink to ~/.config/rofi/themes
	rofiSource := filepath.Join(dotfilesBase, "rofi", "themes")
	rofiDest := filepath.Join(home, ".config", "rofi", "themes")
	if _, err := os.Stat(rofiSource); err == nil {
		fmt.Printf("Linking Rofi themes folder: %s\n", rofiSource)
		if err := ensureSymlink(rofiSource, rofiDest); err != nil {
			fmt.Printf("Warning: failed to link Rofi themes: %v\n", err)
		} else {
			fmt.Println("  - Rofi themes linked to ~/.config/rofi/themes")
		}
	}

	// 2. Ulauncher Themes
	// If ~/.dotfiles/ulauncher/themes exists, symlink to ~/.config/ulauncher/user-themes
	ulauncherSource := filepath.Join(dotfilesBase, "ulauncher", "themes")
	ulauncherDest := filepath.Join(home, ".config", "ulauncher", "user-themes")
	if _, err := os.Stat(ulauncherSource); err == nil {
		fmt.Printf("Linking Ulauncher themes folder: %s\n", ulauncherSource)
		if err := ensureSymlink(ulauncherSource, ulauncherDest); err != nil {
			fmt.Printf("Warning: failed to link Ulauncher themes: %v\n", err)
		} else {
			fmt.Println("  - Ulauncher themes linked to ~/.config/ulauncher/user-themes")
		}
	}

	// 3. GTK Themes and Theme-specific logic
	for _, theme := range selectedThemes {
		fmt.Printf("\nChecking for specific %s theme files...\n", theme)
		
		// Check for GTK theme in ~/.dotfiles/gtk/<theme>
		gtkSource := filepath.Join(dotfilesBase, "gtk", theme)
		if _, err := os.Stat(gtkSource); err == nil {
			gtkDest := filepath.Join(home, ".themes", theme)
			fmt.Printf("Linking GTK theme: %s\n", theme)
			if err := ensureSymlink(gtkSource, gtkDest); err != nil {
				fmt.Printf("Warning: failed to link GTK theme %s: %v\n", theme, err)
			} else {
				fmt.Printf("  - GTK theme %s linked to ~/.themes/%s\n", theme, theme)
			}
		}

		// Additional theme-specific logic could be added here if needed, 
		// such as individual file symlinks if the folder-level link is not enough.
	}

	fmt.Println("\nThemes application complete!")
	return nil
}

// ensureSymlink ensures the destination parent directory exists, 
// removes any existing file or directory at the destination, 
// and creates a new symbolic link.
func ensureSymlink(source, dest string) error {
	// Create destination parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	// Remove existing destination (file, symlink, or directory) to avoid conflicts
	if _, err := os.Lstat(dest); err == nil {
		if err := os.RemoveAll(dest); err != nil {
			return fmt.Errorf("failed to remove existing destination: %v", err)
		}
	}

	// Create the symlink
	return os.Symlink(source, dest)
}
