package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

// SetupFileManagers installs Thunar, Yazi, and other file management essentials.
func SetupFileManagers(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	fmt.Println("\n--- File Management Essentials ---")

	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Install File Management Essentials?").
				Description("Installs Thunar (GUI), Yazi (TUI), and Archive Manager with previews.").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if !confirm {
		fmt.Println("Skipping File Management Essentials setup.")
		return nil
	}

	if sysInfo.OS == "arch" || sysInfo.OS == "manjaro" {
		fmt.Println("Installing File Management suite via yay...")
		packages := []string{
			"thunar", "thunar-archive-plugin", "thunar-volman",
			"gvfs", "gvfs-mtp", "gvfs-afc",
			"tumbler", "ffmpegthumbnailer", "poppler-glib", "libgsf",
			"yazi", "file-roller",
		}

		// Use yay as requested
		if err := pkgmanager.RunCommand("yay", append([]string{"-S", "--noconfirm"}, packages...)...); err != nil {
			fmt.Printf("Warning: Failed to install via yay: %v. Falling back to pacman...\n", err)
			if err := manager.Install(packages...); err != nil {
				return fmt.Errorf("failed to install file management suite: %v", err)
			}
		}
	} else {
		fmt.Println("Installing Thunar, Yazi, and gvfs...")
		// Standard install for other distros
		if err := manager.Install("thunar", "yazi", "gvfs"); err != nil {
			return fmt.Errorf("failed to install file managers: %v", err)
		}
	}

	// Hyprland Keybindings
	if sysInfo.DE == "hyprland" {
		home, _ := os.UserHomeDir()
		hyprConf := filepath.Join(home, ".config", "hypr", "hyprland.conf")

		// Check if config exists
		if _, err := os.Stat(hyprConf); err == nil {
			fmt.Println("Adding file manager keybindings to Hyprland config...")

			keybindings := []string{
				"bind = $mainMod, E, exec, thunar",
				"bind = $mainMod SHIFT, E, exec, kitty yazi",
			}

			for _, kb := range keybindings {
				if err := appendToFileIfMissing(hyprConf, kb); err != nil {
					fmt.Printf("Warning: Failed to add keybinding to %s: %v\n", hyprConf, err)
				}
			}
		} else {
			fmt.Printf("Note: Hyprland config not found at %s. Skipping keybindings.\n", hyprConf)
		}
	}

	fmt.Println("File Management Essentials setup complete!")
	return nil
}
