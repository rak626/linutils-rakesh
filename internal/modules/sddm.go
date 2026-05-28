package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

// SetupSDDM configures the SDDM login manager with a custom theme.
func SetupSDDM(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	var installSDDM bool
	confirm := huh.NewConfirm().
		Title("SDDM Login Manager").
		Description("Do you want to install and enable SDDM as your login manager?").
		Value(&installSDDM)

	err := confirm.Run()
	if err != nil {
		return err
	}

	if !installSDDM {
		fmt.Println("Skipping SDDM setup.")
		return nil
	}

	fmt.Printf("\n--- Setting up SDDM on %s ---\n", sysInfo.OS)

	// Disabling other display managers to avoid conflicts
	fmt.Println("Attempting to disable other display managers (gdm, lightdm, lxdm, slim)...")
	dms := []string{"gdm", "lightdm", "lxdm", "slim"}
	for _, dm := range dms {
		// We use RunCommand directly and ignore errors as these services might not exist
		_ = pkgmanager.RunCommand("sudo", "systemctl", "disable", dm)
	}

	switch sysInfo.OS {
	case "arch", "manjaro":
		fmt.Println("Installing SDDM and sugar-candy theme via yay...")
		// Packages: sddm, sddm-sugar-candy-git, qt5-graphicaleffects, qt5-quickcontrols2, qt5-svg
		// Note: sddm-sugar-candy-git is an AUR package, so we use yay.
		if err := pkgmanager.RunCommand("yay", "-S", "--noconfirm", "sddm", "sddm-sugar-candy-git", "qt5-graphicaleffects", "qt5-quickcontrols2", "qt5-svg"); err != nil {
			fmt.Printf("Warning: Failed to install via yay: %v. Falling back to pacman for sddm and dependencies...\n", err)
			if err := manager.Install("sddm", "qt5-graphicaleffects", "qt5-quickcontrols2", "qt5-svg"); err != nil {
				return err
			}
			fmt.Println("Note: sddm-sugar-candy-git (AUR) was not installed. Theme configuration may fail.")
		}

		// Configure theme
		confDir := "/etc/sddm.conf.d"
		confFile := filepath.Join(confDir, "default.conf")
		content := "[Theme]\nCurrent=sugar-candy\n"

		fmt.Println("Configuring SDDM theme (sugar-candy)...")
		if err := pkgmanager.RunCommand("sudo", "mkdir", "-p", confDir); err != nil {
			return err
		}

		tempFile := "/tmp/sddm_default.conf"
		if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write temporary config file: %v", err)
		}
		if err := pkgmanager.RunCommand("sudo", "cp", tempFile, confFile); err != nil {
			return fmt.Errorf("failed to copy config file to %s: %v", confFile, err)
		}

	case "fedora":
		fmt.Println("Installing SDDM and Fedora theme...")
		// Fedora uses sddm-themes-fedora or similar
		if err := manager.Install("sddm", "sddm-themes-fedora"); err != nil {
			return err
		}
	case "debian", "ubuntu", "pop", "linuxmint":
		fmt.Println("Installing SDDM and Breeze theme...")
		// Ubuntu/Debian often use sddm-theme-breeze
		if err := manager.Install("sddm", "sddm-theme-breeze"); err != nil {
			return err
		}
	default:
		fmt.Println("Installing SDDM...")
		if err := manager.Install("sddm"); err != nil {
			return err
		}
	}

	fmt.Println("Enabling SDDM service...")
	if err := pkgmanager.RunCommand("sudo", "systemctl", "enable", "sddm"); err != nil {
		return fmt.Errorf("failed to enable sddm service: %v", err)
	}

	fmt.Println("SDDM setup complete. Please reboot for changes to take effect.")
	return nil
}
