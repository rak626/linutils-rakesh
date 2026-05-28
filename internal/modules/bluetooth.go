package modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

// SetupBluetoothAndAudio implements Omarchy-style Bluetooth & Audio Management.
func SetupBluetoothAndAudio(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	var confirm bool
	err := huh.NewConfirm().
		Title("Bluetooth & Audio Management").
		Description("Install Omarchy-style Bluetooth (bluetui) and Audio (wiremix) managers?").
		Affirmative("Install").
		Negative("Skip").
		Value(&confirm).
		Run()

	if err != nil {
		return err
	}

	if !confirm {
		return nil
	}

	fmt.Println("\n--- Setting up Bluetooth & Audio Management ---")

	if sysInfo.OS == "arch" || sysInfo.OS == "manjaro" {
		fmt.Println("Installing bluetui, wiremix, bluez, bluez-utils via yay...")
		// Use yay as requested for Arch-based systems
		if err := pkgmanager.RunCommand("yay", "-S", "--noconfirm", "bluetui", "wiremix", "bluez", "bluez-utils"); err != nil {
			fmt.Printf("Warning: Failed to install via yay: %v. Falling back to pacman for base packages...\n", err)
			manager.Install("bluez", "bluez-utils")
		}
	} else {
		fmt.Printf("Installing standard packages for %s...\n", sysInfo.OS)
		// Fallback for other distros
		pkgs := []string{"bluez", "bluez-utils", "blueman", "pavucontrol"}
		if err := manager.Install(pkgs...); err != nil {
			return err
		}
	}

	fmt.Println("Enabling and starting bluetooth service...")
	if err := pkgmanager.RunCommand("sudo", "systemctl", "enable", "--now", "bluetooth"); err != nil {
		fmt.Printf("Warning: Failed to enable bluetooth service: %v\n", err)
	}

	// Automate Hyprland configuration if possible
	if sysInfo.DE == "hyprland" {
		home, err := os.UserHomeDir()
		if err == nil {
			hyprConf := filepath.Join(home, ".config/hypr/hyprland.conf")
			if _, err := os.Stat(hyprConf); err == nil {
				fmt.Println("Detected Hyprland, adding keybinds to hyprland.conf...")
				appendToFileIfMissing(hyprConf, "bind = $mainMod CTRL, B, exec, kitty --class floating -e bluetui")
				appendToFileIfMissing(hyprConf, "bind = $mainMod CTRL, A, exec, kitty --class floating -e wiremix")
			} else {
				fmt.Println("\nAdd these binds to your hyprland.conf:")
				fmt.Println("bind = $mainMod CTRL, B, exec, kitty --class floating -e bluetui")
				fmt.Println("bind = $mainMod CTRL, A, exec, kitty --class floating -e wiremix")
			}
		}
	}

	fmt.Println("\nBluetooth & Audio management setup complete.")
	return nil
}
