package modules

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func SetupHyprland(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	if sysInfo.OS != "arch" && sysInfo.OS != "manjaro" {
		fmt.Printf("Hyprland setup is currently only optimized for Arch-based distributions. Detected: %s\n", sysInfo.OS)
		return nil
	}

	fmt.Printf("\n--- Setting up Hyprland on %s (%s) ---\n", sysInfo.OS, sysInfo.SessionType)

	pkgs := []string{
		"hyprland", "waybar", "wofi", "alacritty", "hyprpaper", "grim", "slurp", "satty",
		"wl-clipboard", "mako", "swaylock-effects", "xdg-desktop-portal-hyprland",
		"polkit-kde-agent", "qt5-wayland", "qt6-wayland",
		"pipewire", "wireplumber", "pipewire-pulse",
	}

	if err := manager.Install(pkgs...); err != nil {
		return err
	}

	fmt.Println("Hyprland setup complete.")
	return nil
}
