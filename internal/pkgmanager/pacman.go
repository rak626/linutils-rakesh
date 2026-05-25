package pkgmanager

import (
	"fmt"
	"os/exec"
)

type PacmanManager struct{}

func (m *PacmanManager) Install(packages ...string) error {
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

	args := append([]string{"-S", "--noconfirm"}, toInstall...)
	fmt.Printf("Running: sudo pacman %v\n", args)
	return RunCommand("sudo", append([]string{"pacman"}, args...)...)
}

func (m *PacmanManager) Remove(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}
	args := append([]string{"-Rs", "--noconfirm"}, packages...)
	fmt.Printf("Running: sudo pacman %v\n", args)
	return RunCommand("sudo", append([]string{"pacman"}, args...)...)
}

func (m *PacmanManager) Update() error {
	fmt.Println("Running: sudo pacman -Sy")
	return RunCommand("sudo", "pacman", "-Sy")
}

func (m *PacmanManager) Upgrade() error {
	fmt.Println("Running: sudo pacman -Syu --noconfirm")
	return RunCommand("sudo", "pacman", "-Syu", "--noconfirm")
}

func (m *PacmanManager) IsInstalled(pkg string) bool {
	err := exec.Command("pacman", "-Qs", pkg).Run()
	return err == nil
}
