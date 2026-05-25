package config

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

type GnomeConfigurator struct {
	SysInfo system.Info
}

func (c *GnomeConfigurator) Setup(manager pkgmanager.PackageManager) error {
	fmt.Printf("Optimizing Gnome %s on %s...\n", c.SysInfo.DEVersion, c.SysInfo.OS)

	var pkgs []string
	switch c.SysInfo.OS {
	case "arch", "manjaro":
		pkgs = []string{"gnome-tweaks", "gnome-shell-extensions"}
	case "fedora":
		pkgs = []string{"gnome-tweaks", "gnome-extensions-app"}
	case "debian", "ubuntu":
		pkgs = []string{"gnome-tweaks", "gnome-shell-extensions"}
	}

	if err := manager.Install(pkgs...); err != nil {
		return err
	}

	fmt.Println("Gnome setup complete.")
	return nil
}
