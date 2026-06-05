package modules

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

// ConfigureHyprlandExtras handles the installation of additional Hyprland tools.
func ConfigureHyprlandExtras(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Granular Hyprland Setup ---")

	// 1. Install additional Hyprland tools
	pkgs := []string{"hypridle", "hyprlock", "hyprpaper", "hyprsunset"}
	fmt.Printf("Installing Hyprland tools: %v\n", pkgs)

	if err := manager.Install(pkgs...); err != nil {
		fmt.Printf("Note: Some packages might not be available in your distribution's repositories: %v\n", err)
	}

	return nil
}
