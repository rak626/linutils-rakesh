package modules

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func SetupI3(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	fmt.Printf("\n--- Setting up i3wm on %s (%s) ---\n", sysInfo.OS, sysInfo.SessionType)

	var pkgs []string
	switch sysInfo.OS {
	case "arch", "manjaro":
		pkgs = []string{"i3-wm", "i3status", "i3lock", "dmenu", "kitty", "feh"}
	case "fedora":
		pkgs = []string{"i3", "i3status", "i3lock", "dmenu", "kitty", "feh"}
	case "debian", "ubuntu":
		pkgs = []string{"i3", "i3status", "i3lock", "dmenu", "kitty", "feh"}
	}

	if err := manager.Install(pkgs...); err != nil {
		return err
	}

	fmt.Println("i3wm setup complete.")
	return nil
}
