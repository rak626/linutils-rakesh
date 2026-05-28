package modules

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func InstallIconAssets(manager pkgmanager.PackageManager) error {
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Install Icon & Cursor Assets").
				Description("This will install various icon themes and cursor sets (Papirus, Rose Pine, Catppuccin, Gruvbox).").
				Value(&confirm),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if !confirm {
		return nil
	}

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
