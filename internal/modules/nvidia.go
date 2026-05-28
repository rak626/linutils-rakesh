package modules

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	var installErr error
	switch sysInfo.OS {
	case "debian", "ubuntu", "pop", "linuxmint":
		installErr = setupNvidiaDebian(manager, sysInfo)
	case "fedora":
		installErr = setupNvidiaFedora(manager)
	case "arch", "manjaro":
		installErr = setupNvidiaArch(manager)
	default:
		return fmt.Errorf("unsupported distribution for automated NVIDIA setup: %s", sysInfo.OS)
	}

	if installErr != nil {
		return installErr
	}

	// Post-install logic: Environment Variables for Hyprland
	if sysInfo.DE == "hyprland" {
		return setupNvidiaHyprlandEnv()
	}

	return nil
}

func setupNvidiaDebian(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	if sysInfo.OS == "ubuntu" || sysInfo.OS == "pop" || sysInfo.OS == "linuxmint" {
		fmt.Println("Detected Ubuntu-based system. Using ubuntu-drivers autoinstall...")
		return pkgmanager.RunCommand("sudo", "ubuntu-drivers", "autoinstall")
	}

	fmt.Println("Detected Debian system. Installing nvidia-driver...")
	return manager.Install("nvidia-driver", "nvidia-settings", "nvidia-xconfig")
}

func setupNvidiaFedora(manager pkgmanager.PackageManager) error {
	fmt.Println("Installing NVIDIA drivers via dnf...")
	// Note: Assumes RPM Fusion is already enabled by Initial Setup
	return manager.Install("akmod-nvidia", "xorg-x11-drv-nvidia-cuda", "nvidia-settings", "libva-nvidia-driver")
}

func setupNvidiaArch(manager pkgmanager.PackageManager) error {
	fmt.Println("Installing NVIDIA drivers via pacman (DKMS)...")
	// Use nvidia-dkms for better kernel update support
	return manager.Install("nvidia-dkms", "nvidia-utils", "nvidia-settings", "libva-nvidia-driver")
}

func setupNvidiaHyprlandEnv() error {
	fmt.Println("\n--- Configuring NVIDIA Environment Variables for Hyprland ---")
	
	home, _ := os.UserHomeDir()
	hyprConfig := filepath.Join(home, ".config", "hypr", "hyprland.conf")

	if _, err := os.Stat(hyprConfig); os.IsNotExist(err) {
		fmt.Println("Hyprland config not found. Skipping environment variable injection.")
		return nil
	}

	envVars := []string{
		"env = LIBVA_DRIVER_NAME,nvidia",
		"env = XDG_SESSION_TYPE,wayland",
		"env = GBM_BACKEND,nvidia-drm",
		"env = __GLX_VENDOR_LIBRARY_NAME,nvidia",
		"env = WLR_NO_HARDWARE_CURSORS,1",
	}

	fmt.Println("Injecting Wayland/NVIDIA variables into hyprland.conf...")
	for _, env := range envVars {
		if err := appendToFileIfMissing(hyprConfig, env); err != nil {
			return fmt.Errorf("failed to append %s: %v", env, err)
		}
	}

	fmt.Println("NVIDIA environment variables configured successfully!")
	return nil
}
