package modules

import (
	"fmt"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func SetupFlatpak(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	fmt.Println("\n--- Configuring Flatpak ---")

	var pkgs []string
	pkgs = append(pkgs, "flatpak")

	// Distro-specific GNOME Software plugin
	switch sysInfo.OS {
	case "debian", "ubuntu", "pop", "linuxmint":
		pkgs = append(pkgs, "gnome-software-plugin-flatpak")
	case "fedora":
		// Usually pre-installed on Fedora, but ensures it's there
		pkgs = append(pkgs, "gnome-software") 
	case "arch", "manjaro":
		pkgs = append(pkgs, "gnome-software")
	}

	fmt.Printf("Installing Flatpak and GNOME integration for %s...\n", sysInfo.OS)
	if err := manager.Install(pkgs...); err != nil {
		return fmt.Errorf("failed to install flatpak packages: %v", err)
	}

	// Add Flathub repository
	fmt.Println("Adding Flathub repository...")
	err := pkgmanager.RunCommand("flatpak", "remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo")
	if err != nil {
		return fmt.Errorf("failed to add flathub remote: %v", err)
	}

	fmt.Println("Flatpak configuration complete. Note: You might need to restart GNOME Software or your session to see Flatpaks.")
	return nil
}
