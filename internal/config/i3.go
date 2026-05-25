package config

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

type I3Configurator struct {
	SysInfo system.Info
}

func (c *I3Configurator) Setup(manager pkgmanager.PackageManager) error {
	fmt.Printf("Setting up i3wm on %s...\n", c.SysInfo.OS)

	var pkgs []string
	switch c.SysInfo.OS {
	case "arch", "manjaro":
		pkgs = []string{"i3-wm", "i3status", "i3lock", "dmenu", "alacritty"}
	case "fedora":
		pkgs = []string{"i3", "i3status", "i3lock", "dmenu", "alacritty"}
	case "debian", "ubuntu":
		pkgs = []string{"i3", "i3status", "i3lock", "suckless-tools", "alacritty"}
	}

	if err := manager.Install(pkgs...); err != nil {
		return err
	}

	fmt.Println("i3wm installation complete.")
	return nil
}
