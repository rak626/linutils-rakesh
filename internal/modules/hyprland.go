package modules

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func SetupHyprland(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	fmt.Printf("\n--- Setting up Hyprland on %s (%s) ---\n", sysInfo.OS, sysInfo.SessionType)

	var pkgs []string
	switch sysInfo.OS {
	case "arch", "manjaro":
		pkgs = []string{
			"hyprland", "waybar", "wofi", "alacritty", "hyprpaper", "grim", "slurp",
			"wl-clipboard", "mako", "swaylock-effects", "xdg-desktop-portal-hyprland",
			"polkit-kde-agent", "qt5-wayland", "qt6-wayland",
			"pipewire", "wireplumber", "pipewire-pulse", "pavucontrol",
		}
	case "fedora":
		pkgs = []string{"hyprland", "waybar", "wofi", "alacritty", "hyprpaper", "grim", "slurp"}
	case "debian", "ubuntu":
		pkgs = []string{"hyprland", "waybar", "wofi", "alacritty", "hyprpaper", "grim", "slurp"}
	}

	if err := manager.Install(pkgs...); err != nil {
		return err
	}

	fmt.Println("Hyprland setup complete.")
	return nil
}
