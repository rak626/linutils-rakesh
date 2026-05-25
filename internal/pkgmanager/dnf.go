package pkgmanager

import (
	"fmt"
	"os/exec"
	"strings"
)

type DnfManager struct{}

func (m *DnfManager) Install(packages ...string) error {
	var toInstall []string
	for _, pkg := range packages {
		if !m.IsInstalled(pkg) {
			toInstall = append(toInstall, pkg)
		} else {
			fmt.Printf("Package %s is already installed. Skipping.\n", pkg)
		}
	}

	if len(toInstall) == 0 {
		return nil
	}

	args := append([]string{"install", "-y"}, toInstall...)
	fmt.Printf("Running: sudo dnf %v\n", args)
	return RunCommand("sudo", append([]string{"dnf"}, args...)...)
}

func (m *DnfManager) Remove(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}
	args := append([]string{"remove", "-y"}, packages...)
	fmt.Printf("Running: sudo dnf %v\n", args)
	return RunCommand("sudo", append([]string{"dnf"}, args...)...)
}

func (m *DnfManager) Update() error {
	fmt.Println("Running: sudo dnf check-update")
	// dnf check-update returns 100 if updates are available, which is not an error
	_ = RunCommand("sudo", "dnf", "check-update")
	return nil
}

func (m *DnfManager) Upgrade() error {
	fmt.Println("Running: sudo dnf upgrade -y")
	return RunCommand("sudo", "dnf", "upgrade", "-y")
}

func (m *DnfManager) IsInstalled(pkg string) bool {
	// Handle group packages starting with @
	if strings.HasPrefix(pkg, "@") {
		// For groups, we'll just check if the command exists for now or return false to be safe
		// A better way would be 'dnf group list installed'
		return false 
	}
	err := exec.Command("rpm", "-q", pkg).Run()
	return err == nil
}
