package config

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

type HyprlandConfigurator struct {
	SysInfo system.Info
}

func (c *HyprlandConfigurator) Setup(manager pkgmanager.PackageManager) error {
	fmt.Printf("Setting up Hyprland on %s (%s)...\n", c.SysInfo.OS, c.SysInfo.SessionType)

	var pkgs []string
	switch c.SysInfo.OS {
	case "arch", "manjaro":
		pkgs = []string{"hyprland", "waybar", "wofi", "kitty", "swaybg", "grim", "slurp"}
	case "fedora":
		pkgs = []string{"hyprland", "waybar", "wofi", "kitty", "swaybg", "grim", "slurp"}
	case "debian", "ubuntu":
		pkgs = []string{"hyprland", "waybar", "wofi", "kitty", "swaybg", "grim", "slurp"}
	}

	if err := manager.Install(pkgs...); err != nil {
		return err
	}

	fmt.Println("Hyprland setup complete.")
	return nil
}
