package modules

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
)

func SetupNvidia(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	fmt.Println("\n--- NVIDIA Driver Setup ---")

	// Check if an NVIDIA GPU is present
	lspci, err := exec.Command("lspci").Output()
	if err == nil && !strings.Contains(strings.ToLower(string(lspci)), "nvidia") {
		fmt.Println("No NVIDIA GPU detected via lspci. Skipping installation.")
		return nil
	}

	switch sysInfo.OS {
	case "debian", "ubuntu", "pop", "linuxmint":
		return setupNvidiaDebian(manager, sysInfo)
	case "fedora":
		return setupNvidiaFedora(manager)
	case "arch", "manjaro":
		return setupNvidiaArch(manager)
	default:
		return fmt.Errorf("unsupported distribution for automated NVIDIA setup: %s", sysInfo.OS)
	}
}

func setupNvidiaDebian(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	if sysInfo.OS == "ubuntu" || sysInfo.OS == "pop" || sysInfo.OS == "linuxmint" {
		fmt.Println("Detected Ubuntu-based system. Using ubuntu-drivers autoinstall...")
		return pkgmanager.RunCommand("sudo", "ubuntu-drivers", "autoinstall")
	}

	fmt.Println("Detected Debian system. Ensuring non-free is enabled and installing nvidia-driver...")
	// Note: This is a simplified version. Enabling non-free usually requires editing sources.list
	// For now, we assume the user has non-free enabled or we try to install it.
	return manager.Install("nvidia-driver", "nvidia-settings", "nvidia-xconfig")
}

func setupNvidiaFedora(manager pkgmanager.PackageManager) error {
	fmt.Println("Ensuring RPM Fusion is enabled for Fedora...")
	
	// Enable RPM Fusion Free
	err := pkgmanager.RunCommand("sudo", "dnf", "install", "-y", "https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm")
	if err != nil {
		fmt.Printf("Warning: Failed to enable RPM Fusion Free (might already be enabled): %v\n", err)
	}

	// Enable RPM Fusion Nonfree
	err = pkgmanager.RunCommand("sudo", "dnf", "install", "-y", "https://mirrors.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm")
	if err != nil {
		fmt.Printf("Warning: Failed to enable RPM Fusion Nonfree (might already be enabled): %v\n", err)
	}

	fmt.Println("Installing NVIDIA drivers via dnf...")
	return manager.Install("akmod-nvidia", "xorg-x11-drv-nvidia-cuda", "nvidia-settings")
}

func setupNvidiaArch(manager pkgmanager.PackageManager) error {
	fmt.Println("Installing NVIDIA drivers via pacman...")
	return manager.Install("nvidia", "nvidia-utils", "nvidia-settings")
}
