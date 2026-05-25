package modules

import (
	"fmt"
	"os/exec"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
)

type GnomeExtension struct {
	PackageName string
	ExtensionID string
}

var BloatExtensions = []GnomeExtension{
	{PackageName: "gnome-extensions-app", ExtensionID: ""},
	{PackageName: "gnome-shell-extension-apps-menu", ExtensionID: "apps-menu@gnome-shell-extensions.gcampax.github.com"},
	{PackageName: "gnome-shell-extension-background-logo", ExtensionID: "background-logo@fedorahosted.org"},
	{PackageName: "gnome-shell-extension-launch-new-instance", ExtensionID: "launch-new-instance@gnome-shell-extensions.gcampax.github.com"},
	{PackageName: "gnome-shell-extension-places-menu", ExtensionID: "places-menu@gnome-shell-extensions.gcampax.github.com"},
	{PackageName: "gnome-shell-extension-window-list", ExtensionID: "window-list@gnome-shell-extensions.gcampax.github.com"},
}

func GetBloatExtensions() []string {
	exts := make([]string, len(BloatExtensions))
	for i, e := range BloatExtensions {
		exts[i] = e.PackageName
	}
	return exts
}

func RemoveExtensions(manager pkgmanager.PackageManager, packages []string) error {
	hasGnomeExt := pkgmanager.IsCommandAvailable("gnome-extensions")

	for _, pkg := range packages {
		// Find the corresponding ExtensionID
		var extID string
		for _, e := range BloatExtensions {
			if e.PackageName == pkg {
				extID = e.ExtensionID
				break
			}
		}

		// Disable extension if it has an ID
		if extID != "" && hasGnomeExt {
			fmt.Printf("Disabling extension: %s\n", extID)
			exec.Command("gnome-extensions", "disable", extID).Run()
		}
	}

	// Remove packages
	var toRemove []string
	for _, pkg := range packages {
		if manager.IsInstalled(pkg) {
			toRemove = append(toRemove, pkg)
		}
	}

	if len(toRemove) > 0 {
		fmt.Printf("Removing packages: %v\n", toRemove)
		return manager.Remove(toRemove...)
	}
	
	return nil
}

func InstallExtensions(manager pkgmanager.PackageManager, packages []string) error {
	// Install packages
	if len(packages) > 0 {
		fmt.Printf("Installing packages: %v\n", packages)
		if err := manager.Install(packages...); err != nil {
			return err
		}
	}

	hasGnomeExt := pkgmanager.IsCommandAvailable("gnome-extensions")
	if !hasGnomeExt {
		return nil
	}

	for _, pkg := range packages {
		// Find the corresponding ExtensionID
		var extID string
		for _, e := range BloatExtensions {
			if e.PackageName == pkg {
				extID = e.ExtensionID
				break
			}
		}

		// Enable extension if it has an ID
		if extID != "" {
			fmt.Printf("Enabling extension: %s\n", extID)
			exec.Command("gnome-extensions", "enable", extID).Run()
		}
	}
	return nil
}
