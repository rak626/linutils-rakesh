package modules

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

func SetupGnomePerformance() error {
	if !pkgmanager.IsCommandAvailable("gsettings") {
		return fmt.Errorf("gsettings command not found. This module only works on GNOME")
	}

	fmt.Println("\n--- Applying GNOME Performance & UI Improvements ---")

	settings := [][]string{
		{"set", "org.gnome.desktop.interface", "enable-hot-corners", "false"},
		{"set", "org.gnome.desktop.interface", "enable-animations", "false"},
		{"set", "org.gnome.desktop.interface", "show-battery-percentage", "true"},
		{"set", "org.gnome.mutter", "overlay-key", "''"},
	}

	for _, s := range settings {
		fmt.Printf("Setting %s %s to %s...\n", s[1], s[2], s[3])
		pkgmanager.RunCommand("gsettings", s...)
	}

	fmt.Println("GNOME performance improvements applied.")
	return nil
}
