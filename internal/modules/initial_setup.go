package modules

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

func RunInitialSetup(manager pkgmanager.PackageManager, sysInfo system.Info) error {
	items := []tui.ListItem{
		{Key: "fedora", Name: "Fedora Initial Setup"},
		{Key: "debian", Name: "Debian Initial Setup"},
		{Key: "ubuntu", Name: "Ubuntu Initial Setup"},
		{Key: "arch", Name: "Arch Initial Setup"},
	}

	// Auto-select the current OS item for convenience
	for i, item := range items {
		if item.Key == sysInfo.OS {
			items[i].Selected = true
		}
	}

	action, results, err := tui.RunListUI("OS Initial Setup", items)
	if err != nil || action == "" {
		return err
	}

	for _, item := range results {
		if !item.Selected {
			continue
		}

		switch item.Key {
		case "fedora":
			setupFedora(manager)
		case "debian":
			setupDebian(manager)
		case "ubuntu":
			setupUbuntu(manager)
		case "arch":
			setupArch(manager)
		}
	}

	return nil
}

func setupFedora(manager pkgmanager.PackageManager) {
	fmt.Println("\n--- Fedora Initial Setup ---")

	// 1. DNF Speedup
	fmt.Println("1. Optimizing DNF (max_parallel_downloads=10, fastestmirror=True)...")
	dnfConfig := `[main]
gpgcheck=True
installonly_limit=3
clean_requirements_on_remove=True
best=False
skip_if_unavailable=True
max_parallel_downloads=10
fastestmirror=True
`
	err := os.WriteFile("/tmp/dnf.conf", []byte(dnfConfig), 0644)
	if err == nil {
		pkgmanager.RunCommand("sudo", "cp", "/tmp/dnf.conf", "/etc/dnf/dnf.conf")
	}

	// 2. RPM Fusion
	fmt.Println("2. Enabling RPM Fusion (Free & Nonfree)...")
	pkgmanager.RunCommand("sudo", "dnf", "install", "-y", "https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm")
	pkgmanager.RunCommand("sudo", "dnf", "install", "-y", "https://mirrors.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm")

	// 3. DNS Configuration
	setupDNS()

	// 4. Update System
	fmt.Println("4. Updating the whole system...")
	manager.Update()
	manager.Upgrade()

	fmt.Println("\nFedora setup complete. Reboot is recommended.")
}

func setupDebian(manager pkgmanager.PackageManager) {
	fmt.Println("\n--- Debian Initial Setup ---")

	// 1. Enable Sudo (if not present) and Non-Free
	fmt.Println("1. Ensuring non-free and contrib components are enabled...")
	pkgmanager.RunCommand("sudo", "apt", "install", "-y", "software-properties-common")
	pkgmanager.RunCommand("sudo", "add-apt-repository", "contrib", "non-free", "non-free-firmware")

	// 2. DNS Configuration
	setupDNS()

	// 3. Update System
	fmt.Println("3. Updating system...")
	manager.Update()
	manager.Upgrade()

	fmt.Println("\nDebian setup complete.")
}

func setupUbuntu(manager pkgmanager.PackageManager) {
	fmt.Println("\n--- Ubuntu Initial Setup ---")

	// 1. Enable Multiverse and Universe
	fmt.Println("1. Enabling Universe and Multiverse repositories...")
	pkgmanager.RunCommand("sudo", "add-apt-repository", "-y", "universe")
	pkgmanager.RunCommand("sudo", "add-apt-repository", "-y", "multiverse")

	// 2. DNS Configuration
	setupDNS()

	// 3. Update System
	fmt.Println("3. Updating system...")
	manager.Update()
	manager.Upgrade()

	fmt.Println("\nUbuntu setup complete.")
}

func setupArch(manager pkgmanager.PackageManager) {
	fmt.Println("\n--- Arch Initial Setup ---")

	// 1. Parallel Downloads in Pacman
	fmt.Println("1. Optimizing Pacman (ParallelDownloads = 5)...")
	pkgmanager.RunCommand("sudo", "sed", "-i", "s/#ParallelDownloads = 5/ParallelDownloads = 10/", "/etc/pacman.conf")

	// 2. Reflector for fast mirrors
	fmt.Println("2. Installing and running Reflector for fastest mirrors...")
	manager.Install("reflector")
	pkgmanager.RunCommand("sudo", "reflector", "--latest", "20", "--protocol", "https", "--sort", "rate", "--save", "/etc/pacman.d/mirrorlist")

	// 3. DNS Configuration
	setupDNS()

	// 4. Update System
	fmt.Println("4. Updating system (Pacman)...")
	manager.Update()
	manager.Upgrade()

	fmt.Println("\nArch setup complete.")
}

func setupDNS() {
	fmt.Println("3. Configuring Google (8.8.8.8) and Cloudflare (1.1.1.1) DNS via nmcli...")
	// Get active connection name
	out, err := exec.Command("bash", "-c", "nmcli -t -f NAME,TYPE connection show --active | grep ethernet | head -1 | cut -d: -f1").Output()
	if err != nil || len(out) == 0 {
		out, err = exec.Command("bash", "-c", "nmcli -t -f NAME,TYPE connection show --active | grep wifi | head -1 | cut -d: -f1").Output()
	}

	connName := strings.TrimSpace(string(out))
	if connName != "" {
		fmt.Printf("Updating connection: %s\n", connName)
		pkgmanager.RunCommand("sudo", "nmcli", "connection", "modify", connName, "ipv4.dns", "1.1.1.1,8.8.8.8")
		pkgmanager.RunCommand("sudo", "nmcli", "connection", "modify", connName, "ipv4.ignore-auto-dns", "yes")
		pkgmanager.RunCommand("sudo", "nmcli", "connection", "up", connName)
	} else {
		fmt.Println("Warning: Could not detect an active NetworkManager connection to apply DNS settings.")
	}
}
