package pkgmanager

import (
	"fmt"
	"os/exec"
)

type AptManager struct{}

func (m *AptManager) Install(packages ...string) error {
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
	fmt.Printf("Running: sudo apt %v\n", args)
	return RunCommand("sudo", append([]string{"apt"}, args...)...)
}

func (m *AptManager) Remove(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}
	args := append([]string{"purge", "-y"}, packages...)
	fmt.Printf("Running: sudo apt %v\n", args)
	return RunCommand("sudo", append([]string{"apt"}, args...)...)
}

func (m *AptManager) Update() error {
	fmt.Println("Running: sudo apt update")
	return RunCommand("sudo", "apt", "update")
}

func (m *AptManager) Upgrade() error {
	fmt.Println("Running: sudo apt upgrade -y")
	return RunCommand("sudo", "apt", "upgrade", "-y")
}

func (m *AptManager) IsInstalled(pkg string) bool {
	err := exec.Command("dpkg", "-l", pkg).Run()
	return err == nil
}
