package modules

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rakesh/linutils-rakesh/internal/pkgmanager"
	"github.com/rakesh/linutils-rakesh/internal/system"
	"github.com/rakesh/linutils-rakesh/internal/tui"
)

var (
	alertStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF0000")).
			Padding(1, 4).
			MarginTop(1).
			MarginBottom(1)
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

	if isRebootRequired(sysInfo) {
		fmt.Println("\n" + alertStyle.Render("REBOOT REQUIRED: Significant system updates were applied. Please reboot your system now."))
	}

	return nil
}

func isRebootRequired(sysInfo system.Info) bool {
	switch sysInfo.OS {
	case "fedora":
		// dnf needs-restarting -r returns 1 if reboot is needed
		err := exec.Command("dnf", "needs-restarting", "-r").Run()
		if err != nil {
			return true
		}
	case "debian", "ubuntu", "pop", "linuxmint":
		if _, err := os.Stat("/var/run/reboot-required"); err == nil {
			return true
		}
	case "arch", "manjaro":
		// On Arch, if the running kernel's module directory is gone, a new kernel was installed
		out, err := exec.Command("uname", "-r").Output()
		if err == nil {
			kernelVer := strings.TrimSpace(string(out))
			moduleDir := fmt.Sprintf("/usr/lib/modules/%s", kernelVer)
			if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
				return true
			}
		}
	}
	return false
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

# --- Optimization Tweaks ---
max_parallel_downloads=10
fastestmirror=True
metadata_expire=86400
deltarpm=True
defaultyes=True
keepcache=True
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

	fmt.Println("\nFedora setup complete.")
}

func setupDebian(manager pkgmanager.PackageManager) {
	fmt.Println("\n--- Debian Initial Setup ---")

	// 1. Apt Speedup & Configuration
	fmt.Println("1. Optimizing Apt (Parallel downloads and better caching)...")
	aptConfig := `Binary::apt::APT::Keep-Downloaded-Packages "true";
Acquire::Languages "none";
Acquire::ParallelDownloads "10";
`
	err := os.WriteFile("/tmp/99parallel", []byte(aptConfig), 0644)
	if err == nil {
		pkgmanager.RunCommand("sudo", "cp", "/tmp/99parallel", "/etc/apt/apt.conf.d/99parallel")
	}

	// 2. Enable Contrib, Non-Free and Non-Free-Firmware
	fmt.Println("2. Ensuring non-free and contrib components are enabled...")
	pkgmanager.RunCommand("sudo", "apt", "install", "-y", "software-properties-common")
	// For Debian 12+, non-free-firmware is a separate component
	pkgmanager.RunCommand("sudo", "sed", "-i", "s/main$/main contrib non-free non-free-firmware/g", "/etc/apt/sources.list")
	// Also try the standard command as fallback/addition
	pkgmanager.RunCommand("sudo", "add-apt-repository", "-y", "contrib")
	pkgmanager.RunCommand("sudo", "add-apt-repository", "-y", "non-free")
	pkgmanager.RunCommand("sudo", "add-apt-repository", "-y", "non-free-firmware")

	// 3. DNS Configuration
	setupDNS()

	// 4. Update System
	fmt.Println("4. Updating system...")
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
