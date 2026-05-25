package modules

import (
	"fmt"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

var BloatApps = []string{
	"gnome-contacts",
	"gnome-weather",
	"gnome-maps",
	"totem",
	"rhythmbox",
	"cheese",
	"simple-scan",
	"yelp",
	"gnome-tour",
	"mediawriter",
	"snapshot",
}

func GetBloatApps(sysInfo system.Info) []string {
	apps := make([]string, len(BloatApps))
	copy(apps, BloatApps)
	if sysInfo.OS == "fedora" || sysInfo.OS == "arch" || sysInfo.OS == "manjaro" {
		apps = append(apps, "gnome-software")
	}
	if sysInfo.OS == "fedora" {
		apps = append(apps, "PackageKit", "PackageKit-command-not-found")
	}
	return apps
}

func RemoveApps(manager pkgmanager.PackageManager, apps []string) error {
	var toRemove []string
	for _, p := range apps {
		if manager.IsInstalled(p) {
			toRemove = append(toRemove, p)
		} else {
			fmt.Printf("App %s is not installed. Skipping removal.\n", p)
		}
	}

	if len(toRemove) == 0 {
		return nil
	}

	fmt.Printf("Removing apps: %v\n", toRemove)
	return manager.Remove(toRemove...)
}

func InstallApps(manager pkgmanager.PackageManager, apps []string) error {
	fmt.Printf("Installing (Resetting) apps: %v\n", apps)
	return manager.Install(apps...)
}
