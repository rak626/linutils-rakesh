package modules

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func InstallIconAssets(manager pkgmanager.PackageManager) error {
	fmt.Println("\n--- Installing Icon & Cursor Assets ---")

	sysInfo := system.GetSystemInfo()
	var packages []string

	if sysInfo.OS == "arch" || sysInfo.OS == "manjaro" {
		packages = []string{
			"papirus-icon-theme",
			"rose-pine-cursor",
			"rose-pine-icons",
			"catppuccin-cursors-macchiato",
			"gruvbox-plus-icon-pack-git",
		}
	} else {
		packages = []string{
			"papirus-icon-theme",
		}
	}

	return manager.Install(packages...)
}
